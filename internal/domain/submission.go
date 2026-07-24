// Package domain defines the monitor's framework-independent data model.
package domain

import (
	"fmt"
	"strings"
	"time"
)

// Health is the normalized color name returned by asc status.
type Health string

const (
	// HealthRed indicates rejection or a blocking review issue.
	HealthRed Health = "red"
	// HealthYellow indicates waiting or work in progress.
	HealthYellow Health = "yellow"
	// HealthGreen indicates a successful terminal state.
	HealthGreen Health = "green"
	// HealthBlue indicates an informational state.
	HealthBlue Health = "blue"
	// HealthPurple indicates a nonstandard upstream informational state.
	HealthPurple Health = "purple"
	// HealthGray is the safe fallback for unknown health values.
	HealthGray Health = "gray"
)

// Submission is one App Store review submission displayed as a card.
type Submission struct {
	ID            string    `json:"id"`
	AppID         string    `json:"app_id"`
	AppName       string    `json:"app_name"`
	BundleID      string    `json:"bundle_id"`
	Platform      string    `json:"platform"`
	Version       string    `json:"version"`
	VersionID     string    `json:"version_id"`
	Health        Health    `json:"health"`
	ReviewState   string    `json:"review_state"`
	AppStoreState string    `json:"app_store_state"`
	NextAction    string    `json:"next_action"`
	AppStoreURL   string    `json:"app_store_url"`
	ReviewURL     string    `json:"review_url"`
	TestFlightURL string    `json:"test_flight_url"`
	SubmittedAt   time.Time `json:"submitted_at"`
	CreatedAt     time.Time `json:"created_at"`
	LastSeenAt    time.Time `json:"last_seen_at"`
	LastChangedAt time.Time `json:"last_changed_at"`
	InFlight      bool      `json:"in_flight"`
	Acknowledged  bool      `json:"acknowledged"`
	Retained      bool      `json:"retained"`
}

// Key returns the stable identity used for persistence and selection.
func (s Submission) Key() string {
	if s.ID != "" {
		return s.ID
	}
	return s.AppID + ":" + s.Platform
}

// Fingerprint returns fields whose change should alert the user.
func (s Submission) Fingerprint() string {
	return strings.Join([]string{
		string(s.Health),
		s.ReviewState,
		s.AppStoreState,
		s.NextAction,
		fmt.Sprintf("%t", s.InFlight),
	}, "\x1f")
}

// StatusLabel returns the most useful concise status for display and speech.
func (s Submission) StatusLabel() string {
	if s.AppStoreState != "" {
		return humanizeState(s.AppStoreState)
	}
	if s.ReviewState != "" {
		return humanizeState(s.ReviewState)
	}
	return "Unknown"
}

// Terminal reports whether the submission has reached an outcome.
func (s Submission) Terminal() bool {
	if s.ReviewState == "COMPLETE" || !s.InFlight && s.Health == HealthGreen {
		return true
	}
	state := strings.ToUpper(s.AppStoreState)
	return strings.Contains(state, "READY_FOR_") ||
		strings.Contains(state, "REJECT") ||
		state == "DEVELOPER_REJECTED"
}

// NormalizeHealth converts asc color names to the supported palette.
func NormalizeHealth(value string) Health {
	switch Health(strings.ToLower(strings.TrimSpace(value))) {
	case HealthRed:
		return HealthRed
	case HealthYellow:
		return HealthYellow
	case HealthGreen:
		return HealthGreen
	case HealthBlue:
		return HealthBlue
	case HealthPurple:
		return HealthPurple
	default:
		return HealthGray
	}
}

func humanizeState(value string) string {
	words := strings.Fields(strings.ReplaceAll(strings.ToLower(value), "_", " "))
	for index := range words {
		words[index] = strings.ToUpper(words[index][:1]) + words[index][1:]
	}
	return strings.Join(words, " ")
}
