package config_test

import (
	"testing"

	"github.com/nix-united/golang-gin-boilerplate/internal/config"

	"github.com/caarlos0/env/v11"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplicationConfig_AuthJWTSecret(t *testing.T) {
	parse := func(environment map[string]string) (config.ApplicationConfig, error) {
		var cfg config.ApplicationConfig
		err := env.ParseWithOptions(&cfg, env.Options{Environment: environment})
		return cfg, err
	}

	t.Run("It should fail when AUTH_JWT_SECRET is missing", func(t *testing.T) {
		_, err := parse(map[string]string{})
		require.Error(t, err)
		assert.ErrorContains(t, err, "AUTH_JWT_SECRET")
	})

	t.Run("It should fail when AUTH_JWT_SECRET is empty", func(t *testing.T) {
		_, err := parse(map[string]string{"AUTH_JWT_SECRET": ""})
		require.Error(t, err)
		assert.ErrorContains(t, err, "AUTH_JWT_SECRET")
	})

	t.Run("It should parse a non-empty AUTH_JWT_SECRET", func(t *testing.T) {
		cfg, err := parse(map[string]string{"AUTH_JWT_SECRET": "secret"})
		require.NoError(t, err)
		assert.Equal(t, "secret", cfg.AuthConfig.Secret)
	})
}
