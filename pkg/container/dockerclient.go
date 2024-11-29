package codecontainer

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

type dockerClient struct {
	ContainerClient
	client *client.Client
}

func NewDockerClient(opts *client.Opt) (ContainerClient, error) {
	var cli *client.Client
	var err error

	if opts != nil {
		cli, err = client.NewClientWithOpts(*opts)
	} else {
		// TODO: Make sure you set this env when the server starts up.
		cli, err = client.NewClientWithOpts(client.FromEnv)
	}
	if err != nil {
		return &dockerClient{}, fmt.Errorf("failed to initiliaze docker client: %w", err)
	}

	return &dockerClient{
		client: cli,
	}, nil
}

func (d *dockerClient) GetContainers(ctx context.Context, opts *container.ListOptions) ([]Container, error) {
	var containersList []Container

	containers, err := d.client.ContainerList(ctx, *opts)
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

// 1. Use volume mounts. Mount a specific file - the upload file. Then run the container with a custom command.
// TODO: Add Cgroups support
// TODO: Improve security: drop all the capabilities and use only those that are absolutely necessary.
func (d *dockerClient) CreateAndStartContainer(ctx context.Context, code *Code) (string, error) {
	res, err := d.client.ContainerCreate(ctx, &container.Config{
		Cmd:   getLanguageRunCmd(code),
		Image: getLanguageContainerImage(code.Language),
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: getHostLanguageCodePath(code.Language),
				Target: TargetMountPath,
			},
		},
	}, nil, nil, getContainerName(code))

	if err != nil {
		return "", fmt.Errorf("failed to create a container: %w", err)
	}

	if err = d.client.ContainerStart(ctx, res.ID, container.StartOptions{}); err != nil {
		return res.ID, fmt.Errorf("failed to start the container after creating: %w", err)
	}

	return res.ID, nil
}

// TODO: Implementation. Just delete all the containers that are in deleted/exited state to avoid memory saturation issues in the host.
func (d *dockerClient) FreeUpZombieContainers(ctx context.Context) error {
	return nil
}

// TODO: Need to stop the user from printing infinite times - use cgroup?
func (d *dockerClient) GetContainerOutput(ctx context.Context, code *Code) (string, error) {
	f, err := os.Open(getCodeFilePath(code) + ".out")
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

func getHostLanguageCodePath(lang string) string {
	switch lang {
	case "cpp":
		return CppCodePath
	case "golang":
		return GolangCodePath
	}
	return ""
}

func getContainerName(code *Code) string {
	return code.FileName + "-" + uuid.New().String()
}

// Only supports Golang and CPP for now. TODO: Add python, java, etc.
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

func getExecutablePath(code *Code) string {
	// Assumption is that file name doesn't contain extension name
	return filepath.Join(BaseContainerCodeExecutablePath, code.FileName)
}

func getCodeFilePath(code *Code) string {
	var fileExtension string
	switch code.Language {
	case "cpp":
		fileExtension = "cpp"
	case "golang":
		fileExtension = "go"
	}
	return filepath.Join(TargetMountPath, code.FileName+"."+fileExtension)
}

func getCodeCompilationCmd(code *Code) string {
	switch code.Language {
	case "cpp":
		return fmt.Sprintf("g++ %s", getCodeFilePath(code))
	case "golang":
		return fmt.Sprintf("go build -o %s %s", getExecutablePath(code), getCodeFilePath(code))
	}
	return ""
}

func getLanguageRunCmd(code *Code) []string {
	codeFilePath := getCodeFilePath(code)
	executablePath := getExecutablePath(code)
	codeCompilationCmd := getCodeCompilationCmd(code)

	return []string{
		"sh", "-c",
		fmt.Sprintf("%s; .%s > %s.out", codeCompilationCmd, executablePath, codeFilePath),
	}
}