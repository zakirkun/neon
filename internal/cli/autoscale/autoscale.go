package autoscale

import "github.com/spf13/cobra"

func NewAutoscaleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "autoscale",
		Short: "Manajemen auto scaling untuk services",
	}
}
