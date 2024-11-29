package main

import (
	"context"
	"fmt"
	"net/http"
	codecontainer "remote-code-engine/pkg/container"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func StartServer() error {
	r := gin.Default()

	server := &http.Server{
		Addr:         ":9000",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	cli, err := codecontainer.NewDockerClient(nil)
	if err != nil {
		panic(err)
	}

	RegisterRoutes(r, cli)
	return server.ListenAndServe()
}

type Request struct {
	EncodedCode string `json: "code"`
	Language    string `json: "language"`
}

type Response struct {
	Output string `json: "output"`
}

func RegisterRoutes(r *gin.Engine, client codecontainer.ContainerClient) {
	r.POST("/api/v1/submit", func(ctx *gin.Context) {
		var req Request
		if err := ctx.BindJSON(&req); err != nil {
			return
		}
		fmt.Printf("%+v\n", req)

		code := &codecontainer.Code{
			EncodedCode: req.EncodedCode,
			Language:    req.Language,
			FileName:    uuid.New().String(), // without any extension at the end.
		}

		fmt.Println("successfully wrote the content to the file")
		id, err := client.CreateAndStartContainer(context.Background(), code)
		if err != nil {
			panic(err)
		}

		fmt.Println(id)

		output, err := client.GetContainerOutput(ctx, code)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, Response{
			Output: output,
		})
	})
}

func main() {
	StartServer()
}

// We will convert the whole code into base64 string along with the language and pass these two things to the server.
