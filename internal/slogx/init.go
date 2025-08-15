package slogx

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/nix-united/golang-gin-boilerplate/internal/config"
)

const permissionsForLoggsFile = 0o644

func Init(config config.LoggerConfig) (err error) {
	writer := io.Writer(os.Stdout)
	if config.File != "" {

		writer, err = os.OpenFile(config.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, permissionsForLoggsFile)
		if err != nil {
			return fmt.Errorf("open file %s: %w", config.File, err)
		}
	}

	level := slog.LevelDebug
	if config.Level != "" {
		if err = level.UnmarshalText([]byte(config.Level)); err != nil {
			return fmt.Errorf("parse log level %s: %w", config.Level, err)
		}
	}

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("get hostname: %w", err)
	}

	jsonHandler := slog.NewJSONHandler(writer, &slog.HandlerOptions{AddSource: config.AddSource, Level: level})

	traceHandler := newTraceHandler(jsonHandler)

	logger := slog.New(traceHandler).
		With("application", config.Application).
		With("hostname", hostname)

	slog.SetDefault(logger)

	return nil
}
