package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tee-nullpointer/go-common-kit/pkg/logger"
	"go.uber.org/zap"
)

type GinServer struct {
	mode   string
	router *gin.Engine
	server *http.Server
}

func NewGinServer(mode string) *GinServer {
	gin.SetMode(mode)
	return &GinServer{
		mode:   mode,
		router: gin.New(),
	}
}

func (s *GinServer) SetupRouter(setup func(r *gin.Engine), middlewares ...gin.HandlerFunc) {
	s.router.Use(gin.Recovery())
	s.router.Use(middlewares...)
	setup(s.router)
}

func (s *GinServer) Start(host string, port string) {
	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: s.router,
	}
	logger.Info("Gin Server starting", zap.String("address", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

func (s *GinServer) Shutdown() {
	if err := s.server.Close(); err != nil {
		logger.Error("Failed to stop Gin Server", zap.Error(err))
	}
	logger.Info("Gin Server stopped")
}

func (s *GinServer) GracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	logger.Info("Initiating graceful shutdown...")
	if err := s.server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
		return
	}
	logger.Info("Server gracefully stopped")
}
