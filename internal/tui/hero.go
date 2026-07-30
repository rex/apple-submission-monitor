package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

const heroGap = 3

type heroArt struct {
	name       bannerArt
	right      bannerArt
	nameWidth  int
	rightWidth int
	font       string
}

func prepareHero(
	name string,
	rightText string,
	width int,
	maxHeight int,
	nameRejected bool,
	rightBloody bool,
) heroArt {
	for _, candidate := range candidateFonts(nameRejected) {
		nameFigure := renderFigure(name, candidate.name)
		rightFigure := renderFigure(rightText, candidate.name)
		if len(nameFigure) == 0 || len(rightFigure) == 0 {
			continue
		}
		nameArt := styleFigure(nameFigure, name, candidate.bloody)
		rightArt := styleFigure(rightFigure, rightText, rightBloody)
		nameWidth := figureWidth(nameArt.figure)
		rightWidth := figureWidth(rightArt.figure)
		if nameWidth+heroGap+rightWidth > width ||
			max(len(nameArt.figure), len(rightArt.figure)) > maxHeight {
			continue
		}
		return heroArt{
			name: nameArt, right: rightArt,
			nameWidth: nameWidth, rightWidth: rightWidth,
			font: candidate.name,
		}
	}

	nameFallback := truncate(strings.ToUpper(strings.TrimSpace(name)), max(1, width/2))
	rightFallback := truncate(rightText, max(1, width/3))
	return heroArt{
		name:       bannerArt{fallback: nameFallback},
		right:      bannerArt{fallback: rightFallback},
		nameWidth:  runewidth.StringWidth(nameFallback),
		rightWidth: runewidth.StringWidth(rightFallback),
	}
}

func (art heroArt) render(
	width int,
	nameColors palette,
	rightColors palette,
	frame uint64,
) []string {
	nameLines := art.name.render(max(1, art.nameWidth), nameColors, frame)
	rightLines := art.right.render(max(1, art.rightWidth), rightColors, frame+40)
	height := max(len(nameLines), len(rightLines))
	lines := make([]string, height)
	for row := range height {
		left := ""
		if row < len(nameLines) {
			left = nameLines[row]
		}
		right := ""
		if row < len(rightLines) {
			right = rightLines[row]
		}
		lines[row] = padBetween(left, right, width)
		if lipgloss.Width(lines[row]) > width {
			lines[row] = left
		}
	}
	return lines
}
