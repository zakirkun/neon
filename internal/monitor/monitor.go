package monitor

import (
	"context"
	"encoding/json"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/zakirkun/neon/internal/docker"
	"github.com/zakirkun/neon/internal/logger"
)

type Metrics struct {
	CPU    float64
	Memory float64
	IO     float64
}

type Monitor struct {
	client *docker.Client
}

func NewMonitor(client *docker.Client) *Monitor {
	return &Monitor{client: client}
}

func (m *Monitor) WatchService(ctx context.Context, serviceID string) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			tasks, err := m.client.TaskList(ctx, types.TaskListOptions{
				Filters: filters.NewArgs(filters.Arg("service", serviceID)),
			})
			if err != nil {
				logger.Error(err, "Failed to get tasks")
				continue
			}

			for _, task := range tasks {
				if task.Status.State != swarm.TaskStateRunning {
					logger.Warnf("Task %s is in state: %s", task.ID[:12], task.Status.State)
					continue
				}

				stats, err := m.getContainerStats(ctx, task.Status.ContainerStatus.ContainerID)
				if err != nil {
					logger.Error(err, "Failed to get container stats")
					continue
				}

				logger.Infof("Task %s - CPU: %.2f%%, Memory: %.2f MB",
					task.ID[:12], stats.CPU, stats.Memory/1024/1024)
			}
		}
	}
}

func (m *Monitor) getContainerStats(ctx context.Context, containerID string) (*Metrics, error) {
	stats, err := m.client.ContainerStats(ctx, containerID, false)
	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()

	var metrics Metrics
	decoder := json.NewDecoder(stats.Body)
	if err := decoder.Decode(&metrics); err != nil {
		return nil, err
	}

	return &metrics, nil
}
