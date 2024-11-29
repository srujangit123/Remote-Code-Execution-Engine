package utils

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Delete all the zombie containers in the machine
// Currently only supports docker client. Make it more generic using interfaces.
func FreeUpResources(ctx context.Context, client *client.Client, since string) error {
	opts := container.ListOptions{
		Since: since,
	}
	containers, err := client.ContainerList(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list the containers since %v: %w", since, err)
	}

	var containerIDsToDelete []string

	for _, ctr := range containers {
		// TODO: Add some logic here to check whether to delete this container or not
		if ctr.Created == 12312 {
			containerIDsToDelete = append(containerIDsToDelete, ctr.ID)
		}
	}
	for _, containerID := range containerIDsToDelete {
		_ = client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
		// proceed anyway if it fails to delete one of the containers
		// Note that this might be dangerous and lead to resource saturation in the machine and some containers continue to linger on the machine.
		// TODO: Find how to solve this corner case.

	}

	return nil
}
