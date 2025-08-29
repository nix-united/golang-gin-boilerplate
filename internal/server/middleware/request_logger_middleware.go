package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//go:generate go tool mockgen -source=$GOFILE -destination=request_logger_middleware_mock_test.go -package=${GOPACKAGE}_test -typed=true
//go:generate go tool mockgen -destination=io_writer_mock_test.go -package=${GOPACKAGE}_test -typed=true io Writer

type tracer interface {
	Start(ctx context.Context) (context.Context, error)
}

// RequestLoggerMiddleware is a logging middleware that generated trace ID for each request.
type RequestLoggerMiddleware struct {
	tracer tracer
}

func NewRequestLoggerMiddleware(tracer tracer) *RequestLoggerMiddleware {
	return &RequestLoggerMiddleware{tracer: tracer}
}

// handle creates trace and logs request information.
func (m *RequestLoggerMiddleware) Handle(c *gin.Context) {
	start := time.Now()

	ctx, err := m.tracer.Start(c.Request.Context())
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "Failed to start a trace", "err", err.Error())

		return
	}

	c.Request = c.Request.WithContext(ctx)

	c.Next()

	level := slog.LevelInfo
	if c.Writer.Status() >= http.StatusInternalServerError {
		level = slog.LevelError
	}

	attrs := []any{
		slog.Group("http",
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"path", c.FullPath(),
			"duration_ms", time.Since(start).Milliseconds(),
		),
	}

	if len(c.Errors) > 0 {
		attrs = append(attrs, "error", strings.Join(c.Errors.Errors(), "\n"))
	}

	slog.Log(ctx, level, "Processed request", attrs...)
}
