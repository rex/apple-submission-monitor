package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

type metadataItem struct {
	text  string
	color string
	link  string
}

func (item metadataItem) render(width int) string {
	text := truncate(item.text, width)
	if item.link != "" && text == item.text {
		text = hyperlink(item.link, text)
	}
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(item.color)).
		Render(text)
}

func metadataRails(
	card domain.Submission,
	width int,
	colors palette,
) []string {
	version := card.Version
	if version == "" {
		version = "UNKNOWN"
	}
	platform := card.Platform
	if platform == "" {
		platform = "UNKNOWN"
	}
	appID := card.AppID
	if appID == "" {
		appID = "UNKNOWN"
	}

	first := metadataRail(width, []metadataItem{
		{text: "● " + strings.ToUpper(card.StatusLabel()), color: colors.bright},
		{text: "VERSION " + version, color: "#DDE7F5"},
		{text: "PLATFORM " + platform, color: "#DDE7F5"},
		{text: "APP ID " + appID, color: "#DDE7F5"},
	})
	second := metadataRail(width, []metadataItem{
		buildCheck(card),
		booleanCheck("NO BLOCKERS", card.BlockersKnown, card.BlockerCount == 0),
		booleanCheck("REVIEW INFO", card.ReviewKnown, card.ReviewDetails),
		booleanCheck("IN FLIGHT", true, card.InFlight),
	})

	submitted := "SUBMITTED UNKNOWN"
	if !card.SubmittedAt.IsZero() {
		submitted = "SUBMITTED " + strings.ToUpper(
			card.SubmittedAt.Local().Format("Jan 2 3:04 PM"),
		)
	}
	link := card.ReviewURL
	if link == "" {
		link = card.AppStoreURL
	}
	third := metadataRail(width, []metadataItem{
		{text: "BUNDLE " + firstNonEmptyText(card.BundleID, "UNKNOWN"), color: "#A7B4C8"},
		{text: submitted, color: "#A7B4C8"},
		{text: "NEXT " + firstNonEmptyText(card.NextAction, "UNKNOWN"), color: "#A7B4C8"},
		{text: "↗ APP STORE CONNECT [O]", color: "#C4B5FD", link: link},
	})
	return []string{first, second, third}
}

func metadataRail(width int, items []metadataItem) string {
	if width <= 0 || len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		value := items[0].render(width)
		return value + strings.Repeat(" ", max(0, width-lipgloss.Width(value)))
	}
	total := 0
	for _, item := range items {
		total += lipgloss.Width(item.text)
	}
	if total+len(items)-1 <= width {
		remaining := width - total
		gaps := max(1, len(items)-1)
		baseGap := remaining / gaps
		extra := remaining % gaps
		var row strings.Builder
		for index, item := range items {
			row.WriteString(item.render(lipgloss.Width(item.text)))
			if index == len(items)-1 {
				continue
			}
			gap := baseGap
			if index < extra {
				gap++
			}
			row.WriteString(strings.Repeat(" ", gap))
		}
		return row.String()
	}

	gaps := len(items) - 1
	available := max(1, width-gaps)
	slotWidth := max(1, available/len(items))
	extra := max(0, available-slotWidth*len(items))
	rendered := make([]string, len(items))
	for index, item := range items {
		widthForItem := slotWidth
		if index < extra {
			widthForItem++
		}
		value := item.render(widthForItem)
		rendered[index] = value + strings.Repeat(
			" ",
			max(0, widthForItem-lipgloss.Width(value)),
		)
	}
	return strings.Join(rendered, " ")
}

func booleanCheck(label string, known bool, passed bool) metadataItem {
	switch {
	case !known:
		return metadataItem{text: "? " + label, color: "#8B949E"}
	case passed:
		return metadataItem{text: "✓ " + label, color: "#5AF78E"}
	default:
		return metadataItem{text: "✕ " + label, color: "#FF6B81"}
	}
}

func buildCheck(card domain.Submission) metadataItem {
	if !card.BuildKnown {
		return metadataItem{text: "? BUILD UNKNOWN", color: "#8B949E"}
	}
	if card.BuildValid() {
		return metadataItem{text: "✓ BUILD VALID", color: "#5AF78E"}
	}
	return metadataItem{
		text:  fmt.Sprintf("✕ BUILD %s", firstNonEmptyText(card.BuildState, "INVALID")),
		color: "#FF6B81",
	}
}

func firstNonEmptyText(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
