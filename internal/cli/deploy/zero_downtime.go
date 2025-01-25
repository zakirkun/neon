package deploy

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/swarm"
	"github.com/spf13/cobra"
	"github.com/zakirkun/neon/internal/docker"
	"github.com/zakirkun/neon/internal/logger"
)

func newZeroDowntimeCmd() *cobra.Command {
	var (
		replicas    uint64
		updateDelay time.Duration
		image       string
	)

	cmd := &cobra.Command{
		Use:   "rolling [service-name]",
		Short: "Deploy service with zero downtime using rolling update",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := docker.NewClient()
			if err != nil {
				return err
			}

			// Get existing service
			service, _, err := client.ServiceInspectWithRaw(context.Background(), args[0], types.ServiceInspectOptions{})
			if err != nil {
				return fmt.Errorf("service not found: %v", err)
			}

			// Prepare update config
			updateConfig := &swarm.UpdateConfig{
				Parallelism:   1,
				Delay:         time.Duration(updateDelay),
				Order:         "start-first", // Start new container before stopping old one
				FailureAction: "rollback",
				Monitor:       time.Duration(5 * time.Second),
			}

			// Prepare rollback config
			rollbackConfig := &swarm.UpdateConfig{
				Parallelism:   1,
				Delay:         time.Duration(updateDelay),
				Order:         "stop-first",
				FailureAction: "pause",
			}

			// Update service spec
			service.Spec.Mode.Replicated.Replicas = &replicas
			service.Spec.TaskTemplate.ContainerSpec.Image = image
			service.Spec.UpdateConfig = updateConfig
			service.Spec.RollbackConfig = rollbackConfig

			// Health check config
			service.Spec.TaskTemplate.ContainerSpec.Healthcheck = &container.HealthConfig{
				Test:     []string{"CMD-SHELL", "curl -f http://localhost/health || exit 1"},
				Interval: time.Duration(5 * time.Second),
				Timeout:  time.Duration(3 * time.Second),
				Retries:  3,
			}

			logger.Info("Starting zero-downtime deployment...")

			// Update service
			response, err := client.ServiceUpdate(
				context.Background(),
				args[0],
				service.Version,
				service.Spec,
				types.ServiceUpdateOptions{},
			)
			if err != nil {
				return fmt.Errorf("deployment failed: %v", err)
			}

			if len(response.Warnings) > 0 {
				logger.Warn("Deployment warnings: " + strings.Join(response.Warnings, ", "))
			}

			logger.Info("Zero-downtime deployment completed successfully")
			return nil
		},
	}

	cmd.Flags().Uint64VarP(&replicas, "replicas", "r", 3, "Number of replicas")
	cmd.Flags().DurationVarP(&updateDelay, "update-delay", "d", 10*time.Second, "Delay between updates")
	cmd.Flags().StringVarP(&image, "image", "i", "", "New image to deploy")
	cmd.MarkFlagRequired("image")

	return cmd
}
