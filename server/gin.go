package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GinServer struct {
	mode   string
	router *gin.Engine
	server *http.Server
	logger *zap.Logger
}

func NewGinServer(mode string) *GinServer {
	gin.SetMode(mode)
	return &GinServer{
		mode:   mode,
		router: gin.New(),
		logger: zap.L(),
	}
}

func (s *GinServer) Start(host string, port string) {
	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: s.router,
	}
	s.logger.Info("Gin Server starting", zap.String("address", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Fatal("Failed to start server", zap.Error(err))
	}
}

func (s *GinServer) Shutdown() {
	if err := s.server.Close(); err != nil {
		s.logger.Error("Failed to stop Gin Server", zap.Error(err))
	}
	s.logger.Info("Gin Server stopped")
}

func (s *GinServer) GracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.logger.Info("Initiating graceful shutdown...")
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("Server forced to shutdown", zap.Error(err))
		return
	}
	s.logger.Info("Server gracefully stopped")
}

func (s *GinServer) GetRouter() *gin.Engine {
	return s.router
}
