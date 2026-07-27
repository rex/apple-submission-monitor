package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/require"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

func TestRenderBannerUsesReadableWordmarkAndFits(t *testing.T) {
	t.Parallel()

	lines := renderBanner("Synthetic Alpha", 32, healthPalette(domain.HealthGreen), false)
	require.Len(t, lines, 3)
	for _, line := range lines {
		require.LessOrEqual(t, lipgloss.Width(line), 32)
	}
	rendered := strings.Join(lines, "")
	require.Contains(t, rendered, "SYNTHETIC ALPHA")
	require.NotContains(t, rendered, "███")
}

func TestInterpolateHex(t *testing.T) {
	t.Parallel()

	require.Equal(t, "#808080", interpolateHex("#000000", "#FFFFFF", 0.5))
	require.Equal(t, "invalid", interpolateHex("invalid", "#FFFFFF", 0.5))
}
