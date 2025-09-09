package common_interceptor

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	TraceIDKey = "trace_id"
	LoggerKey  = "logger"
)

func TraceUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	traceID := uuid.NewString()
	reqLogger := zap.L().With(zap.String("trace_id", traceID))
	ctx = context.WithValue(ctx, TraceIDKey, traceID)
	ctx = context.WithValue(ctx, LoggerKey, reqLogger)
	resp, err = handler(ctx, req)
	return resp, err
}
