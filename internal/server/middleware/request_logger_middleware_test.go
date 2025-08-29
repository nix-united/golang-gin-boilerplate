package middleware_test

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nix-united/golang-gin-boilerplate/internal/server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRequestLoggerMiddleware(t *testing.T) {
	testCases := map[string]struct {
		responseStatusCode  int
		errDuringProcessing error
		assertLogMessage    func(t *testing.T, gotLogMessage map[string]any)
	}{
		"It should log default message": {
			responseStatusCode: http.StatusCreated,
			assertLogMessage: func(t *testing.T, gotLogMessage map[string]any) {
				t.Helper()

				assert.Equal(t, slog.LevelInfo.String(), gotLogMessage["level"])
				assert.Equal(t, "Processed request", gotLogMessage["msg"])

				httpAttributes, ok := gotLogMessage["http"].(map[string]any)
				require.True(t, ok, "missing http attribute")

				assert.Equal(t, http.MethodPost, httpAttributes["method"])
				assert.Equal(t, float64(http.StatusCreated), httpAttributes["status"])
				assert.Equal(t, "/some-endpoint", httpAttributes["path"])
			},
		},
		"It should use ERROR level when status code >= 500": {
			responseStatusCode: http.StatusInternalServerError,
			assertLogMessage: func(t *testing.T, gotLogMessage map[string]any) {
				t.Helper()

				assert.Equal(t, slog.LevelError.String(), gotLogMessage["level"])
				assert.Equal(t, "Processed request", gotLogMessage["msg"])

				httpAttributes, ok := gotLogMessage["http"].(map[string]any)
				require.True(t, ok, "missing http attribute")

				assert.Equal(t, http.MethodPost, httpAttributes["method"])
				assert.Equal(t, float64(http.StatusInternalServerError), httpAttributes["status"])
				assert.Equal(t, "/some-endpoint", httpAttributes["path"])
			},
		},
		"It should log error message when it occurred during processing": {
			responseStatusCode:  http.StatusConflict,
			errDuringProcessing: errors.New("some error during processing"),
			assertLogMessage: func(t *testing.T, gotLogMessage map[string]any) {
				t.Helper()

				assert.Equal(t, slog.LevelWarn.String(), gotLogMessage["level"])
				assert.Equal(t, "Processed request", gotLogMessage["msg"])
				assert.Equal(t, "some error during processing", gotLogMessage["error"])

				httpAttributes, ok := gotLogMessage["http"].(map[string]any)
				require.True(t, ok, "missing http attribute")

				assert.Equal(t, http.MethodPost, httpAttributes["method"])
				assert.Equal(t, float64(http.StatusConflict), httpAttributes["status"])
				assert.Equal(t, "/some-endpoint", httpAttributes["path"])
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			writer := NewMockWriter(ctrl)
			tracer := NewMocktracer(ctrl)

			requestLoggerMiddleware := middleware.NewRequestLoggerMiddleware(tracer)

			mockedLogger := slog.New(slog.NewJSONHandler(writer, nil))
			defaultLogger := slog.Default()
			slog.SetDefault(mockedLogger)
			t.Cleanup(func() {
				slog.SetDefault(defaultLogger)
			})

			engine := gin.New()
			engine.POST("/some-endpoint", requestLoggerMiddleware.Handle, func(c *gin.Context) {
				if testCase.errDuringProcessing != nil {
					c.Error(testCase.errDuringProcessing)
				}

				c.Status(testCase.responseStatusCode)
			})

			recorder := httptest.NewRecorder()

			httpRequest := httptest.NewRequest(http.MethodPost, "/some-endpoint", http.NoBody)

			tracer.
				EXPECT().
				Start(gomock.Any()).
				DoAndReturn(func(ctx context.Context) (context.Context, error) {
					return ctx, nil
				})

			writer.
				EXPECT().
				Write(gomock.Any()).
				DoAndReturn(func(gotRawLogMessage []byte) (int, error) {
					var gotLogMessage map[string]any
					err := json.Unmarshal(gotRawLogMessage, &gotLogMessage)
					require.NoError(t, err)

					testCase.assertLogMessage(t, gotLogMessage)

					return len(gotRawLogMessage), nil
				})

			engine.ServeHTTP(recorder, httpRequest)

			assert.Equal(t, testCase.responseStatusCode, recorder.Result().StatusCode)
		})
	}
}
