package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/rex/apple-submission-monitor/internal/config"
)

func TestLoadCreatesPrivateDefault(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "nested", "config.yaml")
	cfg, err := config.Load(path)
	require.NoError(t, err)
	require.Equal(t, 30*time.Second, cfg.PollInterval)
	require.Equal(t, filepath.Join(filepath.Dir(path), "state.json"), cfg.StatePath)

	info, err := os.Stat(path)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestLoadRejectsUnsafePollingRate(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.yaml")
	content := `poll_interval: 1s
discovery_interval: 5m
command_timeout: 25s
animation_interval: 400ms
discovery_workers: 4
asc_path: asc
announcements:
  approved: approved
  rejected: rejected
  status_changed: changed
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	_, err := config.Load(path)
	require.ErrorContains(t, err, "poll_interval must be at least")
}

func TestEnvironmentOverridesPollInterval(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	t.Setenv("ASM_POLL_INTERVAL", "45s")

	cfg, err := config.Load(path)
	require.NoError(t, err)
	require.Equal(t, 45*time.Second, cfg.PollInterval)
	require.NotContains(t, cfg.String(), path)
}
