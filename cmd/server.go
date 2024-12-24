package main

import (
	"context"
	"net/http"
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

func StartServer(cli codecontainer.ContainerClient) error {
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

	RegisterRoutes(r, cli)
	return server.ListenAndServe()
}

type Request struct {
	EncodedCode string `json:"code"`
	Language    string `json:"language"`
}

type Response struct {
	Output string `json:"output"`
}

func main() {
	defer logger.Sync()
	_ = ParseFlags()
	cli, err := codecontainer.NewDockerClient(nil, logger)
	if err != nil {
		logger.Error("failed to create a docker client",
			zap.Error(err),
		)
		panic(err)
	}

	go cli.FreeUpZombieContainers(context.Background())

	StartServer(cli)
}
