package container

import (
	"context"
	"io"
	"os"

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
	return m.client.ContainerList(ctx, container.ListOptions{All: true})
}

func (m *Manager) Start(ctx context.Context, containerID string) error {
	return m.client.ContainerStart(ctx, containerID, container.StartOptions{})
}

func (m *Manager) Stop(ctx context.Context, containerID string) error {
	return m.client.ContainerStop(ctx, containerID, container.StopOptions{})
}

func (m *Manager) Remove(ctx context.Context, containerID string, force bool) error {
	return m.client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: force})
}

func (m *Manager) Logs(ctx context.Context, containerID string) error {
	logs, err := m.client.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return err
	}
	defer logs.Close()

	_, err = io.Copy(os.Stdout, logs)
	return err
}
