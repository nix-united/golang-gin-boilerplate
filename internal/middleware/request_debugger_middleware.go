package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
)

// requestDebuggerMiddleware is a logging middleware that logs request and response bodies with DEBUG level logs.
// Warning: Do not use this middleware with endpoints containing sensitive information.
//
// Based on the Content-Type, it determines how the body will be formatted.
// If the content type is application/json, the body will be logged as JSON; otherwise, it will be logged as a string.
type requestDebuggerMiddleware struct{}

func NewRequestDebuggerMiddleware() gin.HandlerFunc {
	return (&requestDebuggerMiddleware{}).handle
}

func (m *requestDebuggerMiddleware) handle(c *gin.Context) {
	ctx := c.Request.Context()

	if !slog.Default().Enabled(ctx, slog.LevelDebug) {
		return
	}

	requestBody, err := m.getRequestBody(c)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read request body for request debugging", "err", err.Error())

		return
	}

	responseBodyGetter := m.getResponseBodyGetter(c)

	c.Next()

	var attrs []any
	if requestBody != nil {
		attrs = append(attrs, "request_body", requestBody)
	}

	responseBody := responseBodyGetter(c)
	if responseBody != nil {
		attrs = append(attrs, "response_body", responseBody)
	}

	message := "Request/response data"
	if len(attrs) == 0 {
		message = "Request/response without any data"
	}

	slog.DebugContext(c.Request.Context(), message, attrs...)

	return
}

func (m *requestDebuggerMiddleware) getRequestBody(c *gin.Context) (any, error) {
	if c.Request.Body == nil {
		return nil, nil
	}

	rawRequestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, fmt.Errorf("read request body: %w", err)
	}

	c.Request.Body = io.NopCloser(bytes.NewReader(rawRequestBody))

	if strings.HasPrefix(c.Request.Header.Get("Content-Type"), "application/json") {
		return json.RawMessage(rawRequestBody), nil
	}

	return string(rawRequestBody), nil
}

func (m *requestDebuggerMiddleware) getResponseBodyGetter(c *gin.Context) func(c *gin.Context) any {
	storer := newCapturingResponseWriter(c.Writer)
	c.Writer = storer

	return func(c *gin.Context) any {
		if storer.response == nil {
			return nil
		}

		if strings.HasPrefix(c.Request.Response.Header.Get("Content-Type"), "application/json") {
			return json.RawMessage(storer.response)
		}

		return string(storer.response)
	}
}

// capturingResponseWriter stores the written response by the handler into its field.
// This is used to automate response logging.
type capturingResponseWriter struct {
	gin.ResponseWriter

	response []byte
}

func newCapturingResponseWriter(writer gin.ResponseWriter) *capturingResponseWriter {
	return &capturingResponseWriter{ResponseWriter: writer}
}

func (w *capturingResponseWriter) Write(response []byte) (int, error) {
	w.response = response

	n, err := w.ResponseWriter.Write(response)
	if err != nil {
		return n, fmt.Errorf("write response after capture: %w", err)
	}

	return n, nil
}
