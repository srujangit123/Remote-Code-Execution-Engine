package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	dockerC "remote-code-engine/pkg"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

func SetUpDocker() {

}

func StartServer(c *client.Client) error {
	r := gin.Default()

	server := &http.Server{
		Addr:         ":9000",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	RegisterRoutes(r, c)
	return server.ListenAndServe()
}

type UserCode struct {
	Code     string `json: "code"`
	Language string `json: "language"`
}

type Response struct {
	Output string `json: "output"`
}

func RegisterRoutes(r *gin.Engine, client *client.Client) {
	r.POST("/api/v1/submit", func(ctx *gin.Context) {
		var code UserCode
		if err := ctx.BindJSON(&code); err != nil {
			return
		}
		fmt.Printf("%+v\n", code)

		f, err := os.Create("/home/srujan/Documents/code/cpp/uploaded.cpp")
		if err != nil {
			panic(err)
		}

		data, err := base64.StdEncoding.DecodeString(code.Code)
		if err != nil {
			log.Fatal("error:", err)
		}

		fmt.Printf("%q\n", data)

		_, err = f.Write([]byte(data))
		if err != nil {
			panic(err)
		}

		fmt.Println("successfully wrote the content to the file")
		id, err := dockerC.CreateContainer(context.Background(), client, &dockerC.Code{
			Language: "cpp",
			FileName: "uploaded.cpp",
		})
		if err != nil {
			panic(err)
		}

		fmt.Println(id)

		// containers, err := dockerC.GetContainers(context.Background(), client, &container.ListOptions{})
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Printf("Running containers: %+v\n", containers)

		output, err := dockerC.GetCodeOutput(context.Background(), client, &dockerC.Code{
			Language: code.Language,
			FileName: "uploaded.cpp",
		})
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, Response{
			Output: output,
		})
	})
}

func main() {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	StartServer(client)

	///////
	containers, err := dockerC.GetContainers(context.Background(), client, &container.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Running containers: %+v\n", containers)

	id, err := dockerC.CreateContainer(context.Background(), client, &dockerC.Code{
		Language: "cpp",
		FileName: "sample.txt",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(id)

	containers, err = dockerC.GetContainers(context.Background(), client, &container.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Running containers: %+v\n", containers)

}

// We will convert the whole code into base64 string along with the language and pass these two things to the server.
