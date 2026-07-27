package tui

import (
	"fmt"
	"time"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

type submissionAgeBand uint8

const (
	ageBandGreen submissionAgeBand = iota
	ageBandYellow
	ageBandRed
)

func submissionAge(card domain.Submission, now time.Time) time.Duration {
	if card.SubmittedAt.IsZero() {
		return 0
	}
	age := now.Sub(card.SubmittedAt)
	return max(age, 0)
}

func submissionAgeText(card domain.Submission, now time.Time) string {
	age := submissionAge(card, now)
	if age >= 24*time.Hour {
		return fmt.Sprintf("%dD %dH", int(age.Hours()/24), int(age.Hours())%24)
	}
	return fmt.Sprintf("%dH %dM", int(age.Hours()), int(age.Minutes())%60)
}

func classifySubmissionAge(age time.Duration) submissionAgeBand {
	switch {
	case age <= 48*time.Hour:
		return ageBandGreen
	case age < 120*time.Hour:
		return ageBandYellow
	default:
		return ageBandRed
	}
}

func submissionAgePalette(band submissionAgeBand) palette {
	switch band {
	case ageBandYellow:
		return palette{
			base: "#FFC857", bright: "#FFF0A6", shadow: "#6B3F12",
			gradient: []string{"#FFF0A6", "#FFD166", "#FF9F1C", "#FFCF70", "#FFF0A6"},
		}
	case ageBandRed:
		return palette{
			base: "#FF174D", bright: "#FF8A8A", shadow: "#4A0615",
			gradient: []string{"#FF8A8A", "#FF174D", "#9E002B", "#5C0018", "#FF8A8A"},
		}
	default:
		return palette{
			base: "#2AF598", bright: "#B7FBD8", shadow: "#123E31",
			gradient: []string{"#B7FBD8", "#5AF78E", "#00D7B9", "#B7F34A", "#B7FBD8"},
		}
	}
}
