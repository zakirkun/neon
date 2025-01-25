package swarm

import (
	"context"
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/zakirkun/neon/internal/docker"
)

type Manager struct {
	client *docker.Client
}

func NewManager(client *docker.Client) *Manager {
	return &Manager{client: client}
}

func (m *Manager) InitSwarm(ctx context.Context, advertiseAddr string) (string, error) {
	req := swarm.InitRequest{
		ListenAddr:    "0.0.0.0:2377",
		AdvertiseAddr: advertiseAddr,
	}
	return m.client.SwarmInit(ctx, req)
}

func (m *Manager) JoinSwarm(ctx context.Context, remoteAddr, token string) error {
	req := swarm.JoinRequest{
		ListenAddr:    "0.0.0.0:2377",
		AdvertiseAddr: "eth0",
		RemoteAddrs:   []string{remoteAddr},
		JoinToken:     token,
	}
	return m.client.SwarmJoin(ctx, req)
}

func (m *Manager) LeaveSwarm(ctx context.Context, force bool) error {
	return m.client.SwarmLeave(ctx, force)
}

func (m *Manager) ListNodes(ctx context.Context) ([]swarm.Node, error) {
	return m.client.NodeList(ctx, types.NodeListOptions{})
}

func (m *Manager) ListServices(ctx context.Context) ([]swarm.Service, error) {
	return m.client.ServiceList(ctx, types.ServiceListOptions{})
}

func (m *Manager) ScaleService(ctx context.Context, serviceID string, replicas uint64) error {
	service, _, err := m.client.ServiceInspectWithRaw(ctx, serviceID, types.ServiceInspectOptions{})
	if err != nil {
		return err
	}

	service.Spec.Mode.Replicated.Replicas = &replicas

	_, err = m.client.ServiceUpdate(ctx, serviceID, service.Version, service.Spec, types.ServiceUpdateOptions{})
	return err
}

type ServiceStats struct {
	Replicas uint64
	CPUUsage float64
}

func (m *Manager) GetServiceStats(ctx context.Context, serviceID string) (*ServiceStats, error) {
	service, _, err := m.client.ServiceInspectWithRaw(ctx, serviceID, types.ServiceInspectOptions{})
	if err != nil {
		return nil, err
	}

	// TODO: Implementasi pengambilan CPU usage dari container stats
	return &ServiceStats{
		Replicas: getReplicaCount(&service),
		CPUUsage: calculateCPUUsage(ctx, m.client, serviceID), // Implementasi pengambilan CPU usage dari container stats
	}, nil
}

func getReplicaCount(service *swarm.Service) uint64 {
	if service.Spec.Mode.Replicated != nil && service.Spec.Mode.Replicated.Replicas != nil {
		return *service.Spec.Mode.Replicated.Replicas
	}
	return 0
}

func calculateCPUUsage(ctx context.Context, client *docker.Client, serviceID string) float64 {
	// Dapatkan daftar task untuk service
	tasks, err := client.TaskList(ctx, types.TaskListOptions{
		Filters: filters.NewArgs(filters.Arg("service", serviceID)),
	})
	if err != nil {
		return 0.0
	}

	var totalCPUUsage float64
	var activeContainers int

	// Hitung rata-rata CPU usage dari semua container yang aktif
	for _, task := range tasks {
		if task.Status.State == swarm.TaskStateRunning {
			stats, err := client.ContainerStats(ctx, task.Status.ContainerStatus.ContainerID, false)
			if err != nil {
				continue
			}
			defer stats.Body.Close()

			var statsJSON types.StatsJSON
			if err := json.NewDecoder(stats.Body).Decode(&statsJSON); err != nil {
				continue
			}

			cpuDelta := float64(statsJSON.CPUStats.CPUUsage.TotalUsage - statsJSON.PreCPUStats.CPUUsage.TotalUsage)
			systemDelta := float64(statsJSON.CPUStats.SystemUsage - statsJSON.PreCPUStats.SystemUsage)

			if systemDelta > 0.0 && cpuDelta > 0.0 {
				cpuPercent := (cpuDelta / systemDelta) * float64(len(statsJSON.CPUStats.CPUUsage.PercpuUsage)) * 100.0
				totalCPUUsage += cpuPercent
				activeContainers++
			}
		}
	}

	if activeContainers > 0 {
		return totalCPUUsage / float64(activeContainers)
	}
	return 0.0
}
