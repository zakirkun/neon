package docker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/pkg/archive"
	"github.com/zakirkun/neon/internal/config"
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

func (d *Deployer) cloneRepo(repoURL string) (string, error) {
	// TODO: Implementasi Git clone
	// Untuk sementara kita return dummy path
	return "/tmp/repo", nil
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
