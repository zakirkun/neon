package container

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zakirkun/neon/internal/docker"
	"github.com/zakirkun/neon/internal/docker/container"
)

func NewContainerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "container",
		Short: "Manajemen Docker container",
	}

	cmd.AddCommand(
		newListCmd(),
		newStartCmd(),
		newStopCmd(),
		newRemoveCmd(),
		newLogsCmd(),
	)

	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Menampilkan daftar container",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := docker.NewClient()
			if err != nil {
				return err
			}

			manager := container.NewManager(client)
			containers, err := manager.List(context.Background())
			if err != nil {
				return err
			}

			for _, c := range containers {
				fmt.Printf("ID: %s, Name: %s, Status: %s\n", c.ID[:12], c.Names[0], c.Status)
			}

			return nil
		},
	}
}

// Implementasi command lainnya...
