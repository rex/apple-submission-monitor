package asc_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/rex/apple-submission-monitor/internal/asc"
)

func TestCommandRunnerExecutesArgumentVector(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "fake-asc")
	script := "#!/bin/sh\nprintf '%s\\n' \"$*\"\n"
	require.NoError(t, os.WriteFile(path, []byte(script), 0o600))
	// #nosec G302 -- synthetic test fixture must be executable.
	require.NoError(t, os.Chmod(path, 0o700))
	runner := asc.NewRunner(path, 5*time.Second)

	output, err := runner.Run(context.Background(), "status", "--output", "json")
	require.NoError(t, err)
	require.Equal(t, "status --output json\n", string(output))
}

func TestCommandRunnerReturnsSanitizedFailure(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "fake-asc")
	script := "#!/bin/sh\nprintf 'private synthetic output' >&2\nexit 1\n"
	require.NoError(t, os.WriteFile(path, []byte(script), 0o600))
	// #nosec G302 -- synthetic test fixture must be executable.
	require.NoError(t, os.Chmod(path, 0o700))
	runner := asc.NewRunner(path, 5*time.Second)

	_, err := runner.Run(context.Background(), "status", "--app", "app-alpha")
	require.ErrorIs(t, err, asc.ErrCommandFailed)
	require.NotContains(t, err.Error(), "app-alpha")
	require.NotContains(t, err.Error(), "private synthetic output")
}
