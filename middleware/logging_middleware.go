package common_middleware

import (
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tee-nullpointer/go-common-kit/pkg/logger"
	"go.uber.org/zap"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := logger.GetLogger(c.Request.Context())
		start := time.Now()
		l.Info("HTTP Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
		responseBody := &bytes.Buffer{}
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           responseBody,
		}
		c.Writer = writer

		c.Next()

		duration := time.Since(start)

		l.Info("HTTP Response",
			zap.Int("status_code", c.Writer.Status()),
			zap.Duration("response_time", duration),
			zap.Int("response_size", responseBody.Len()),
			zap.String("response_body", responseBody.String()),
		)
	}
}
