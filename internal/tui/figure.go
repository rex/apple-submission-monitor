package tui

import (
	"strings"
	"unicode"

	"github.com/lsferreira42/figlet-go/figlet"
	"github.com/mattn/go-runewidth"
)

type bannerFont struct {
	name   string
	bloody bool
}

func prepareBanner(name string, width int, maxHeight int, rejected bool) bannerArt {
	figure, style := fittedFigure(name, width, maxHeight, rejected)
	if len(figure) == 0 {
		return bannerArt{fallback: truncate(strings.ToUpper(strings.TrimSpace(name)), width)}
	}
	if style.bloody {
		figure = distortFigure(figure)
		figure = appendBloodDrips(figure, width, name)
		return bannerArt{figure: figure, bloody: true}
	}
	return bannerArt{figure: withShadow(figure)}
}

func fittedFigure(name string, width int, maxHeight int, rejected bool) ([]string, bannerFont) {
	fonts := []bannerFont{
		{name: "banner3"},
		{name: "doom"},
		{name: "standard"},
		{name: "smshadow"},
		{name: "mini"},
	}
	if rejected {
		fonts = []bannerFont{
			{name: "poison", bloody: true},
			{name: "sblood", bloody: true},
			{name: "banner3", bloody: true},
			{name: "mini", bloody: true},
		}
	}

	cleanName := cleanFigureName(name)
	for _, candidate := range fonts {
		figure, err := figlet.Render(
			cleanName,
			figlet.WithFont(candidate.name),
			figlet.WithWidth(1000),
		)
		if err != nil {
			continue
		}
		lines := normalizeFigure(figure)
		effectHeight := 1
		effectWidth := 1
		if candidate.bloody {
			effectHeight = 3
			effectWidth = 2
		}
		if figureFits(lines, width-effectWidth, maxHeight-effectHeight) {
			return lines, candidate
		}
	}
	return nil, bannerFont{}
}

func cleanFigureName(name string) string {
	name = strings.ToUpper(strings.TrimSpace(name))
	name = strings.Map(func(character rune) rune {
		switch {
		case character >= 32 && character <= 126:
			return character
		case unicode.IsSpace(character):
			return ' '
		default:
			return '?'
		}
	}, name)
	return runewidth.Truncate(name, 64, "")
}

func normalizeFigure(figure string) []string {
	lines := strings.Split(strings.TrimRight(figure, "\n"), "\n")
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	for index := range lines {
		lines[index] = strings.TrimRight(lines[index], " ")
	}
	return lines
}

func figureFits(lines []string, width int, height int) bool {
	if len(lines) == 0 || len(lines) > height || width < 1 {
		return false
	}
	for _, line := range lines {
		if runewidth.StringWidth(line) > width {
			return false
		}
	}
	return true
}

func withShadow(figure []string) []string {
	width := figureWidth(figure)
	result := make([]string, len(figure)+1)
	for row := range result {
		var line strings.Builder
		for column := 0; column <= width; column++ {
			face := figureInk(figure, row, column)
			shadow := figureInk(figure, row-1, column-1)
			switch {
			case face:
				line.WriteRune('█')
			case shadow:
				line.WriteRune('▓')
			default:
				line.WriteRune(' ')
			}
		}
		result[row] = strings.TrimRight(line.String(), " ")
	}
	return result
}

func distortFigure(figure []string) []string {
	result := make([]string, len(figure))
	for row, line := range figure {
		offset := []int{1, 0, 2, 0}[row%4]
		result[row] = strings.Repeat(" ", offset) + line
	}
	return result
}

func appendBloodDrips(figure []string, width int, name string) []string {
	drips := make([][]rune, 3)
	for row := range drips {
		drips[row] = make([]rune, width)
		for column := range drips[row] {
			drips[row][column] = ' '
		}
	}

	seed := 0
	for _, character := range name {
		seed += int(character)
	}
	dripCount := 0
	for column := 0; column < min(width, figureWidth(figure)); column++ {
		if !columnHasInk(figure, column) || (column+seed)%11 != 0 {
			continue
		}
		length := 1 + (column+seed)%3
		for row := 0; row < length && row < len(drips); row++ {
			if column >= len(drips[row]) {
				continue
			}
			drips[row][column] = '│'
		}
		lastRow := min(length-1, len(drips)-1)
		if column < len(drips[lastRow]) {
			drips[lastRow][column] = '●'
		}
		dripCount++
	}
	if dripCount == 0 {
		for column := 0; column < min(width, figureWidth(figure)); column++ {
			if columnHasInk(figure, column) {
				drips[0][column] = '●'
				break
			}
		}
	}

	result := append([]string(nil), figure...)
	for _, drip := range drips {
		result = append(result, strings.TrimRight(string(drip), " "))
	}
	return result
}

func figureWidth(figure []string) int {
	width := 0
	for _, line := range figure {
		width = max(width, len([]rune(line)))
	}
	return width
}

func figureInk(figure []string, row int, column int) bool {
	if row < 0 || row >= len(figure) || column < 0 {
		return false
	}
	line := []rune(figure[row])
	return column < len(line) && line[column] != ' '
}

func columnHasInk(figure []string, column int) bool {
	for _, line := range figure {
		characters := []rune(line)
		if column < len(characters) && characters[column] != ' ' {
			return true
		}
	}
	return false
}
