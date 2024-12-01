package main

import (
	"net/http"
	codecontainer "remote-code-engine/pkg/container"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RegisterRoutes(r *gin.Engine, client codecontainer.ContainerClient) {
	r.POST("/api/v1/submit", func(ctx *gin.Context) {
		var req Request
		if err := ctx.BindJSON(&req); err != nil {
			return
		}

		logger.Info("received a request at",
			zap.String("route", "/api/v1/submit"),
			zap.Any("request params", req),
		)

		code := &codecontainer.Code{
			EncodedCode: req.EncodedCode,
			Language:    req.Language,
			FileName:    uuid.New().String(), // without any extension at the end.
		}

		logger.Info("created a container create request",
			zap.String("Language", code.Language),
			zap.String("FileName", code.FileName),
		)

		output, err := client.ExecuteCode(ctx, code)
		if err != nil {
			logger.Error("Error executing code", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute code"})
		}

		logger.Info("request completed")
		ctx.JSON(http.StatusOK, Response{
			Output: output,
		})
	})
}
