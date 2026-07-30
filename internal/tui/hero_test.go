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

	hero := prepareHero("Synthetic Alpha", "8D 17H", 156, 18, false, true)
	require.NotEmpty(t, hero.font)
	require.NotEmpty(t, hero.name.figure)
	require.NotEmpty(t, hero.right.figure)
	require.True(t, hero.right.bloody)
	require.Contains(t, strings.Join(hero.right.figure, ""), "●")

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

func TestApprovedHeroUsesGiantCleanVictoryText(t *testing.T) {
	t.Parallel()

	hero := prepareHero("Synthetic Alpha", "APPROVED", 180, 18, false, false)
	require.NotEmpty(t, hero.font)
	require.NotEmpty(t, hero.right.figure)
	require.False(t, hero.right.bloody)

	lines := hero.render(180, approvalPalette(), approvalPalette(), 24)
	for _, line := range lines {
		require.Equal(t, 180, lipgloss.Width(line))
	}
}
