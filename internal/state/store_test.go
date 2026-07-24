package state_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rex/apple-submission-monitor/internal/domain"
	"github.com/rex/apple-submission-monitor/internal/state"
)

func TestStoreRoundTripUsesPrivatePermissions(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "state", "state.json")
	store := state.NewStore(path)
	cards := []domain.Submission{{
		ID:           "review-alpha",
		AppID:        "app-alpha",
		AppName:      "Synthetic Alpha",
		Health:       domain.HealthYellow,
		Acknowledged: true,
	}}

	require.NoError(t, store.Save(cards))
	loaded, err := store.Load()
	require.NoError(t, err)
	require.Equal(t, cards, loaded)

	info, err := os.Stat(path)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestStoreRejectsUnknownVersion(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "state.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"version":99,"cards":[]}`), 0o600))

	_, err := state.NewStore(path).Load()
	require.ErrorContains(t, err, "unsupported monitor state version")
}

func TestStoreMissingFileIsEmpty(t *testing.T) {
	t.Parallel()

	cards, err := state.NewStore(filepath.Join(t.TempDir(), "missing.json")).Load()
	require.NoError(t, err)
	require.Empty(t, cards)
}
