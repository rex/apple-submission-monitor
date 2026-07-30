package tui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/require"
)

func TestMetadataRailsFillWidthAndShowAppChecks(t *testing.T) {
	t.Parallel()

	card := testModel(&fakeMonitor{}).cards[0]
	rails := metadataRails(card, 120, healthPalette(card.Health))
	require.Len(t, rails, 3)
	for _, rail := range rails {
		require.Equal(t, 120, lipgloss.Width(rail))
	}
	require.Contains(t, rails[0], "APP ID app-alpha")
	require.Contains(t, rails[1], "✓ BUILD VALID")
	require.Contains(t, rails[1], "✓ NO BLOCKERS")
	require.Contains(t, rails[1], "✓ REVIEW INFO")
	require.Contains(t, rails[1], "✓ IN FLIGHT")
}

func TestApprovedMetadataShowsOutcomeWithoutFalseInFlightFailure(t *testing.T) {
	t.Parallel()

	card := terminalCard()
	rails := metadataRails(card, 120, approvalPalette())
	require.Contains(t, rails[0], "APPROVED")
	require.Contains(t, rails[1], "✓ REVIEW COMPLETE")
	require.NotContains(t, rails[1], "✕ IN FLIGHT")
}
