package network

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
	"github.com/zakirkun/neon/internal/docker"
)

func NewNetworkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Manage Docker networks",
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
		Short: "List Docker networks",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := docker.NewClient()
			if err != nil {
				return err
			}

			networks, err := client.NetworkList(context.Background(), types.NetworkListOptions{})
			if err != nil {
				return err
			}

			fmt.Printf("%-40s %-15s %-15s %-20s\n", "NAME", "DRIVER", "SCOPE", "SUBNET")
			fmt.Println(strings.Repeat("-", 90))

			for _, net := range networks {
				subnet := ""
				if len(net.IPAM.Config) > 0 {
					subnet = net.IPAM.Config[0].Subnet
				}
				fmt.Printf("%-40s %-15s %-15s %-20s\n", net.Name, net.Driver, net.Scope, subnet)
			}

			return nil
		},
	}
}

func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create [name]",
		Short: "Create Docker network",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := docker.NewClient()
			if err != nil {
				return err
			}
			_, err = client.NetworkCreate(context.Background(), args[0], types.NetworkCreate{})
			return err
		},
	}
}

func newRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm [network-id]",
		Short: "Remove Docker network",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := docker.NewClient()
			if err != nil {
				return err
			}
			return client.NetworkRemove(context.Background(), args[0])
		},
	}
}
