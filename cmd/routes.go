package main

import (
	"net/http"
	"remote-code-engine/pkg/config"
	codecontainer "remote-code-engine/pkg/container"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterRoutes(r *gin.Engine, client codecontainer.ContainerClient, config *config.ImageConfig) {
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
			EncodedCode:    req.EncodedCode,
			EncodedInput:   req.EncodedInput,
			Language:       req.Language,
			LanguageConfig: config.GetLanguageConfig(req.Language),
		}

		logger.Info("created a code execution request",
			zap.Any("Language", code.Language),
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
