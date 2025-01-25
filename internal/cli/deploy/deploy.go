package deploy

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zakirkun/neon/internal/config"
	"github.com/zakirkun/neon/internal/docker"
)

var (
	configPath string
	repoURL    string
)

func NewDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy [flags] [repo-url]",
		Short: "Deploy aplikasi ke Docker Swarm",
		Long: `Deploy aplikasi dari repository GitHub ke Docker Swarm cluster.
Contoh: neon deploy https://github.com/user/repo`,
		RunE: runDeploy,
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "config/config.yaml", "Path ke file konfigurasi")
	cmd.Flags().StringVarP(&repoURL, "repo", "r", "", "URL repository GitHub")

	return cmd
}

func runDeploy(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("gagal memuat konfigurasi: %v", err)
	}

	// Inisialisasi Docker client
	client, err := docker.NewClient()
	if err != nil {
		return fmt.Errorf("gagal menginisialisasi Docker client: %v", err)
	}

	ctx := context.Background()

	// Proses deployment
	deployer := docker.NewDeployer(client, cfg)
	err = deployer.Deploy(ctx, repoURL)
	if err != nil {
		return fmt.Errorf("gagal melakukan deployment: %v", err)
	}

	return nil
}
