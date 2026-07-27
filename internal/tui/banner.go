package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

func renderBanner(name string, width int, colors palette, phase bool) []string {
	wordmark := runewidth.Truncate(
		strings.ToUpper(strings.TrimSpace(name)),
		max(1, width-2),
		"…",
	)
	face := gradientLine("◆ "+wordmark, colors, phase)
	shadow := lipgloss.NewStyle().Foreground(lipgloss.Color(colors.shadow)).
		Render("  " + wordmark)
	ruleWidth := min(max(1, runewidth.StringWidth(wordmark)+1), max(1, width-2))
	rule := lipgloss.NewStyle().Foreground(lipgloss.Color(colors.base)).
		Render("╰" + strings.Repeat("─", ruleWidth) + "╯")
	return []string{face, shadow, rule}
}

func gradientLine(line string, colors palette, phase bool) string {
	var rendered strings.Builder
	characters := []rune(line)
	for index, character := range characters {
		if character == ' ' {
			rendered.WriteRune(character)
			continue
		}
		position := 0.0
		if len(characters) > 1 {
			position = float64(index) / float64(len(characters)-1)
		}
		if phase {
			position = 1 - position
		}
		color := interpolateHex(colors.gradientA, colors.gradientB, position)
		rendered.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Bold(true).
			Render(string(character)))
	}
	return rendered.String()
}
