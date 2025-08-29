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

// RequestDebuggerMiddleware is a logging middleware that logs request and
// response bodies with DEBUG level logs.
//
// Warning: Do not use this middleware with endpoints containing sensitive
// information.
//
// Based on the Content-Type, it determines how the body will be formatted.
// If the content type is application/json, the body will be logged as JSON;
// otherwise, it will be logged as a string.
type RequestDebuggerMiddleware struct{}

func NewRequestDebuggerMiddleware() *RequestDebuggerMiddleware {
	return &RequestDebuggerMiddleware{}
}

// Handle captures, when DEBUG logging is enabled, the incoming request body
// and the outgoing response body and emits them in a structured debug log
// after the handler chain finishes.
func (m *RequestDebuggerMiddleware) Handle(c *gin.Context) {
	ctx := c.Request.Context()

	if !slog.Default().Enabled(ctx, slog.LevelDebug) {
		return
	}

	requestBody, err := m.getRequestBody(c)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read request body for request debugging", "err", err)

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
}

// getRequestBody reads the request body and returns a value for logging.
// It restores c.Request.Body so downstream handlers can read it again.
func (m *RequestDebuggerMiddleware) getRequestBody(c *gin.Context) (any, error) {
	if c.Request.Body == nil {
		return nil, nil
	}

	rawRequestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, fmt.Errorf("read request body: %w", err)
	}

	c.Request.Body = io.NopCloser(bytes.NewReader(rawRequestBody))

	if len(rawRequestBody) == 0 {
		return nil, nil
	}

	if strings.HasPrefix(c.Request.Header.Get("Content-Type"), "application/json") {
		return json.RawMessage(rawRequestBody), nil
	}

	return string(rawRequestBody), nil
}

// getResponseBodyGetter wraps the current gin.ResponseWriter with a capturing
// writer that buffers the response payload as it is written.
//
// It returns a closure that, when called after handler execution,
// yields a value suitable for logging.
func (m *RequestDebuggerMiddleware) getResponseBodyGetter(c *gin.Context) func(c *gin.Context) any {
	storer := newCapturingResponseWriter(c.Writer)
	c.Writer = storer

	return func(c *gin.Context) any {
		if storer.response == nil {
			return nil
		}

		if strings.HasPrefix(c.Writer.Header().Get("Content-Type"), "application/json") {
			return json.RawMessage(storer.response)
		}

		return string(storer.response)
	}
}

// capturingResponseWriter decorates a gin.ResponseWriter to capture the
// response body bytes written by handlers. It forwards all writes to the
// underlying writer while storing a copy in memory for later logging.
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
