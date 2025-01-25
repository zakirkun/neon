package autoscale

import (
	"context"
	"time"

	"github.com/zakirkun/neon/internal/docker"
	"github.com/zakirkun/neon/internal/docker/swarm"
)

type ScalingRule struct {
	ServiceID     string
	MinReplicas   uint64
	MaxReplicas   uint64
	CPUThreshold  float64
	ScaleUpStep   uint64
	ScaleDownStep uint64
	Cooldown      time.Duration
}

type Manager struct {
	client       *docker.Client
	swarmManager *swarm.Manager
	rules        map[string]ScalingRule
	lastScale    map[string]time.Time
}

func NewManager(client *docker.Client) *Manager {
	return &Manager{
		client:       client,
		swarmManager: swarm.NewManager(client),
		rules:        make(map[string]ScalingRule),
		lastScale:    make(map[string]time.Time),
	}
}

func (m *Manager) AddRule(rule ScalingRule) {
	m.rules[rule.ServiceID] = rule
}

func (m *Manager) StartMonitoring(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			m.checkAndScale(ctx)
		}
	}
}

func (m *Manager) checkAndScale(ctx context.Context) {
	for serviceID, rule := range m.rules {
		// Skip if in cooldown period
		if lastScale, ok := m.lastScale[serviceID]; ok {
			if time.Since(lastScale) < rule.Cooldown {
				continue
			}
		}

		stats, err := m.getServiceStats(ctx, serviceID)
		if err != nil {
			continue
		}

		currentReplicas := stats.replicas
		cpuUsage := stats.cpuUsage

		if cpuUsage > rule.CPUThreshold && currentReplicas < rule.MaxReplicas {
			// Scale up
			newReplicas := min(currentReplicas+rule.ScaleUpStep, rule.MaxReplicas)
			err := m.swarmManager.ScaleService(ctx, serviceID, newReplicas)
			if err == nil {
				m.lastScale[serviceID] = time.Now()
			}
		} else if cpuUsage < rule.CPUThreshold/2 && currentReplicas > rule.MinReplicas {
			// Scale down
			newReplicas := max(currentReplicas-rule.ScaleDownStep, rule.MinReplicas)
			err := m.swarmManager.ScaleService(ctx, serviceID, newReplicas)
			if err == nil {
				m.lastScale[serviceID] = time.Now()
			}
		}
	}
}

type serviceStats struct {
	replicas uint64
	cpuUsage float64
}

func (m *Manager) getServiceStats(ctx context.Context, serviceID string) (*serviceStats, error) {
	// Implementasi pengambilan statistik service
	// TODO: Implementasi detail monitoring metrics
	return &serviceStats{}, nil
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}
