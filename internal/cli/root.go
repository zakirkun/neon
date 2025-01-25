package cli

import (
	"github.com/spf13/cobra"
	"github.com/zakirkun/neon/internal/cli/autoscale"
	"github.com/zakirkun/neon/internal/cli/container"
	"github.com/zakirkun/neon/internal/cli/deploy"
	"github.com/zakirkun/neon/internal/cli/swarm"
)

var rootCmd = &cobra.Command{
	Use:   "neon",
	Short: "Neon - DevOps Management Tool",
	Long: `Neon adalah tools untuk membantu manajemen DevOps seperti
deployment, scaling, dan pengelolaan Docker/Docker Swarm.`,
}

func init() {
	rootCmd.AddCommand(
		deploy.NewDeployCmd(),
		container.NewContainerCmd(),
		swarm.NewSwarmCmd(),
		autoscale.NewAutoscaleCmd(),
	)
}

func Execute() error {
	return rootCmd.Execute()
}
