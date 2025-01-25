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

func newStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start [container-id]",
		Short: "Menjalankan container",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("container ID diperlukan")
			}

			client, err := docker.NewClient()
			if err != nil {
				return err
			}

			manager := container.NewManager(client)
			return manager.Start(context.Background(), args[0])
		},
	}
}

func newStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop [container-id]",
		Short: "Menghentikan container",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("container ID diperlukan")
			}

			client, err := docker.NewClient()
			if err != nil {
				return err
			}

			manager := container.NewManager(client)
			return manager.Stop(context.Background(), args[0])
		},
	}
}

func newRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm [container-id]",
		Short: "Menghapus container",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("container ID diperlukan")
			}

			client, err := docker.NewClient()
			if err != nil {
				return err
			}

			manager := container.NewManager(client)
			return manager.Remove(context.Background(), args[0], true)
		},
	}
}

func newLogsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logs [container-id]",
		Short: "Menampilkan logs container",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("container ID diperlukan")
			}

			client, err := docker.NewClient()
			if err != nil {
				return err
			}

			manager := container.NewManager(client)
			return manager.Logs(context.Background(), args[0])
		},
	}
}

// Implementasi command lainnya...
