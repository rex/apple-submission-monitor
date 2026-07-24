package tui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculateGridUsesWholeViewportEvenly(t *testing.T) {
	t.Parallel()

	layout := calculateGrid(121, 39, 5)
	require.Len(t, layout.rects, 5)
	require.GreaterOrEqual(t, layout.columns, 1)
	require.GreaterOrEqual(t, layout.rows, 1)

	minimumWidth, maximumWidth := layout.rects[0].width, layout.rects[0].width
	minimumHeight, maximumHeight := layout.rects[0].height, layout.rects[0].height
	for index, area := range layout.rects {
		minimumWidth = min(minimumWidth, area.width)
		maximumWidth = max(maximumWidth, area.width)
		minimumHeight = min(minimumHeight, area.height)
		maximumHeight = max(maximumHeight, area.height)
		require.Equal(t, index, layout.hit(area.x, area.y))
		require.LessOrEqual(t, area.x+area.width, 121)
		require.LessOrEqual(t, area.y+area.height, 39-footerHeight)
	}
	require.LessOrEqual(t, maximumWidth-minimumWidth, 1)
	require.LessOrEqual(t, maximumHeight-minimumHeight, 1)
	require.Equal(t, -1, layout.hit(-1, -1))
}

func TestCalculateGridHandlesEmptyAndTinyInput(t *testing.T) {
	t.Parallel()

	require.Empty(t, calculateGrid(80, 24, 0).rects)
	layout := calculateGrid(1, 1, 1)
	require.Len(t, layout.rects, 1)
	require.Equal(t, 1, layout.rects[0].width)
}
