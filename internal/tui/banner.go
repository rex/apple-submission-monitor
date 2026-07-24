package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

func renderBanner(name string, width int, colors palette) []string {
	maxCharacters := max(1, (width+1)/4)
	characters := []rune(strings.ToUpper(strings.TrimSpace(name)))
	if len(characters) > maxCharacters {
		characters = characters[:maxCharacters]
		characters[len(characters)-1] = '~'
	}

	plain := make([]string, 5)
	for _, character := range characters {
		glyph, ok := blockFont[character]
		if !ok {
			glyph = blockFont['?']
		}
		for row := range plain {
			if plain[row] != "" {
				plain[row] += " "
			}
			plain[row] += glyph[row]
		}
	}
	for row, line := range plain {
		plain[row] = gradientLine(runewidth.Truncate(line, width, ""), colors)
	}
	return plain
}

func gradientLine(line string, colors palette) string {
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
		color := interpolateHex(colors.gradientA, colors.gradientB, position)
		rendered.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Bold(true).
			Render(string(character)))
	}
	return rendered.String()
}

var blockFont = map[rune][5]string{
	'A': {"███", "█ █", "███", "█ █", "█ █"},
	'B': {"██ ", "█ █", "██ ", "█ █", "██ "},
	'C': {"███", "█  ", "█  ", "█  ", "███"},
	'D': {"██ ", "█ █", "█ █", "█ █", "██ "},
	'E': {"███", "█  ", "██ ", "█  ", "███"},
	'F': {"███", "█  ", "██ ", "█  ", "█  "},
	'G': {"███", "█  ", "█ █", "█ █", "███"},
	'H': {"█ █", "█ █", "███", "█ █", "█ █"},
	'I': {"███", " █ ", " █ ", " █ ", "███"},
	'J': {"███", "  █", "  █", "█ █", "███"},
	'K': {"█ █", "█ █", "██ ", "█ █", "█ █"},
	'L': {"█  ", "█  ", "█  ", "█  ", "███"},
	'M': {"█ █", "███", "███", "█ █", "█ █"},
	'N': {"█ █", "███", "███", "███", "█ █"},
	'O': {"███", "█ █", "█ █", "█ █", "███"},
	'P': {"███", "█ █", "███", "█  ", "█  "},
	'Q': {"███", "█ █", "█ █", "███", "  █"},
	'R': {"██ ", "█ █", "██ ", "█ █", "█ █"},
	'S': {"███", "█  ", "███", "  █", "███"},
	'T': {"███", " █ ", " █ ", " █ ", " █ "},
	'U': {"█ █", "█ █", "█ █", "█ █", "███"},
	'V': {"█ █", "█ █", "█ █", "█ █", " █ "},
	'W': {"█ █", "█ █", "███", "███", "█ █"},
	'X': {"█ █", "█ █", " █ ", "█ █", "█ █"},
	'Y': {"█ █", "█ █", " █ ", " █ ", " █ "},
	'Z': {"███", "  █", " █ ", "█  ", "███"},
	'0': {"███", "█ █", "█ █", "█ █", "███"},
	'1': {" ██", "  █", "  █", "  █", "███"},
	'2': {"███", "  █", "███", "█  ", "███"},
	'3': {"███", "  █", " ██", "  █", "███"},
	'4': {"█ █", "█ █", "███", "  █", "  █"},
	'5': {"███", "█  ", "███", "  █", "███"},
	'6': {"███", "█  ", "███", "█ █", "███"},
	'7': {"███", "  █", " █ ", " █ ", " █ "},
	'8': {"███", "█ █", "███", "█ █", "███"},
	'9': {"███", "█ █", "███", "  █", "███"},
	' ': {"   ", "   ", "   ", "   ", "   "},
	'-': {"   ", "   ", "███", "   ", "   "},
	'.': {"   ", "   ", "   ", "   ", " █ "},
	'&': {"██ ", "█ █", " ██", "█ █", " ██"},
	'~': {"   ", "█ █", " █ ", "   ", "   "},
	'?': {"███", "  █", " ██", "   ", " █ "},
}
