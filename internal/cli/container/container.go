package container

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	dContainer "github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
	"github.com/zakirkun/neon/internal/docker"
	"github.com/zakirkun/neon/internal/docker/container"
)

func NewContainerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "container",
		Short: "Manage Docker containers",
	}

	cmd.AddCommand(
		newListCmd(),
		newStartCmd(),
		newStopCmd(),
		newRemoveCmd(),
		newLogsCmd(),
		newInspectCmd(),
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
	var tail int
	var follow bool

	cmd := &cobra.Command{
		Use:   "logs [container-id]",
		Short: "Fetch the logs of a container",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := docker.NewClient()
			if err != nil {
				return err
			}

			options := dContainer.LogsOptions{
				ShowStdout: true,
				ShowStderr: true,
				Follow:     follow,
				Tail:       strconv.Itoa(tail),
			}

			logs, err := client.ContainerLogs(context.Background(), args[0], options)
			if err != nil {
				return err
			}
			defer logs.Close()

			_, err = io.Copy(os.Stdout, logs)
			return err
		},
	}

	cmd.Flags().IntVarP(&tail, "tail", "n", 100, "Number of lines to show from the end of the logs")
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")

	return cmd
}

func newInspectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "inspect [container-id]",
		Short: "Display detailed information on a container",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := docker.NewClient()
			if err != nil {
				return err
			}

			container, err := client.ContainerInspect(context.Background(), args[0])
			if err != nil {
				return err
			}

			data, err := json.MarshalIndent(container, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(data))
			return nil
		},
	}
}

// Implementasi command lainnya...
