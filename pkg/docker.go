package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type Container struct {
	Image  string
	ID     string
	Status string
}

type Code struct {
	FileName string
	Language string
}

type Volume map[string]struct{}

func GetContainers(ctx context.Context, client *client.Client, opts *container.ListOptions) ([]Container, error) {
	var containersList []Container

	containers, err := client.ContainerList(ctx, *opts)
	if err != nil {
		return containersList, fmt.Errorf("failed to get the list of containers: %w", err)
	}

	for _, ctr := range containers {
		container := Container{
			Image:  ctr.Image,
			ID:     ctr.ID,
			Status: ctr.Status,
		}

		containersList = append(containersList, container)
	}

	return containersList, nil
}

func getLanguageContainerImage(lang string) string {
	switch lang {
	case "golang":
		return GolangCodePath
	case "cpp":
		return cppContainerImage

	default:
		return ""
	}

}

func getLanguageCodePath(lang string) string {
	basePath := "/home/srujan/Documents/code/"
	return basePath + lang
}

func getLanguageRunCmd(code *Code) []string {
	_ = code
	return []string{
		"sh", "-c",
		fmt.Sprintf("ls /container/code/; cat /container/code/%s; g++ /container/code/%s -o /container/code/1; ./container/code/1 > /container/code/%s.out", code.FileName, code.FileName, code.FileName),
	}
}

// Different docker images for different languages. Returns container ID and error if any.
// How to create a container?
// 1. Use volume mounts. Mount a specific file - the upload file. Then run the container with a custom command.
// TODO: Add Cgroups support
// TODO: Improve security: drop all the capabilities and use only those that are absolutely necessary.
func CreateContainer(ctx context.Context, client *client.Client, code *Code) (string, error) {
	res, err := client.ContainerCreate(ctx, &container.Config{
		Cmd:   getLanguageRunCmd(code),
		Image: getLanguageContainerImage(code.Language),
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: getLanguageCodePath(code.Language),
				Target: "/container/code",
			},
		},
	}, nil, nil, "sample-container-name-123123123")

	if err != nil {
		return "", fmt.Errorf("failed to create a container: %w", err)
	}

	if err = client.ContainerStart(ctx, res.ID, container.StartOptions{}); err != nil {
		return res.ID, fmt.Errorf("failed to start the container after creating: %w", err)
	}

	return res.ID, nil
}

func GetCodeOutput(ctx context.Context, client *client.Client, code *Code) (string, error) {
	f, err := os.Open("/home/srujan/Documents/code/cpp/" + code.FileName + ".out")
	if err != nil {
		return "", fmt.Errorf("error while opening the file: %w", err)
	}
	defer f.Close()
	fileContent, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("error reading content from the file: %w", err)
	}
	return string(fileContent), nil
}

// Delete all the zombie containers in the machine
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

func CreateContainerAndRun(ctx context.Context, client *client.Client, code Code) {
	// 1. Create a container with the docker image of a language
	// 2. Just mount the volume. For now we just upload the code to the machine.
	// 3. Execute the container.
}
