package volume

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/spf13/cobra"
	"github.com/zakirkun/neon/internal/docker"
)

func NewVolumeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "Manage Docker volumes",
	}

	cmd.AddCommand(
		newListCmd(),
		newCreateCmd(),
		newRemoveCmd(),
	)

	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List Docker volumes",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := docker.NewClient()
			if err != nil {
				return err
			}

			volumes, err := client.VolumeList(context.Background(), volume.ListOptions{
				Filters: filters.NewArgs(),
			})
			if err != nil {
				return err
			}

			fmt.Printf("%-40s %-20s %-60s\n", "NAME", "DRIVER", "MOUNTPOINT")
			fmt.Println(strings.Repeat("-", 120))

			for _, vol := range volumes.Volumes {
				fmt.Printf("%-40s %-20s %-60s\n", vol.Name, vol.Driver, vol.Mountpoint)
			}

			return nil
		},
	}
}

func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create [name]",
		Short: "Create Docker volume",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := docker.NewClient()
			if err != nil {
				return err
			}
			_, err = client.VolumeCreate(context.Background(), volume.CreateOptions{Name: args[0]})
			return err
		},
	}
}

func newRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm [volume-name]",
		Short: "Remove Docker volume",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := docker.NewClient()
			if err != nil {
				return err
			}
			return client.VolumeRemove(context.Background(), args[0], false)
		},
	}
}
