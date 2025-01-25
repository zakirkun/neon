package backup

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/zakirkun/neon/internal/docker"
	"github.com/zakirkun/neon/internal/logger"
)

type Manager struct {
	client    *docker.Client
	backupDir string
}

func NewManager(client *docker.Client, backupDir string) *Manager {
	return &Manager{
		client:    client,
		backupDir: backupDir,
	}
}

func (m *Manager) BackupVolume(ctx context.Context, volumeName string) error {
	// Create backup container
	resp, err := m.client.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"tar", "czf", "/backup/data.tar.gz", "/data"},
		Volumes: map[string]struct{}{
			"/data":   {},
			"/backup": {},
		},
	}, &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/data", volumeName),
			fmt.Sprintf("%s:/backup", m.backupDir),
		},
	}, nil, nil, "")
	if err != nil {
		return fmt.Errorf("failed to create backup container: %v", err)
	}

	// Start backup process
	if err := m.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start backup container: %v", err)
	}

	logger.Infof("Backup of volume %s completed", volumeName)
	return nil
}
