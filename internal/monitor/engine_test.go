package monitor_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/rex/apple-submission-monitor/internal/domain"
	"github.com/rex/apple-submission-monitor/internal/monitor"
)

type fakeSource struct {
	discovered []domain.Submission
	refreshed  []domain.Submission
}

func (f *fakeSource) Discover(context.Context) ([]domain.Submission, error) {
	return append([]domain.Submission(nil), f.discovered...), nil
}

func (f *fakeSource) Refresh(
	context.Context,
	[]domain.Submission,
) ([]domain.Submission, error) {
	return append([]domain.Submission(nil), f.refreshed...), nil
}

type memoryStore struct {
	cards []domain.Submission
	saves int
}

func (s *memoryStore) Load() ([]domain.Submission, error) {
	return append([]domain.Submission(nil), s.cards...), nil
}

func (s *memoryStore) Save(cards []domain.Submission) error {
	s.cards = append([]domain.Submission(nil), cards...)
	s.saves++
	return nil
}

func TestEngineLifecycle(t *testing.T) {
	t.Parallel()

	waiting := syntheticWaiting()
	source := &fakeSource{discovered: []domain.Submission{waiting}}
	store := &memoryStore{}
	engine := monitor.NewEngine(source, store)

	initial, err := engine.Bootstrap(context.Background())
	require.NoError(t, err)
	require.Empty(t, initial.Changes)
	require.True(t, initial.Cards[0].Acknowledged)

	approved := waiting
	approved.Health = domain.HealthGreen
	approved.ReviewState = "COMPLETE"
	approved.AppStoreState = "READY_FOR_DISTRIBUTION"
	approved.Outcome = "APPROVED"
	approved.InFlight = false
	source.refreshed = []domain.Submission{approved}
	source.discovered = nil

	changed, err := engine.Refresh(context.Background(), initial.Cards, true)
	require.NoError(t, err)
	require.Len(t, changed.Changes, 1)
	require.False(t, changed.Cards[0].Acknowledged)
	require.True(t, changed.Cards[0].Retained)

	acknowledged, err := engine.Acknowledge(changed.Cards, waiting.Key())
	require.NoError(t, err)
	require.True(t, acknowledged[0].Acknowledged)

	removed, err := engine.Remove(acknowledged, waiting.Key())
	require.NoError(t, err)
	require.Empty(t, removed)
	require.GreaterOrEqual(t, store.saves, 4)
}

func TestEngineRetainsMissingKnownCard(t *testing.T) {
	t.Parallel()

	waiting := syntheticWaiting()
	source := &fakeSource{}
	store := &memoryStore{cards: []domain.Submission{waiting}}
	engine := monitor.NewEngine(source, store)

	result, err := engine.Bootstrap(context.Background())
	require.NoError(t, err)
	require.Equal(t, []domain.Submission{waiting}, result.Cards)
}

func TestRemoveRequiresAcknowledgedOutcome(t *testing.T) {
	t.Parallel()

	waiting := syntheticWaiting()
	engine := monitor.NewEngine(&fakeSource{}, &memoryStore{})

	_, err := engine.Remove([]domain.Submission{waiting}, waiting.Key())
	require.ErrorContains(t, err, "acknowledge terminal card")
}

func syntheticWaiting() domain.Submission {
	submitted := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	return domain.Submission{
		ID:            "review-alpha",
		AppID:         "app-alpha",
		AppName:       "Synthetic Alpha",
		Platform:      "IOS",
		Health:        domain.HealthYellow,
		ReviewState:   "WAITING_FOR_REVIEW",
		AppStoreState: "WAITING_FOR_REVIEW",
		SubmittedAt:   submitted,
		InFlight:      true,
		Acknowledged:  true,
	}
}
