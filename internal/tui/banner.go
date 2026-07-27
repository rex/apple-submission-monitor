package tui

import (
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

func statusAgeColors() []string {
	return []string{
		"#B6F7D4",
		"#4EE0B7",
		"#FFAD9E",
		"#F5A7D8",
		"#B6F7D4",
	}
}

type bannerArt struct {
	figure   []string
	bloody   bool
	fallback string
}

func renderBanner(
	name string,
	width int,
	maxHeight int,
	colors palette,
	frame uint64,
	rejected bool,
) []string {
	return prepareBanner(name, width, maxHeight, rejected).render(width, colors, frame)
}

func (art bannerArt) render(width int, colors palette, frame uint64) []string {
	if art.fallback != "" {
		return []string{gradientLine(art.fallback, colors, frame)}
	}
	return colorizeFigure(art.figure, width, colors, frame, art.bloody)
}

func (m Model) renderCardBanner(
	card domain.Submission,
	width int,
	height int,
	colors palette,
	frame uint64,
) []string {
	maxHeight := max(1, height-5)
	cached, ok := m.bannerCache[card.Key()]
	if ok && cached.spec == (bannerSpec{
		name:      card.AppName,
		health:    card.Health,
		width:     width,
		maxHeight: maxHeight,
	}) {
		return cached.art.render(width, colors, frame)
	}
	return renderBanner(
		card.AppName,
		width,
		maxHeight,
		colors,
		frame,
		card.Health == domain.HealthRed,
	)
}

func colorizeFigure(
	figure []string,
	width int,
	colors palette,
	frame uint64,
	bloody bool,
) []string {
	glyphs := newGradientGlyphs(colors)
	shadow := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.shadow)).
		Bold(true).
		Render("▓")
	result := make([]string, len(figure))
	for row, line := range figure {
		var rendered strings.Builder
		for column, character := range []rune(line) {
			if character == ' ' {
				rendered.WriteRune(character)
				continue
			}
			position := gradientPosition(column, row, width, frame)
			band := min(int(position*gradientBands), gradientBands-1)
			if character == '▓' {
				rendered.WriteString(shadow)
				continue
			}
			switch {
			case bloody && character == '│':
				rendered.WriteString(glyphs.drip[band])
			case bloody && character == '●':
				rendered.WriteString(glyphs.drop[band])
			default:
				rendered.WriteString(glyphs.block[band])
			}
		}
		result[row] = rendered.String()
	}
	return result
}

const gradientBands = 32

type gradientGlyphs struct {
	block [gradientBands]string
	drip  [gradientBands]string
	drop  [gradientBands]string
}

func newGradientGlyphs(colors palette) gradientGlyphs {
	var glyphs gradientGlyphs
	for index := range gradientBands {
		position := float64(index) / float64(gradientBands-1)
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color(gradientColor(colors.gradient, position))).
			Bold(true)
		glyphs.block[index] = style.Render("█")
		glyphs.drip[index] = style.Render("│")
		glyphs.drop[index] = style.Render("●")
	}
	return glyphs
}

func gradientLine(line string, colors palette, frame uint64) string {
	return colorizeFigure([]string{line}, max(1, runewidth.StringWidth(line)), colors, frame, true)[0]
}

func gradientText(line string, stops []string, frame uint64) string {
	var rendered strings.Builder
	width := max(1, runewidth.StringWidth(line)-1)
	column := 0
	for _, character := range line {
		if character == ' ' {
			rendered.WriteRune(character)
			column++
			continue
		}
		position := statusAgeGradientPosition(column, width, frame)
		rendered.WriteString(lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(gradientColor(stops, position))).
			Render(string(character)))
		column += max(1, runewidth.RuneWidth(character))
	}
	return rendered.String()
}

func gradientPosition(column int, row int, width int, frame uint64) float64 {
	travel := float64(frame%120) / 120
	position := float64(column)/float64(max(1, width-1)) + float64(row)*0.035 + travel
	return math.Mod(position, 1)
}

func statusAgeGradientPosition(column int, width int, frame uint64) float64 {
	travel := float64(frame%180) / 180
	position := 1 - float64(column)/float64(max(1, width)) + travel
	return math.Mod(position, 1)
}

func renderStatusRail(
	card domain.Submission,
	width int,
	colors palette,
	frame uint64,
	now time.Time,
) string {
	detailedAge := statusAgeText(card, now, true)
	compactAge := statusAgeText(card, now, false)
	age := detailedAge
	if lipgloss.Width(age) > max(8, width/2) {
		age = compactAge
	}
	age = truncate(age, max(1, width-4))

	badgeWidth := max(1, width-lipgloss.Width(age)-1)
	badge := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colors.base)).
		Render(truncate("● "+strings.ToUpper(card.StatusLabel()), badgeWidth))
	return padBetween(badge, gradientText(age, statusAgeColors(), frame), width)
}

func statusAgeText(card domain.Submission, now time.Time, detailed bool) string {
	started := card.LastChangedAt
	label := "IN STATUS"
	if started.IsZero() {
		started = card.SubmittedAt
		label = "IN REVIEW"
	}
	if started.IsZero() {
		return ""
	}
	age := strings.ToUpper(elapsed(started, now))
	if !detailed {
		return "◷ " + age
	}
	return "◷ " + label + " · " + age
}

func gradientColor(stops []string, position float64) string {
	if len(stops) == 0 {
		return "#FFFFFF"
	}
	if len(stops) == 1 {
		return stops[0]
	}
	scaled := position * float64(len(stops)-1)
	index := min(int(scaled), len(stops)-2)
	return interpolateHex(stops[index], stops[index+1], scaled-float64(index))
}
