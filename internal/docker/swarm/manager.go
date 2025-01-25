package swarm

import (
	"context"

	"github.com/docker/docker/api/types"
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
