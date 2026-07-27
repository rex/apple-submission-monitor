package tui

import (
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

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
	maxHeight := max(1, height-4)
	cached, ok := m.bannerCache[card.Key()]
	if ok &&
		cached.spec.name == card.AppName &&
		cached.spec.health == card.Health &&
		cached.spec.width == width &&
		cached.spec.maxHeight == maxHeight {
		return cached.art.render(
			width,
			colors,
			submissionAgePalette(cached.spec.ageBand),
			frame,
		)
	}
	now := time.Now()
	age := submissionAge(card, now)
	ageBand := classifySubmissionAge(age)
	return prepareHero(
		card.AppName,
		submissionAgeText(card, now),
		width,
		maxHeight,
		card.Health == domain.HealthRed,
		ageBand,
	).render(width, colors, submissionAgePalette(ageBand), frame)
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

func gradientPosition(column int, row int, width int, frame uint64) float64 {
	travel := float64(frame%120) / 120
	position := float64(column)/float64(max(1, width-1)) + float64(row)*0.035 + travel
	return math.Mod(position, 1)
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
