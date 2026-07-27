package tui

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

func elapsed(value time.Time, now time.Time) string {
	if value.IsZero() {
		return "unknown"
	}
	duration := now.Sub(value)
	if duration < 0 {
		duration = 0
	}
	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		return fmt.Sprintf("%dm", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%dh %dm", int(duration.Hours()), int(duration.Minutes())%60)
	default:
		return fmt.Sprintf("%dd %dh", int(duration.Hours()/24), int(duration.Hours())%24)
	}
}

func truncate(value string, width int) string {
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "\r", " ")
	value = strings.Map(func(character rune) rune {
		if unicode.IsControl(character) {
			return -1
		}
		return character
	}, value)
	if width <= 0 {
		return ""
	}
	return runewidth.Truncate(value, width, "…")
}

func padBetween(left string, right string, width int) string {
	remaining := width - lipgloss.Width(left) - lipgloss.Width(right)
	if remaining < 1 {
		return truncate(left, width)
	}
	return left + strings.Repeat(" ", remaining) + right
}
