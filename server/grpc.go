package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
	logger   *zap.Logger
	opts     []grpc.ServerOption
}

func NewGRPCServer(opts ...grpc.ServerOption) *GRPCServer {
	return &GRPCServer{
		server: grpc.NewServer(opts...),
		logger: zap.L(),
		opts:   opts,
	}
}

func (s *GRPCServer) Start(host string, port string) {
	addr := fmt.Sprintf("%s:%s", host, port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Fatal("Failed to listen", zap.Error(err))
	}
	s.listener = ln
	s.logger.Info("gRPC Server starting", zap.String("address", addr))
	if err := s.server.Serve(ln); err != nil {
		s.logger.Fatal("Failed to start gRPC server", zap.Error(err))
	}
}

func (s *GRPCServer) Shutdown() {
	s.server.Stop()
	s.logger.Info("gRPC Server stopped")
}

func (s *GRPCServer) GracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ch := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(ch)
	}()
	select {
	case <-ch:
		s.logger.Info("gRPC Server gracefully stopped")
	case <-ctx.Done():
		s.logger.Warn("gRPC Server graceful shutdown timed out")
	}
}

func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}
