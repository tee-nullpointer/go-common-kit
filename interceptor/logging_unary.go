package common_interceptor

import (
	"context"
	"time"

	"github.com/tee-nullpointer/go-common-kit/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log := logger.GetLogger(ctx)
	start := time.Now()
	log.Info("GRPC Request", zap.String("method", info.FullMethod), zap.Any("req", req))
	resp, err := handler(ctx, req)
	st, _ := status.FromError(err)
	log.Info("GRPC Response", zap.String("method", info.FullMethod), zap.Duration("duration", time.Since(start)), zap.Any("response", resp), zap.String("error", st.Message()))
	return resp, err
}
