package main

import (
	"fmt"
	common_handler "go-common-kit/handler"
	common_middleware "go-common-kit/middleware"
	"go-common-kit/pkg/env"
	"go-common-kit/pkg/logger"
	"go-common-kit/server"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	env.LoadEnvFile()
	logLevel := env.GetEnv("LOG_LEVEL", "info")
	logFormat := env.GetEnv("LOG_FORMAT", "json")
	if err := logger.InitLogger(logLevel, logFormat); err != nil {
		panic(fmt.Sprintf("Fail to initialize logger: %v", err))
	}
	defer logger.Sync()

	ginMode := env.GetEnv("GIN_MODE", "release")
	ginPort := env.GetEnv("GIN_PORT", "8080")

	logger.Info("Starting Gin Http Server",
		zap.String("port", ginPort),
		zap.String("mode", ginMode),
	)

	ginServer := server.NewGinServer(ginMode)
	ginServer.SetupRouter(setupRouter,
		common_middleware.LoggingMiddleware())
	go ginServer.Start("localhost", ginPort)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
	ginServer.GracefulShutdown()
}

func setupRouter(router *gin.Engine) {
	monitor := router.Group("/")
	{
		monitor.GET("/health", common_handler.HealthCheck)
	}
	api := router.Group("/api/v1")
	{
		api.GET("/example", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "example endpoint"})
		})
	}
}
