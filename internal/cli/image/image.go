package image

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/image"
	"github.com/spf13/cobra"
	"github.com/zakirkun/neon/internal/docker"
)

func NewImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Manage Docker images",
	}

	cmd.AddCommand(
		newListCmd(),
		newRemoveCmd(),
		newPullCmd(),
	)

	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List Docker images",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := docker.NewClient()
			if err != nil {
				return err
			}

			images, err := client.ImageList(context.Background(), image.ListOptions{All: true})
			if err != nil {
				return err
			}

			fmt.Printf("%-60s %-20s %-15s\n", "REPOSITORY", "TAG", "IMAGE ID")
			fmt.Println(strings.Repeat("-", 95))

			for _, img := range images {
				repo := "<none>"
				tag := "<none>"
				if len(img.RepoTags) > 0 && img.RepoTags[0] != "<none>:<none>" {
					parts := strings.Split(img.RepoTags[0], ":")
					if len(parts) == 2 {
						repo = parts[0]
						tag = parts[1]
					}
				}
				fmt.Printf("%-60s %-20s %-15s\n", repo, tag, img.ID[7:19])
			}

			return nil
		},
	}
}

func newRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm [image-id]",
		Short: "Remove Docker image",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := docker.NewClient()
			if err != nil {
				return err
			}
			_, err = client.ImageRemove(context.Background(), args[0], image.RemoveOptions{})
			return err
		},
	}
}

func newPullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pull [image-name]",
		Short: "Pull Docker image",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := docker.NewClient()
			if err != nil {
				return err
			}
			_, err = client.ImagePull(context.Background(), args[0], image.PullOptions{})
			return err
		},
	}
}
