package common_interceptor

import (
	"context"

	"github.com/tee-nullpointer/go-common-kit/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoveryUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	log := logger.GetLogger(ctx)
	defer func() {
		if r := recover(); r != nil {
			log.Error("Panic", zap.String("method", info.FullMethod), zap.Any("req", req), zap.Any("panic", r))
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()
	return handler(ctx, req)
}
