package setup

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go/network"
)

func SetupNetwork(ctx context.Context) (string, func(ctx context.Context) error, error) {
	dockerNetwork, err := network.New(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("new network: %w", err)
	}

	shutdown := func(ctx context.Context) error {
		if err := dockerNetwork.Remove(ctx); err != nil {
			return fmt.Errorf("remove network: %w", err)
		}

		return nil
	}

	return dockerNetwork.Name, shutdown, nil
}
