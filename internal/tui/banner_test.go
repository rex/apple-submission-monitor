package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/require"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

func TestRenderBannerIsFiveRowsAndFits(t *testing.T) {
	t.Parallel()

	lines := renderBanner("Synthetic Alpha", 32, healthPalette(domain.HealthGreen))
	require.Len(t, lines, 5)
	for _, line := range lines {
		require.LessOrEqual(t, lipgloss.Width(line), 32)
	}
	require.NotEmpty(t, strings.Join(lines, ""))
}

func TestInterpolateHex(t *testing.T) {
	t.Parallel()

	require.Equal(t, "#808080", interpolateHex("#000000", "#FFFFFF", 0.5))
	require.Equal(t, "invalid", interpolateHex("invalid", "#FFFFFF", 0.5))
}
