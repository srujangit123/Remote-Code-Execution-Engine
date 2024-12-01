package codecontainer

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type dockerClient struct {
	ContainerClient
	client *client.Client
	logger *zap.Logger
}

func NewDockerClient(opts *client.Opt, logger *zap.Logger) (ContainerClient, error) {
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
		logger: logger,
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
func (d *dockerClient) ExecuteCode(ctx context.Context, code *Code) (string, error) {
	err := d.createCodeFileHost(code)
	if err != nil {
		return "", fmt.Errorf("failed to create the code file: %w", err)
	}
	d.logger.Info("successfully created the code file in the host")

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

	d.logger.Info("created the container, waiting for the container to start")
	if err = d.client.ContainerStart(ctx, res.ID, container.StartOptions{}); err != nil {
		return res.ID, fmt.Errorf("failed to start the container after creating: %w", err)
	}

	d.logger.Info("container started, waiting for the container to exit")
	statusCh, errCh := d.client.ContainerWait(ctx, res.ID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			return "", fmt.Errorf("failed to get the container logs: %w", err)
		}
	case status := <-statusCh:
		d.logger.Info("container exited", zap.Any("status", status))
	}

	logs, err := d.client.ContainerLogs(ctx, res.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
	})
	if err != nil {
		return "", fmt.Errorf("failed to read container logs: %w", err)
	}
	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, logs)
	if err != nil {
		return "", fmt.Errorf("error processing the logs: %w", err)
	}

	return stdoutBuf.String() + "\n" + stderrBuf.String(), nil
}

func (d *dockerClient) FreeUpZombieContainers(ctx context.Context) error {
	for {
		pruneResults, err := d.client.ContainersPrune(ctx, filters.Args{})
		if err != nil {
			d.logger.Error("failed to prune containers",
				zap.Error(err),
			)
		}

		d.logger.Info("successfully pruned the containers:",
			zap.Int("#Pruned containers", len(pruneResults.ContainersDeleted)),
		)

		time.Sleep(5 * time.Minute)
	}
}

func (d *dockerClient) createCodeFileHost(code *Code) error {
	codeFilePath := getCodeFilePathHost(code)
	f, err := os.Create(codeFilePath)
	if err != nil {
		panic(err)
	}

	data, err := base64.StdEncoding.DecodeString(code.EncodedCode)
	if err != nil {
		return fmt.Errorf("failed to decode the code text: %w", err)
	}

	n, err := f.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("failed to write the content to the file: %w", err)
	}
	d.logger.Info("wrote the code content to the file",
		zap.String("file path", codeFilePath),
		zap.Int("bytes", n),
	)

	return nil
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
		return GolangContainerImage
	case "cpp":
		return cppContainerImage

	default:
		return ""
	}
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

func getCodeFilePathHost(code *Code) string {
	var fileExtension string
	var codeDirectoryPath string
	switch code.Language {
	case "cpp":
		fileExtension = "cpp"
		codeDirectoryPath = CppCodePath
	case "golang":
		fileExtension = "go"
		codeDirectoryPath = GolangCodePath
	}
	return filepath.Join(codeDirectoryPath, code.FileName+"."+fileExtension)
}

func getLanguageRunCmd(code *Code) []string {
	codeFilePath := getCodeFilePath(code)

	return []string{
		"sh", "-c",
		fmt.Sprintf("/usr/bin/run-code.sh %s %s", code.Language, codeFilePath),
	}
}
