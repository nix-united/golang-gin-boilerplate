package setup

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

const (
	appReadinessTimeout   = 5 * time.Minute
	mysqlReadinessTimeout = time.Minute
	cleanupTimeout        = 2 * time.Minute
)

func CleanupContext(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.WithoutCancel(parent), cleanupTimeout)
}

func terminateFailedContainer(ctx context.Context, container testcontainers.Container) error {
	cleanupCtx, cancel := CleanupContext(ctx)
	defer cancel()

	if err := testcontainers.TerminateContainer(container, testcontainers.StopContext(cleanupCtx)); err != nil {
		return fmt.Errorf("terminate container that failed to start: %w", err)
	}

	return nil
}
