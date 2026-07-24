package tui

import (
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/require"
)

func TestElapsedFormatsUsefulRanges(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 1, 3, 12, 0, 0, 0, time.UTC)
	require.Equal(t, "unknown", elapsed(time.Time{}, now))
	require.Equal(t, "just now", elapsed(now.Add(-30*time.Second), now))
	require.Equal(t, "15m", elapsed(now.Add(-15*time.Minute), now))
	require.Equal(t, "2h 5m", elapsed(now.Add(-2*time.Hour-5*time.Minute), now))
	require.Equal(t, "2d 3h", elapsed(now.Add(-51*time.Hour), now))
	require.Equal(t, "just now", elapsed(now.Add(time.Minute), now))
}

func TestTruncateAndPadBetween(t *testing.T) {
	t.Parallel()

	require.Equal(t, "abc…", truncate("abcdef", 4))
	require.Equal(t, "one line", truncate("one\nline", 20))
	line := padBetween(
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Render("left"),
		"right",
		20,
	)
	require.Equal(t, 20, lipgloss.Width(line))
}
