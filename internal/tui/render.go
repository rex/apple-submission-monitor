package tui

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

// View renders the current terminal dashboard.
func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return "Starting Apple Submission Monitor…"
	}
	if m.width < 36 || m.height < 9 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD166")).
			Render("Terminal too small — resize to at least 36×9")
	}

	sections := []string{
		m.renderHeader(),
		m.renderCards(),
		m.renderFooter(),
	}
	return strings.Join(sections, "\n")
}

func (m Model) renderHeader() string {
	updated := "starting"
	if !m.lastUpdated.IsZero() {
		updated = elapsed(m.lastUpdated, time.Now())
		if updated != "just now" {
			updated += " ago"
		}
	}
	liveColor := "#5AF78E"
	if m.loading {
		liveColor = "#FFD166"
	}
	title := lipgloss.NewStyle().Bold(true).
		Foreground(lipgloss.Color("#E6EDF3")).
		Render(" APPLE SUBMISSION MONITOR ")
	live := lipgloss.NewStyle().Bold(true).
		Foreground(lipgloss.Color(liveColor)).
		Render("● LIVE")
	lineOne := padBetween(title, live+" • "+updated+" ", m.width)

	state := fmt.Sprintf("%d active or retained • poll %s", len(m.cards), m.pollInterval)
	message := m.notice
	if message == "" && m.loading {
		message = "Refreshing through asc…"
	}
	lineTwo := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8B949E")).
		Render(padBetween(" "+state, message+" ", m.width))
	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#30363D")).
		Render(strings.Repeat("─", m.width))
	return strings.Join([]string{lineOne, lineTwo, divider}, "\n")
}

func (m Model) renderCards() string {
	layout := calculateGrid(m.width, m.height, len(m.cards))
	if len(m.cards) == 0 {
		height := max(1, m.height-headerHeight-footerHeight)
		message := "No active App Store review submissions"
		detail := "New submissions will appear automatically."
		if !m.loaded {
			message = "Loading active submissions through asc…"
			detail = "The first discovery can take a moment."
		}
		return lipgloss.Place(
			m.width,
			height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.NewStyle().Bold(true).
				Foreground(lipgloss.Color("#8B949E")).
				Render(message+"\n"+detail),
		)
	}

	rows := make([]string, 0, layout.rows)
	for row := range layout.rows {
		cells := make([]string, 0, layout.columns)
		for column := range layout.columns {
			index := row*layout.columns + column
			area := layout.cell(index)
			if index < len(m.cards) {
				cells = append(cells, m.renderCard(m.cards[index], area, index == m.selected))
			} else {
				cells = append(cells, lipgloss.NewStyle().
					Width(max(0, area.width)).
					Height(max(0, area.height)).
					Render(""))
			}
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m Model) renderCard(card domain.Submission, area rect, selected bool) string {
	innerWidth := max(1, area.width-2)
	innerHeight := max(1, area.height-2)
	colors := healthPalette(card.Health)
	borderColor := colors.base
	background := colors.background
	border := lipgloss.RoundedBorder()
	if selected {
		border = lipgloss.DoubleBorder()
		borderColor = colors.bright
	}
	if !card.Acknowledged && m.pulse {
		borderColor = "#FFFFFF"
		background = colors.flash
	}

	lines := m.cardLines(card, innerWidth, innerHeight, colors)
	body := strings.Join(lines, "\n")
	return lipgloss.NewStyle().
		Width(innerWidth).
		Height(innerHeight).
		Border(border).
		BorderForeground(lipgloss.Color(borderColor)).
		Background(lipgloss.Color(background)).
		Foreground(lipgloss.Color("#E6EDF3")).
		Render(body)
}

func (m Model) cardLines(
	card domain.Submission,
	width int,
	height int,
	colors palette,
) []string {
	var lines []string
	if height >= 12 && width >= 20 {
		lines = append(lines, renderBanner(card.AppName, width, colors)...)
	} else {
		lines = append(lines, lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colors.bright)).
			Render(truncate(card.AppName, width)))
	}

	badge := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colors.base)).
		Render("● " + strings.ToUpper(card.StatusLabel()))
	version := strings.TrimSpace(card.Version + " • " + card.Platform)
	lines = append(lines, padBetween(badge, truncate(version, width/2), width))
	lines = append(lines, truncate("Review: "+humanState(card.ReviewState), width))

	if !card.SubmittedAt.IsZero() {
		lines = append(lines, truncate(
			"Submitted: "+card.SubmittedAt.Local().Format("Jan 2, 3:04 PM")+
				" • "+elapsed(card.SubmittedAt, time.Now()),
			width,
		))
	}
	if card.NextAction != "" {
		lines = append(lines, truncate("Next: "+card.NextAction, width))
	}
	if card.BundleID != "" && len(lines) < height-2 {
		lines = append(lines, truncate("Bundle: "+card.BundleID, width))
	}
	link := card.ReviewURL
	if link == "" {
		link = card.AppStoreURL
	}
	if link != "" && len(lines) < height-1 {
		lines = append(lines, hyperlink(link, "↗ App Store Connect")+"  [o]")
	}

	alert := ""
	switch {
	case !card.Acknowledged:
		alert = "⚡ CHANGED — click or Enter to acknowledge"
	case card.Retained:
		alert = "Outcome retained • D removes this card"
	}
	if alert != "" {
		for len(lines) >= height {
			lines = lines[:len(lines)-1]
		}
		lines = append(lines, lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colors.bright)).
			Render(truncate(alert, width)))
	}
	if len(lines) > height {
		lines = lines[:height]
	}
	return lines
}

func (m Model) renderFooter() string {
	help := "←↑↓→/hjkl select  Enter/click acknowledge  o open  D remove  r refresh  ? help  q quit"
	if m.help {
		help = "Cards flash until acknowledged • completed cards persist until D • all data comes from asc"
	}
	path := "Config: " + m.configPath
	return strings.Join([]string{
		lipgloss.NewStyle().Foreground(lipgloss.Color("#8B949E")).
			Render(" " + truncate(help, m.width-1)),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#484F58")).
			Render(" " + truncate(path, m.width-1)),
	}, "\n")
}

func hyperlink(value string, label string) string {
	parsed, err := url.ParseRequestURI(value)
	if err != nil ||
		parsed.Scheme != "https" ||
		parsed.Host == "" ||
		strings.ContainsAny(value, "\x1b\x07\r\n") {
		return label
	}
	return "\x1b]8;;" + value + "\x07" + label + "\x1b]8;;\x07"
}

func humanState(value string) string {
	if value == "" {
		return "Unknown"
	}
	words := strings.Fields(strings.ReplaceAll(strings.ToLower(value), "_", " "))
	for index := range words {
		words[index] = strings.ToUpper(words[index][:1]) + words[index][1:]
	}
	return strings.Join(words, " ")
}
