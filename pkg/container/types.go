package codecontainer

import (
	"context"

	"github.com/docker/docker/api/types/container"
)

type Container struct {
	Image  string
	ID     string
	Status string
}

type Code struct {
	EncodedCode string
	FileName    string
	Language    string
}

type ContainerClient interface {
	FreeUpZombieContainers(ctx context.Context) error
	ExecuteCode(ctx context.Context, code *Code) (string, error) // Executes and returns the output in the string, error in case of server errors not code errors.
	// TODO: Is this even needed?
	GetContainers(ctx context.Context, opts *container.ListOptions) ([]Container, error) // Remove list options if you want some other container type other than docker
}
