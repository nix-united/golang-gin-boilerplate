package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tracer interface {
	Start(ctx context.Context) (context.Context, error)
}

// requestLoggerMiddleware is a logging middleware that generated trace ID for each request.
type requestLoggerMiddleware struct {
	tracer tracer
}

func NewRequestLoggerMiddleware(tracer tracer) gin.HandlerFunc {
	return (&requestLoggerMiddleware{tracer: tracer}).handle
}

// handle creates trace and logs request information.
func (l *requestLoggerMiddleware) handle(c *gin.Context) {
	ctx, err := l.tracer.Start(c.Request.Context())
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "Failed to start a trace", "err", err.Error())

		return
	}

	c.Request = c.Request.WithContext(ctx)

	c.Next()

	level := slog.LevelInfo
	if c.Request.Response.StatusCode >= http.StatusInternalServerError {
		level = slog.LevelError
	}

	attrs := []any{
		"method", c.Request.Method,
		"status", c.Request.Response.Status,
		"path", c.Request.URL.Path,
	}

	if len(c.Errors) > 0 {
		attrs = append(attrs, "error", c.Errors.JSON())
	}

	slog.Log(c.Request.Context(), level, "Request", slog.Group("http", attrs...))
}
