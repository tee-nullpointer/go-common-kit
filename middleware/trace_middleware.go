package common_middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	TraceIDKey = "trace_id"
	LoggerKey  = "logger"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		reqLogger := zap.L().With(zap.String("trace_id", traceID))
		c.Set(TraceIDKey, traceID)
		c.Set(LoggerKey, LoggerKey)
		ctx := context.WithValue(c.Request.Context(), TraceIDKey, traceID)
		ctx = context.WithValue(ctx, LoggerKey, reqLogger)
		c.Request = c.Request.WithContext(ctx)

		c.Header("X-Trace-ID", traceID)
		c.Next()
	}
}
