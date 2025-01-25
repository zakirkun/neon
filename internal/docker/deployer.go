package docker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/pkg/archive"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/zakirkun/neon/internal/config"
	"github.com/zakirkun/neon/internal/config/compose"
	"github.com/zakirkun/neon/internal/config/deploy"
)

type Deployer struct {
	client *Client
	config *config.Config
}

func NewDeployer(client *Client, cfg *config.Config) *Deployer {
	return &Deployer{
		client: client,
		config: cfg,
	}
}

func (d *Deployer) Deploy(ctx context.Context, repoURL string) error {
	// 1. Clone repository
	repoPath, err := d.cloneRepo(repoURL)
	if err != nil {
		return err
	}
	defer os.RemoveAll(repoPath)

	// 2. Build Docker image
	imageName, err := d.buildImage(ctx, repoPath)
	if err != nil {
		return err
	}

	// 3. Push image ke registry
	err = d.pushImage(ctx, imageName)
	if err != nil {
		return err
	}

	// 4. Deploy ke Swarm
	return d.deployToSwarm(ctx, imageName)
}

func (d *Deployer) DeployFromConfig(ctx context.Context, svc *deploy.ServiceConfig) error {
	// Pull image dari registry
	if err := d.pullImage(ctx, svc.Image); err != nil {
		return fmt.Errorf("gagal pull image: %v", err)
	}

	// Convert ke service spec
	spec := &swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: svc.Name,
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: svc.Image,
				Env:   svc.Environment,
			},
			Resources: &swarm.ResourceRequirements{
				Limits: &swarm.Limit{
					NanoCPUs:    parseCPUs(svc.Deploy.Resources.Limits.CPUs),
					MemoryBytes: parseMemory(svc.Deploy.Resources.Limits.Memory),
				},
			},
			RestartPolicy: &swarm.RestartPolicy{
				Condition:   swarm.RestartPolicyCondition(svc.Deploy.RestartPolicy.Condition),
				MaxAttempts: &svc.Deploy.RestartPolicy.MaxAttempts,
			},
		},
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &svc.Replicas,
			},
		},
		UpdateConfig: &swarm.UpdateConfig{
			Parallelism: svc.Deploy.UpdateConfig.Parallelism,
			Delay:       parseDuration(svc.Deploy.UpdateConfig.Delay),
		},
		EndpointSpec: &swarm.EndpointSpec{
			Ports: convertPorts(svc.Ports),
		},
	}

	_, err := d.client.ServiceCreate(ctx, *spec, types.ServiceCreateOptions{})
	return err
}

func (d *Deployer) DeployComposeService(ctx context.Context, name string, service *compose.Service) error {
	// Build image jika diperlukan
	var imageName string
	if service.Build != nil {
		var err error
		imageName, err = d.buildImage(ctx, service.Build.Context)
		if err != nil {
			return fmt.Errorf("gagal build image: %v", err)
		}
	} else {
		imageName = service.Image
	}

	// Convert port mappings
	ports := make([]swarm.PortConfig, 0)
	for _, portStr := range service.Ports {
		port, err := parsePortConfig(portStr)
		if err != nil {
			return err
		}
		ports = append(ports, port)
	}

	// Convert environment variables
	env := make([]string, 0)
	for k, v := range service.Environment {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	// Create service spec
	replicas := uint64(service.Deploy.Replicas)
	spec := &swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: name,
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image:   imageName,
				Env:     env,
				Command: []string{service.Command},
			},
			Resources: &swarm.ResourceRequirements{
				Limits: &swarm.Limit{
					NanoCPUs:    parseCPUs(service.Deploy.Resources.Limits.CPUs),
					MemoryBytes: parseMemory(service.Deploy.Resources.Limits.Memory),
				},
			},
			RestartPolicy: &swarm.RestartPolicy{
				Condition: swarm.RestartPolicyCondition(service.Deploy.Restart.Condition),
				MaxAttempts: func() *uint64 {
					v := uint64(service.Deploy.Restart.MaxAttempts)
					return &v
				}(),
			},
		},
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &replicas,
			},
		},
		UpdateConfig: &swarm.UpdateConfig{
			Parallelism: uint64(service.Deploy.UpdateConfig.Parallelism),
			Delay:       parseDuration(service.Deploy.UpdateConfig.Delay),
		},
		EndpointSpec: &swarm.EndpointSpec{
			Ports: ports,
		},
	}

	_, err := d.client.ServiceCreate(ctx, *spec, types.ServiceCreateOptions{})
	return err
}

func (d *Deployer) cloneRepo(repoURL string) (string, error) {
	// Implementasi Git clone menggunakan go-git
	return d.cloneRepository(context.Background(), repoURL, "main")
}

func (d *Deployer) cloneRepository(ctx context.Context, repoURL string, branch string) (string, error) {
	// Buat temporary directory untuk menyimpan hasil clone
	tmpDir, err := os.MkdirTemp("", "repo-*")
	if err != nil {
		return "", fmt.Errorf("gagal membuat temporary directory: %w", err)
	}

	// Clone options
	cloneOpts := &git.CloneOptions{
		URL:           repoURL,
		Progress:      os.Stdout,
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	}

	// Lakukan git clone
	_, err = git.PlainCloneContext(ctx, tmpDir, false, cloneOpts)
	if err != nil {
		return "", fmt.Errorf("gagal melakukan git clone: %w", err)
	}

	return tmpDir, nil
}

func (d *Deployer) buildImage(ctx context.Context, repoPath string) (string, error) {
	tar, err := archive.TarWithOptions(repoPath, &archive.TarOptions{})
	if err != nil {
		return "", fmt.Errorf("gagal membuat tar: %v", err)
	}
	defer tar.Close()

	imageName := fmt.Sprintf("%s/%s:latest", d.config.Docker.Registry, filepath.Base(repoPath))
	buildOptions := types.ImageBuildOptions{
		Tags:       []string{imageName},
		Dockerfile: "Dockerfile",
		Remove:     true,
	}

	resp, err := d.client.ImageBuild(ctx, tar, buildOptions)
	if err != nil {
		return "", fmt.Errorf("gagal build image: %v", err)
	}
	defer resp.Body.Close()

	return imageName, nil
}

func (d *Deployer) pushImage(ctx context.Context, imageName string) error {
	authConfig := registry.AuthConfig{
		Username: d.config.Docker.Username,
		Password: d.config.Docker.Password,
	}

	_, err := d.client.ImagePush(ctx, imageName, image.PushOptions{
		RegistryAuth: authConfig.Username + ":" + authConfig.Password,
	})
	if err != nil {
		return fmt.Errorf("gagal push image: %v", err)
	}

	return nil
}

func (d *Deployer) deployToSwarm(ctx context.Context, imageName string) error {
	replicas := uint64(d.config.Deploy.Replicas)
	updateDelay, _ := time.ParseDuration(d.config.Deploy.UpdateDelay)

	serviceSpec := &swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: filepath.Base(imageName),
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: imageName,
			},
		},
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &replicas,
			},
		},
		UpdateConfig: &swarm.UpdateConfig{
			Delay:         updateDelay,
			FailureAction: d.config.Deploy.FailureAction,
		},
	}

	_, err := d.client.ServiceCreate(ctx, *serviceSpec, types.ServiceCreateOptions{})
	return err
}

func (d *Deployer) pullImage(ctx context.Context, images string) error {
	_, err := d.client.ImagePull(ctx, images, image.PullOptions{})
	if err != nil {
		return err
	}
	return nil
}

func parseCPUs(cpus string) int64 {
	if cpus == "" {
		return 0
	}

	// Konversi CPU cores ke nanoCPUs (1 CPU = 1000000000 nanoCPUs)
	// Contoh input: "0.5" untuk setengah core, "2" untuk 2 cores
	var cpu float64
	fmt.Sscanf(cpus, "%f", &cpu)
	return int64(cpu * 1e9)
}

func parseMemory(memory string) int64 {
	if memory == "" {
		return 0
	}

	// Konversi string memory ke bytes
	// Format yang didukung: "1024", "1024b", "1024k", "1024m", "1024g"
	var value float64
	unit := "b"

	if len(memory) > 0 {
		lastChar := memory[len(memory)-1]
		if lastChar == 'b' || lastChar == 'k' || lastChar == 'm' || lastChar == 'g' {
			unit = string(lastChar)
			memory = memory[:len(memory)-1]
		}
	}

	fmt.Sscanf(memory, "%f", &value)

	switch unit {
	case "k":
		value *= 1024
	case "m":
		value *= 1024 * 1024
	case "g":
		value *= 1024 * 1024 * 1024
	}

	return int64(value)
}

func parseDuration(duration string) time.Duration {
	d, _ := time.ParseDuration(duration)
	return d
}

func convertPorts(ports []deploy.PortConfig) []swarm.PortConfig {
	result := make([]swarm.PortConfig, len(ports))
	for i, p := range ports {
		result[i] = swarm.PortConfig{
			TargetPort:    p.Target,
			PublishedPort: p.Published,
			Protocol:      swarm.PortConfigProtocolTCP,
			PublishMode:   swarm.PortConfigPublishModeIngress,
		}
	}
	return result
}

func parsePortConfig(portStr string) (swarm.PortConfig, error) {
	parts := strings.Split(portStr, ":")
	if len(parts) != 2 {
		return swarm.PortConfig{}, fmt.Errorf("format port tidak valid: %s", portStr)
	}

	published, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return swarm.PortConfig{}, fmt.Errorf("port published tidak valid: %v", err)
	}

	target, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return swarm.PortConfig{}, fmt.Errorf("port target tidak valid: %v", err)
	}

	return swarm.PortConfig{
		TargetPort:    uint32(target),
		PublishedPort: uint32(published),
		Protocol:      swarm.PortConfigProtocolTCP,
		PublishMode:   swarm.PortConfigPublishModeIngress,
	}, nil
}
