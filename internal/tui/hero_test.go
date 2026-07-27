package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/require"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

func TestHeroUsesSharedLargeFontAndRightAlignedAge(t *testing.T) {
	t.Parallel()

	hero := prepareHero("Synthetic Alpha", "8D 17H", 156, 18, false, ageBandRed)
	require.NotEmpty(t, hero.font)
	require.NotEmpty(t, hero.name.figure)
	require.NotEmpty(t, hero.age.figure)
	require.True(t, hero.age.bloody)
	require.Contains(t, strings.Join(hero.age.figure, ""), "●")

	lines := hero.render(
		156,
		healthPalette(domain.HealthYellow),
		submissionAgePalette(ageBandRed),
		12,
	)
	for _, line := range lines {
		require.Equal(t, 156, lipgloss.Width(line))
	}
}
