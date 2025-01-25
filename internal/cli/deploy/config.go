package deploy

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zakirkun/neon/internal/config/deploy"
	"github.com/zakirkun/neon/internal/docker"
	"gopkg.in/yaml.v3"
)

func newConfigDeployCmd() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Deploy services dari file konfigurasi",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := os.ReadFile(configFile)
			if err != nil {
				return fmt.Errorf("gagal membaca file konfigurasi: %v", err)
			}

			var config deploy.Config
			if err := yaml.Unmarshal(data, &config); err != nil {
				return fmt.Errorf("gagal parse konfigurasi: %v", err)
			}

			client, err := docker.NewClient()
			if err != nil {
				return err
			}

			ctx := context.Background()
			deployer := docker.NewDeployer(client, nil)

			for _, svc := range config.Services {
				if err := deployer.DeployFromConfig(ctx, &svc); err != nil {
					return fmt.Errorf("gagal deploy service %s: %v", svc.Name, err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&configFile, "file", "f", "config/deploy.yaml", "Path ke file konfigurasi deploy")
	return cmd
}
