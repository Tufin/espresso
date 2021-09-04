package env_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/espresso/env"
)

func TestGetEnvSensitiveOrExit(t *testing.T) {

	const key, secret = "pass", "my-pass"
	require.NoError(t, os.Setenv(key, secret))
	require.Equal(t, secret, env.GetEnvSensitiveOrExit(key))
}

func TestGetEnvWithDefault(t *testing.T) {

	const value = "test-me"
	require.Equal(t, value, env.GetEnvWithDefault("STAM_KEY", value))
}
