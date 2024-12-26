package codecontainer

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"remote-code-engine/pkg/config"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-units"
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
	codeFileName, inputFileName, err := createCodeAndInputFilesHost(code, d.logger)
	if err != nil {
		return "", fmt.Errorf("failed to create code and input files: %w", err)
	}
	d.logger.Info("created code and input files",
		zap.String("code file name", codeFileName),
		zap.String("input file name", inputFileName),
	)

	res, err := d.client.ContainerCreate(ctx, &container.Config{
		Cmd:   getContainerCommand(code, codeFileName, inputFileName),
		Image: code.Image,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: config.GetHostLanguageCodePath(code.Language),
				Target: TargetMountPath,
			},
		},
		// don't let the containers use any network
		NetworkMode: "none",
		RestartPolicy: container.RestartPolicy{
			Name: "no",
		},
		// We are reading the container logs to get the output. So it's better to disable and have a separate thread to delete stale containers
		AutoRemove: false,
		// Drop all the capabilities
		CapDrop:    []string{"ALL"},
		Privileged: false,
		// Set the memory limit to 1GB
		Resources: container.Resources{
			// set 500 MB as the memory limit in bytes
			Memory: 500 * 1024 * 1024,
			Ulimits: []*units.Ulimit{
				{
					Name: "nproc",
					Soft: 64,
					Hard: 128,
				},
				{
					Name: "nofile",
					Soft: 64,
					Hard: 128,
				},
				{
					Name: "core",
					Soft: 0,
					Hard: 0,
				},
				{
					// Maximum file size that can be created by the process (output file in our case)
					Name: "fsize",
					Soft: 20 * 1024 * 1024,
					Hard: 20 * 1024 * 1024,
				},
			},
		},
	}, nil, nil, getContainerName())

	if err != nil {
		return "", fmt.Errorf("failed to create a container: %w", err)
	}

	d.logger.Info("created the container, waiting for the container to start")
	if err = d.client.ContainerStart(ctx, res.ID, container.StartOptions{}); err != nil {
		return res.ID, fmt.Errorf("failed to start the container after creating: %w", err)
	}

	d.logger.Info("container started, waiting for the container to exit")
	statusCh, errCh := d.client.ContainerWait(ctx, res.ID, container.WaitConditionNotRunning)
	ticker := time.NewTicker(MAX_EXECUTION_TIME)

	select {
	case <-ticker.C:
		d.logger.Info("container has been running for more than 60 seconds, killing the container",
			zap.String("container ID", res.ID),
		)
		if err := d.client.ContainerKill(ctx, res.ID, "KILL"); err != nil {
			return "", fmt.Errorf("failed to kill the container: %w", err)
		}
		d.logger.Info("killed the container")
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

	output := stdoutBuf.String()
	if stderrBuf.Len() > 0 {
		output += "\n" + stderrBuf.String()
	}

	return output, nil
}

func deleteStaleFiles(dir string, logger *zap.Logger) error {
	now := time.Now()
	threshold := now.Add(-5 * time.Minute)

	files, err := os.ReadDir(dir)
	if err != nil {
		logger.Error("failed to read the directory", zap.String("directory name", dir))
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			logger.Debug("failed to get file info of the file", zap.String("file name", file.Name()))
			continue
		}
		modTime := fileInfo.ModTime()

		if modTime.Before(threshold) {
			filePath := filepath.Join(dir, file.Name())
			err := os.Remove(filePath)
			if err != nil {
				logger.Error("error deleting file", zap.String("file name", filePath))
			} else {
				logger.Debug("deleted the file", zap.String("file name", filePath))
			}
		}
	}
	return nil
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

		deleteStaleFiles(config.GetHostLanguageCodePath(config.Cpp), d.logger)
		deleteStaleFiles(config.GetHostLanguageCodePath(config.Golang), d.logger)

		time.Sleep(GarbageCollectionTimeWindow)
	}
}

func getContainerName() string {
	return fmt.Sprintf("code-execution-%s", uuid.New().String())
}

func getContainerCommand(code *Code, codeFileName, inputFileName string) []string {
	codeFilePath := getFilePathContainer(TargetMountPath, codeFileName)
	inputFilePath := getFilePathContainer(TargetMountPath, inputFileName)

	command := code.Command
	command = strings.Replace(command, "{{LANGUAGE}}", string(code.Language), -1)
	command = strings.Replace(command, "{{FILE}}", codeFilePath, -1)
	command = strings.Replace(command, "{{INPUT}}", inputFilePath, -1)

	fmt.Println("Command to be executed:", command)

	return []string{
		"sh", "-c",
		command,
	}
}
