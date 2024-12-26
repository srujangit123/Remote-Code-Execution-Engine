package codecontainer

import (
	"context"
	"remote-code-engine/pkg/config"

	"github.com/docker/docker/api/types/container"
)

type Container struct {
	Image  string
	ID     string
	Status string
}

type Code struct {
	EncodedCode  string
	EncodedInput string
	Language     config.Language
	config.LanguageConfig
}

type ContainerClient interface {
	FreeUpZombieContainers(ctx context.Context) error

	// Executes and returns the output in the string, error in case of server errors not code errors.
	ExecuteCode(ctx context.Context, code *Code) (string, error)

	// TODO: Is this even needed?
	// Remove list options if you want some other container type other than docker
	GetContainers(ctx context.Context, opts *container.ListOptions) ([]Container, error)
}
