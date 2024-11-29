package main

import (
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

func StartServer() error {
	r := gin.Default()
	logger.Info("starting the server",
		zap.String("Address", ADDR),
	)

	server := &http.Server{
		Addr:         ":9000",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	cli, err := codecontainer.NewDockerClient(nil, logger)
	if err != nil {
		logger.Error("failed to create a docker client",
			zap.Error(err),
		)
		panic(err)
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
	StartServer()
}

// We will convert the whole code into base64 string along with the language and pass these two things to the server.
