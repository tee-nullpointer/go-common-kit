package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	commonhandler "github.com/tee-nullpointer/go-common-kit/handler"
	commoninterceptor "github.com/tee-nullpointer/go-common-kit/interceptor"
	commonmiddleware "github.com/tee-nullpointer/go-common-kit/middleware"
	"github.com/tee-nullpointer/go-common-kit/pkg/env"
	"github.com/tee-nullpointer/go-common-kit/pkg/logger"
	"github.com/tee-nullpointer/go-common-kit/proto/pb"
	"github.com/tee-nullpointer/go-common-kit/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type pingServiceServer struct {
	pb.UnimplementedPingServiceServer
}

func (p pingServiceServer) Ping(ctx context.Context, request *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "pong"}, nil
}

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

	zap.L().Info("Starting Gin Http Server",
		zap.String("port", ginPort),
		zap.String("mode", ginMode),
	)

	ginServer := server.NewGinServer(ginMode)
	ginRouter := ginServer.GetRouter()
	ginRouter.Use(gin.Recovery(), commonmiddleware.TraceMiddleware(), commonmiddleware.LoggingMiddleware())
	setupRouter(ginRouter)
	go ginServer.Start("localhost", ginPort)

	grpcPort := env.GetEnv("GRPC_PORT", "9090")
	grpcServer := server.NewGRPCServer(
		grpc.UnaryInterceptor(commoninterceptor.ChainUnaryInterceptors(
			commoninterceptor.RecoveryUnaryInterceptor,
			commoninterceptor.TraceUnaryInterceptor,
			commoninterceptor.LoggingUnaryInterceptor,
		)),
	)
	pb.RegisterPingServiceServer(grpcServer.GetServer(), &pingServiceServer{})
	go grpcServer.Start("localhost", grpcPort)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	zap.L().Info("Received shutdown signal", zap.String("signal", sig.String()))
	ginServer.GracefulShutdown()
	grpcServer.GracefulShutdown()
}

func setupRouter(router *gin.Engine) {
	monitor := router.Group("/")
	{
		monitor.GET("/health", commonhandler.HealthCheck)
	}
	api := router.Group("/api/v1")
	{
		api.GET("/example", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "example endpoint"})
		})
	}
}
