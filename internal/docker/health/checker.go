package health

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/client"
	"github.com/zakirkun/neon/internal/logger"
)

type Checker struct {
	client *client.Client
}

func NewChecker(client *client.Client) *Checker {
	return &Checker{
		client: client,
	}
}

func (c *Checker) CheckDockerStatus() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if Docker daemon is running
	_, err := c.client.Ping(ctx)
	if err != nil {
		logger.Error(err, "Docker daemon is not running")
		return fmt.Errorf("docker daemon is not running: %v", err)
	}

	logger.Info("Docker daemon is running")
	return nil
}

func (c *Checker) CheckSwarmStatus() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Check Swarm status
	info, err := c.client.Info(ctx)
	if err != nil {
		logger.Error(err, "Failed to get Docker info")
		return fmt.Errorf("failed to get docker info: %v", err)
	}

	if !info.Swarm.ControlAvailable {
		logger.Warn("Docker Swarm mode is not enabled")
		return fmt.Errorf("docker swarm mode is not enabled")
	}

	logger.Info("Docker Swarm mode is active")
	return nil
}
