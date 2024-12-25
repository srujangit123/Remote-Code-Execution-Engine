package main

import (
	"context"
	"net/http"
	"remote-code-engine/pkg/config"
	codecontainer "remote-code-engine/pkg/container"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger *zap.Logger

const (
	ADDR = ":9000"
)

func init() {
	logger, _ = zap.NewProduction()
}

func StartServer(cli codecontainer.ContainerClient, config *config.ImageConfig) error {
	r := gin.Default()
	logger.Info("starting the server",
		zap.String("Address", ADDR),
	)

	server := &http.Server{
		Addr:         ":9000",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	RegisterRoutes(r, cli, config)
	return server.ListenAndServe()
}

func main() {
	defer logger.Sync()
	// This should be given as a command line argument.
	imageConfig, err := config.LoadConfig("../config.yml")
	if err != nil {
		logger.Error("failed to load the config file",
			zap.Error(err),
		)
		panic(err)
	}

	logger.Debug("loaded the config file",
		zap.Any("config", imageConfig),
	)

	cli, err := codecontainer.NewDockerClient(nil, logger)
	if err != nil {
		logger.Error("failed to create a docker client",
			zap.Error(err),
		)
		panic(err)
	}

	go cli.FreeUpZombieContainers(context.Background())

	StartServer(cli, imageConfig)
}
