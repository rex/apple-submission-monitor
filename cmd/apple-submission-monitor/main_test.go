package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunVersion(t *testing.T) {
	previous := version
	version = "test-version"
	t.Cleanup(func() { version = previous })
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run(context.Background(), []string{"--version"}, &stdout, &stderr)
	require.Equal(t, 0, code)
	require.Equal(t, "test-version\n", stdout.String())
	require.Empty(t, stderr.String())
}

func TestRunRejectsUnknownFlag(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := run(context.Background(), []string{"--unknown"}, &stdout, &stderr)

	require.Equal(t, 2, code)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), "flag provided but not defined")
}

func TestRunSanitizesConfigurationFailure(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "invalid.yaml")
	require.NoError(t, os.WriteFile(path, []byte("poll_interval: invalid"), 0o600))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run(
		context.Background(),
		[]string{"--config", path},
		&stdout,
		&stderr,
	)
	require.Equal(t, 1, code)
	require.Equal(t, "configuration could not be loaded\n", stderr.String())
	require.NotContains(t, stderr.String(), path)
}
