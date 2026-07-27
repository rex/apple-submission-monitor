package tui

import (
	"context"
	"time"

	"github.com/rex/apple-submission-monitor/internal/config"
	"github.com/rex/apple-submission-monitor/internal/domain"
	"github.com/rex/apple-submission-monitor/internal/monitor"
)

type fakeMonitor struct {
	acknowledged  string
	removed       string
	bootResult    monitor.Result
	refreshResult monitor.Result
}

func (f *fakeMonitor) Bootstrap(context.Context) (monitor.Result, error) {
	return f.bootResult, nil
}

func (f *fakeMonitor) Refresh(
	context.Context,
	[]domain.Submission,
	bool,
) (monitor.Result, error) {
	return f.refreshResult, nil
}

func (f *fakeMonitor) Acknowledge(
	cards []domain.Submission,
	key string,
) ([]domain.Submission, error) {
	f.acknowledged = key
	return cards, nil
}

func (f *fakeMonitor) Remove(
	_ []domain.Submission,
	key string,
) ([]domain.Submission, error) {
	f.removed = key
	return nil, nil
}

type fakeSpeaker struct {
	count int
}

func (f *fakeSpeaker) Announce(
	context.Context,
	domain.Submission,
	domain.Submission,
) error {
	f.count++
	return nil
}

type fakeOpener struct {
	value string
}

func (f *fakeOpener) Open(_ context.Context, value string) error {
	f.value = value
	return nil
}

func testModel(engine Monitor) Model {
	card := domain.Submission{
		ID:            "review-alpha",
		AppID:         "app-alpha",
		AppName:       "Synthetic Alpha",
		BundleID:      "test.example.alpha",
		Platform:      "IOS",
		Version:       "1.2.3",
		Health:        domain.HealthYellow,
		ReviewState:   "WAITING_FOR_REVIEW",
		AppStoreState: "WAITING_FOR_REVIEW",
		NextAction:    "Wait for review.",
		ReviewURL:     "https://example.test/review",
		SubmittedAt:   time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		LastChangedAt: time.Now().Add(-2 * time.Hour),
		InFlight:      true,
		BuildState:    "VALID",
		BuildKnown:    true,
		BlockersKnown: true,
		ReviewDetails: true,
		ReviewKnown:   true,
		Acknowledged:  true,
	}
	cfg := config.Config{
		Path:              "synthetic-config.yaml",
		PollInterval:      30 * time.Second,
		DiscoveryInterval: 5 * time.Minute,
		AnimationInterval: 400 * time.Millisecond,
	}
	model := New(context.Background(), engine, &fakeSpeaker{}, &fakeOpener{}, cfg)
	model.cards = []domain.Submission{card}
	model.width = 100
	model.height = 28
	model.loaded = true
	model.loading = false
	model.lastUpdated = time.Now()
	return model
}

func terminalCard() domain.Submission {
	return domain.Submission{
		ID:            "review-terminal",
		AppName:       "Synthetic Terminal",
		Health:        domain.HealthGreen,
		ReviewState:   "COMPLETE",
		AppStoreState: "READY_FOR_DISTRIBUTION",
		ReviewURL:     "https://example.test/review",
		Acknowledged:  true,
		Retained:      true,
	}
}
