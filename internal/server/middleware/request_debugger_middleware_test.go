package middleware_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nix-united/golang-gin-boilerplate/internal/server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestRequestDebuggerMiddleware(t *testing.T) {
	jsonRequestBody := map[string]any{"data": map[string]any{"message": "request-body"}}
	jsonResponseBody := map[string]any{"data": map[string]any{"message": "response-body"}}

	rawJSONRequestBody, err := json.Marshal(jsonRequestBody)
	require.NoError(t, err)

	rawJSONResponseBody, err := json.Marshal(jsonResponseBody)
	require.NoError(t, err)

	jsonHTTPRequest := httptest.NewRequest(http.MethodPost, "/some-endpoint", bytes.NewReader(rawJSONRequestBody))
	jsonHTTPRequest.Header.Set("Content-Type", "application/json")

	textRequestBody := "text-request-body"
	textResponseBody := "text-response-body"

	textHTTPRequest := httptest.NewRequest(http.MethodPost, "/some-endpoint", strings.NewReader(textRequestBody))
	textHTTPRequest.Header.Set("Content-Type", "text/plain")

	emptyHTTPRequest := httptest.NewRequest(http.MethodPost, "/some-endpoint", http.NoBody)

	emptyJSONHttpRequest := httptest.NewRequest(http.MethodPost, "/some-endpoint", http.NoBody)
	emptyJSONHttpRequest.Header.Set("Content-Type", "application/json")

	testCases := map[string]struct {
		inputRequest        *http.Request
		responseContentType string
		responseBody        []byte
		assert              func(t *testing.T, gotLogMessage map[string]any)
	}{
		"It should log request & response JSON body as JSON": {
			inputRequest:        jsonHTTPRequest,
			responseContentType: "application/json",
			responseBody:        rawJSONResponseBody,
			assert: func(t *testing.T, gotLogMessage map[string]any) {
				t.Helper()

				assert.Equal(t, "DEBUG", gotLogMessage["level"])
				assert.Equal(t, "Request/response data", gotLogMessage["msg"])

				assert.Equal(t, jsonRequestBody, gotLogMessage["request_body"])
				assert.Equal(t, jsonResponseBody, gotLogMessage["response_body"])
			},
		},
		"It should log request & response body of not-JSON format as text": {
			inputRequest:        textHTTPRequest,
			responseContentType: "text/plain",
			responseBody:        []byte(textResponseBody),
			assert: func(t *testing.T, gotLogMessage map[string]any) {
				t.Helper()

				assert.Equal(t, textRequestBody, gotLogMessage["request_body"])
				assert.Equal(t, textResponseBody, gotLogMessage["response_body"])
			},
		},
		"It should skip request & response body logging if they're missing": {
			inputRequest: emptyHTTPRequest,
			assert: func(t *testing.T, gotLogMessage map[string]any) {
				t.Helper()

				assert.Equal(t, "Request/response without any data", gotLogMessage["msg"])

				assert.Empty(t, gotLogMessage["request_body"])
				assert.Empty(t, gotLogMessage["response_body"])
			},
		},
		"It should skip request & response body logging if they're missing even if application/json Content-Type specified": {
			inputRequest: emptyJSONHttpRequest,
			assert: func(t *testing.T, gotLogMessage map[string]any) {
				t.Helper()

				assert.Equal(t, "Request/response without any data", gotLogMessage["msg"])

				assert.Empty(t, gotLogMessage["request_body"])
				assert.Empty(t, gotLogMessage["response_body"])
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			writer := NewMockWriter(ctrl)
			requestLoggerMiddleware := middleware.NewRequestDebuggerMiddleware()

			mockedLogger := slog.New(slog.NewJSONHandler(
				writer,
				&slog.HandlerOptions{Level: slog.LevelDebug}),
			)
			defaultLogger := slog.Default()

			slog.SetDefault(mockedLogger)
			t.Cleanup(func() {
				slog.SetDefault(defaultLogger)
			})

			engine := gin.New()
			engine.POST("/some-endpoint", requestLoggerMiddleware.Handle, func(c *gin.Context) {
				_, err := io.ReadAll(c.Request.Body)
				require.NoError(t, err)

				if testCase.responseBody == nil {
					c.Status(http.StatusOK)

					return
				}

				c.Data(http.StatusOK, testCase.responseContentType, testCase.responseBody)
			})

			recorder := httptest.NewRecorder()

			writer.
				EXPECT().
				Write(gomock.Any()).
				DoAndReturn(func(gotRawLogMessage []byte) (int, error) {
					var gotLogMessage map[string]any
					err := json.Unmarshal(gotRawLogMessage, &gotLogMessage)
					require.NoError(t, err)

					testCase.assert(t, gotLogMessage)

					return len(gotRawLogMessage), nil
				})

			engine.ServeHTTP(recorder, testCase.inputRequest)
		})
	}
}
