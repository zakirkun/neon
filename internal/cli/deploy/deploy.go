package deploy

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/zakirkun/neon/internal/config"
	"github.com/zakirkun/neon/internal/docker"
)

var (
	configPath string
	repoURL    string
)

type DeployConfig struct {
	RepoURL     string `yaml:"repo_url"`
	Branch      string `yaml:"branch"`
	ConfigFile  string `yaml:"config_file"`
	Environment string `yaml:"environment"`
}

func NewDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy applications to Docker Swarm",
	}

	cmd.AddCommand(
		newZeroDowntimeCmd(),
		newConfigDeployCmd(),
		newComposeCmd(),
	)

	cmd.Flags().StringVarP(&configPath, "config", "c", "config/config.yaml", "Path ke file konfigurasi")
	cmd.Flags().StringVarP(&repoURL, "repo", "r", "", "URL repository GitHub")

	return cmd
}

func newConfigCmd() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Deploy using configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deployFromConfig(configFile)
		},
	}

	cmd.Flags().StringVarP(&configFile, "file", "f", "deploy.yaml", "Deployment configuration file")
	return cmd
}

// func newComposeCmd() *cobra.Command {
// 	var composeFile string

// 	cmd := &cobra.Command{
// 		Use:   "compose",
// 		Short: "Deploy using docker-compose file",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			return deployFromCompose(composeFile)
// 		},
// 	}

// 	cmd.Flags().StringVarP(&composeFile, "file", "f", "docker-compose.yml", "Docker compose file")
// 	return cmd
// }

func deployFromConfig(configFile string) error {
	// TODO: Implement config-based deployment
	return nil
}

// func deployFromCompose(composeFile string) error {
// 	// TODO: Implement compose-based deployment
// 	return nil
// }

func runDeploy(cmd *cobra.Command, args []string) error {
	// Validasi input
	if repoURL == "" {
		return fmt.Errorf("URL repository harus diisi")
	}

	// Load konfigurasi
	cfg := config.Get()

	// Inisialisasi Docker client
	client, err := docker.NewClient()
	if err != nil {
		return fmt.Errorf("gagal menginisialisasi Docker client: %v", err)
	}

	// Buat context dengan timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Proses deployment
	deployer := docker.NewDeployer(client, cfg)

	fmt.Printf("Memulai deployment dari repository: %s\n", repoURL)

	err = deployer.Deploy(ctx, repoURL)
	if err != nil {
		return fmt.Errorf("gagal melakukan deployment: %v", err)
	}

	fmt.Println("Deployment berhasil dilakukan")
	return nil
}
