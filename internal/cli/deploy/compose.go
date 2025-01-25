package deploy

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zakirkun/neon/internal/config/compose"
	"github.com/zakirkun/neon/internal/docker"
)

func newComposeCmd() *cobra.Command {
	var composePath string

	cmd := &cobra.Command{
		Use:   "compose",
		Short: "Deploy dari Docker Compose file",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load compose file
			composeConfig, err := compose.LoadFromFile(composePath)
			if err != nil {
				return err
			}

			// Inisialisasi Docker client
			client, err := docker.NewClient()
			if err != nil {
				return err
			}

			ctx := context.Background()
			deployer := docker.NewDeployer(client, nil)

			// Deploy setiap service
			for name, service := range composeConfig.Services {
				fmt.Printf("Deploying service: %s\n", name)

				if err := deployer.DeployComposeService(ctx, name, &service); err != nil {
					return fmt.Errorf("gagal deploy service %s: %v", name, err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&composePath, "file", "f", "docker-compose.yml", "Path ke Docker Compose file")
	return cmd
}
