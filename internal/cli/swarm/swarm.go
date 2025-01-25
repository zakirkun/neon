package swarm

import "github.com/spf13/cobra"

func NewSwarmCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "swarm",
		Short: "Manajemen Docker Swarm cluster",
	}
}
