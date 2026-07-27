package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/require"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

func TestRenderBannerUsesLargeFigletArtAndFits(t *testing.T) {
	t.Parallel()

	lines := renderBanner(
		"Synthetic Alpha",
		100,
		10,
		healthPalette(domain.HealthYellow),
		0,
		false,
	)
	require.GreaterOrEqual(t, len(lines), 5)
	for _, line := range lines {
		require.LessOrEqual(t, lipgloss.Width(line), 100)
	}
	rendered := strings.Join(lines, "")
	require.Contains(t, rendered, "█")
	require.Contains(t, rendered, "▓")
}

func TestRejectedBannerIsDistortedAndDripping(t *testing.T) {
	t.Parallel()

	lines := renderBanner(
		"Synthetic Reject",
		140,
		18,
		healthPalette(domain.HealthRed),
		0,
		true,
	)
	require.GreaterOrEqual(t, len(lines), 8)
	rendered := strings.Join(lines, "")
	require.Contains(t, rendered, "█")
	require.Contains(t, rendered, "●")
}

func TestBannerGradientPositionMovesBetweenFrames(t *testing.T) {
	t.Parallel()

	first := gradientPosition(25, 2, 100, 0)
	second := gradientPosition(25, 2, 100, 1)
	require.NotEqual(t, first, second)

	statusFirst := statusAgeGradientPosition(25, 100, 0)
	statusSecond := statusAgeGradientPosition(25, 100, 1)
	require.NotEqual(t, statusFirst, statusSecond)
	require.NotEqual(t, first, statusFirst)
}

func TestGradientTextPreservesContentAndWidth(t *testing.T) {
	t.Parallel()

	const text = "◷ IN STATUS · 2H 5M"
	rendered := gradientText(text, statusAgeColors(), 12)
	require.Equal(t, lipgloss.Width(text), lipgloss.Width(rendered))
	require.Contains(t, rendered, "◷")
	require.Contains(t, rendered, "2")
}

func TestInterpolateHex(t *testing.T) {
	t.Parallel()

	require.Equal(t, "#808080", interpolateHex("#000000", "#FFFFFF", 0.5))
	require.Equal(t, "invalid", interpolateHex("invalid", "#FFFFFF", 0.5))
}

func BenchmarkRenderPreparedBanner(b *testing.B) {
	colors := healthPalette(domain.HealthYellow)
	art := prepareBanner("Synthetic Alpha", 156, 18, false)
	b.ResetTimer()
	for b.Loop() {
		art.render(156, colors, 12)
	}
}
