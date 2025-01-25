package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/zakirkun/neon/internal/docker"
)

type Manager struct {
	client *docker.Client
}

func NewManager(client *docker.Client) *Manager {
	return &Manager{client: client}
}

func (m *Manager) List(ctx context.Context) ([]types.Container, error) {
	return m.client.ContainerList(ctx, types.ContainerListOptions{All: true})
}

func (m *Manager) Start(ctx context.Context, containerID string) error {
	return m.client.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}

func (m *Manager) Stop(ctx context.Context, containerID string) error {
	return m.client.ContainerStop(ctx, containerID, container.StopOptions{})
}

func (m *Manager) Remove(ctx context.Context, containerID string, force bool) error {
	return m.client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: force})
}

func (m *Manager) Logs(ctx context.Context, containerID string) (string, error) {
	logs, err := m.client.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", err
	}
	defer logs.Close()

	// Convert logs to string
	buf := new([]byte)
	_, err = logs.Read(*buf)
	if err != nil {
		return "", err
	}

	return string(*buf), nil
}
