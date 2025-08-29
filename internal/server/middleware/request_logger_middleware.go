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

// tracer defines the minimal contract for starting a trace for an incoming request.
//
// Tracing enables end-to-end observability: you can correlate logs, measure
// latency, and follow a request across services and goroutines.
type tracer interface {
	Start(ctx context.Context) (context.Context, error)
}

// RequestLoggerMiddleware is a middleware that enriches incoming requests with
// a tracing-aware context and logs message after the request is processed.
type RequestLoggerMiddleware struct {
	tracer tracer
}

func NewRequestLoggerMiddleware(tracer tracer) *RequestLoggerMiddleware {
	return &RequestLoggerMiddleware{tracer: tracer}
}

// Handle starts a trace for the request (via m.tracer), attaches the resulting
// context to the request and then logs a structured summary of
// the request/response.
func (m *RequestLoggerMiddleware) Handle(c *gin.Context) {
	start := time.Now()

	ctx, err := m.tracer.Start(c.Request.Context())
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "Failed to start a trace", "err", err)

		return
	}

	c.Request = c.Request.WithContext(ctx)

	c.Next()

	level := slog.LevelInfo
	switch {
	case c.Writer.Status() >= http.StatusInternalServerError:
		level = slog.LevelError
	case c.Writer.Status() >= http.StatusBadRequest:
		level = slog.LevelWarn
	}

	attrs := []any{
		slog.Group("http",
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"path", c.Request.URL.Path,
			"duration_ms", time.Since(start).Milliseconds(),
		),
	}

	if len(c.Errors) > 0 {
		attrs = append(attrs, "error", strings.Join(c.Errors.Errors(), "\n"))
	}

	slog.Log(ctx, level, "Processed request", attrs...)
}
