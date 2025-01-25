package docker

import (
	"fmt"

	"github.com/docker/docker/client"
)

type Client struct {
	*client.Client
}

func NewClient() (*Client, error) {
	// Set API version to match server
	client.WithVersion("1.45")

	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion("1.45"), // Pin API version
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %v", err)
	}

	return &Client{
		Client: cli,
	}, nil
}
