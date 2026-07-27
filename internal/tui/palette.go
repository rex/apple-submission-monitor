package tui

import (
	"fmt"
	"math"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

type palette struct {
	base       string
	bright     string
	background string
	flash      string
	shadow     string
	gradient   []string
}

func healthPalette(health domain.Health) palette {
	switch health {
	case domain.HealthRed:
		return palette{
			base: "#FF4D6D", bright: "#FF8FA3", background: "#250B12",
			flash: "#65152A", shadow: "#4A0615",
			gradient: []string{"#FF8A8A", "#FF174D", "#9E002B", "#5C0018", "#FF8A8A"},
		}
	case domain.HealthYellow:
		return palette{
			base: "#4F8CFF", bright: "#B9DEFF", background: "#020817",
			flash: "#143F82", shadow: "#102A5B",
			gradient: []string{"#25D8FF", "#3F8CFF", "#795CFF", "#E766FF", "#25D8FF"},
		}
	case domain.HealthGreen:
		return palette{
			base: "#5AF78E", bright: "#A7F3D0", background: "#062319",
			flash: "#0C6247", shadow: "#123E31",
			gradient: []string{"#2AF598", "#00D7B9", "#38BDF8", "#B7F34A", "#2AF598"},
		}
	case domain.HealthBlue:
		return palette{
			base: "#60A5FA", bright: "#BFDBFE", background: "#071B35",
			flash: "#124A87", shadow: "#153A70",
			gradient: []string{"#22D3EE", "#3B82F6", "#8B5CF6", "#C4B5FD", "#22D3EE"},
		}
	case domain.HealthPurple:
		return palette{
			base: "#C084FC", bright: "#E9D5FF", background: "#211030",
			flash: "#64338A", shadow: "#512D6D",
			gradient: []string{"#A78BFA", "#D946EF", "#FB7185", "#67E8F9", "#A78BFA"},
		}
	default:
		return palette{
			base: "#8B949E", bright: "#D0D7DE", background: "#11151A",
			flash: "#30363D", shadow: "#30363D",
			gradient: []string{"#8B949E", "#D0D7DE", "#9CA3AF", "#8B949E"},
		}
	}
}

func interpolateHex(start string, end string, position float64) string {
	var startR, startG, startB int
	var endR, endG, endB int
	if _, err := fmt.Sscanf(start, "#%02x%02x%02x", &startR, &startG, &startB); err != nil {
		return start
	}
	if _, err := fmt.Sscanf(end, "#%02x%02x%02x", &endR, &endG, &endB); err != nil {
		return start
	}
	channel := func(a int, b int) int {
		return int(math.Round(float64(a) + (float64(b)-float64(a))*position))
	}
	return fmt.Sprintf("#%02X%02X%02X", channel(startR, endR), channel(startG, endG), channel(startB, endB))
}
