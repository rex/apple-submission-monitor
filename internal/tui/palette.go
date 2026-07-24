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
	gradientA  string
	gradientB  string
}

func healthPalette(health domain.Health) palette {
	switch health {
	case domain.HealthRed:
		return palette{
			base: "#FF4D6D", bright: "#FF8FA3", background: "#250B12",
			flash: "#65152A", gradientA: "#FF4D6D", gradientB: "#FFB3C1",
		}
	case domain.HealthYellow:
		return palette{
			base: "#FFD166", bright: "#FFE29A", background: "#251C08",
			flash: "#6B4D00", gradientA: "#FFB703", gradientB: "#FFF3B0",
		}
	case domain.HealthGreen:
		return palette{
			base: "#5AF78E", bright: "#A7F3D0", background: "#062319",
			flash: "#0C6247", gradientA: "#34D399", gradientB: "#99F6E4",
		}
	case domain.HealthBlue:
		return palette{
			base: "#60A5FA", bright: "#BFDBFE", background: "#071B35",
			flash: "#124A87", gradientA: "#38BDF8", gradientB: "#C4B5FD",
		}
	case domain.HealthPurple:
		return palette{
			base: "#C084FC", bright: "#E9D5FF", background: "#211030",
			flash: "#64338A", gradientA: "#A78BFA", gradientB: "#F0ABFC",
		}
	default:
		return palette{
			base: "#8B949E", bright: "#D0D7DE", background: "#11151A",
			flash: "#30363D", gradientA: "#8B949E", gradientB: "#D0D7DE",
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
