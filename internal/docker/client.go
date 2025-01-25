package docker

import (
	"github.com/docker/docker/client"
)

type Client struct {
	*client.Client
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &Client{cli}, nil
}
