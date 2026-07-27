package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

const heroGap = 3

type heroArt struct {
	name      bannerArt
	age       bannerArt
	nameWidth int
	ageWidth  int
	font      string
}

func prepareHero(
	name string,
	ageText string,
	width int,
	maxHeight int,
	nameRejected bool,
	ageBand submissionAgeBand,
) heroArt {
	for _, candidate := range candidateFonts(nameRejected) {
		nameFigure := renderFigure(name, candidate.name)
		ageFigure := renderFigure(ageText, candidate.name)
		if len(nameFigure) == 0 || len(ageFigure) == 0 {
			continue
		}
		nameArt := styleFigure(nameFigure, name, candidate.bloody)
		ageArt := styleFigure(ageFigure, ageText, ageBand == ageBandRed)
		nameWidth := figureWidth(nameArt.figure)
		ageWidth := figureWidth(ageArt.figure)
		if nameWidth+heroGap+ageWidth > width ||
			max(len(nameArt.figure), len(ageArt.figure)) > maxHeight {
			continue
		}
		return heroArt{
			name: nameArt, age: ageArt,
			nameWidth: nameWidth, ageWidth: ageWidth,
			font: candidate.name,
		}
	}

	nameFallback := truncate(strings.ToUpper(strings.TrimSpace(name)), max(1, width/2))
	ageFallback := truncate(ageText, max(1, width/3))
	return heroArt{
		name:      bannerArt{fallback: nameFallback},
		age:       bannerArt{fallback: ageFallback},
		nameWidth: runewidth.StringWidth(nameFallback),
		ageWidth:  runewidth.StringWidth(ageFallback),
	}
}

func (art heroArt) render(
	width int,
	nameColors palette,
	ageColors palette,
	frame uint64,
) []string {
	nameLines := art.name.render(max(1, art.nameWidth), nameColors, frame)
	ageLines := art.age.render(max(1, art.ageWidth), ageColors, frame+40)
	height := max(len(nameLines), len(ageLines))
	lines := make([]string, height)
	for row := range height {
		left := ""
		if row < len(nameLines) {
			left = nameLines[row]
		}
		right := ""
		if row < len(ageLines) {
			right = ageLines[row]
		}
		lines[row] = padBetween(left, right, width)
		if lipgloss.Width(lines[row]) > width {
			lines[row] = left
		}
	}
	return lines
}
