// Package monitor reconciles asc snapshots with durable user state.
package monitor

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

// Source discovers active submissions and refreshes known cards.
type Source interface {
	Discover(context.Context) ([]domain.Submission, error)
	Refresh(context.Context, []domain.Submission) ([]domain.Submission, error)
}

// Store persists the monitor's durable card state.
type Store interface {
	Load() ([]domain.Submission, error)
	Save([]domain.Submission) error
}

// Change is a meaningful transition that should animate and announce.
type Change struct {
	Before domain.Submission
	After  domain.Submission
}

// Result contains the latest cards, transitions, and a recoverable source error.
type Result struct {
	Cards   []domain.Submission
	Changes []Change
	Err     error
}

// Engine coordinates source refreshes and durable state.
type Engine struct {
	source Source
	store  Store
	now    func() time.Time
}

// NewEngine creates a monitor engine.
func NewEngine(source Source, store Store) *Engine {
	return &Engine{source: source, store: store, now: time.Now}
}

// Bootstrap merges persisted cards, refreshed known cards, and new discoveries.
func (e *Engine) Bootstrap(ctx context.Context) (Result, error) {
	persisted, err := e.store.Load()
	if err != nil {
		return Result{}, err
	}
	result := Result{Cards: persisted}
	if len(persisted) > 0 {
		refreshed, refreshErr := e.source.Refresh(ctx, persisted)
		result = e.reconcile(result.Cards, refreshed)
		result.Err = errors.Join(result.Err, refreshErr)
	}
	discovered, discoveryErr := e.source.Discover(ctx)
	discoveryResult := e.reconcile(result.Cards, discovered)
	result.Cards = discoveryResult.Cards
	result.Changes = append(result.Changes, discoveryResult.Changes...)
	result.Err = errors.Join(result.Err, discoveryErr)
	if err := e.store.Save(result.Cards); err != nil {
		return Result{}, err
	}
	return result, nil
}

// Refresh polls known cards and optionally discovers newly submitted apps.
func (e *Engine) Refresh(
	ctx context.Context,
	current []domain.Submission,
	discover bool,
) (Result, error) {
	observed, refreshErr := e.source.Refresh(ctx, current)
	result := e.reconcile(current, observed)
	result.Err = refreshErr
	if discover {
		discovered, discoveryErr := e.source.Discover(ctx)
		discoveryResult := e.reconcile(result.Cards, discovered)
		result.Cards = discoveryResult.Cards
		result.Changes = append(result.Changes, discoveryResult.Changes...)
		result.Err = errors.Join(result.Err, discoveryErr)
	}
	if err := e.store.Save(result.Cards); err != nil {
		return Result{}, err
	}
	return result, nil
}

// Acknowledge stops animation for a card and persists the choice.
func (e *Engine) Acknowledge(
	cards []domain.Submission,
	key string,
) ([]domain.Submission, error) {
	updated := clone(cards)
	found := false
	for index := range updated {
		if updated[index].Key() == key {
			updated[index].Acknowledged = true
			found = true
			break
		}
	}
	if !found {
		return cards, fmt.Errorf("submission not found")
	}
	if err := e.store.Save(updated); err != nil {
		return cards, err
	}
	return updated, nil
}

// Remove deletes an acknowledged terminal or retained card.
func (e *Engine) Remove(
	cards []domain.Submission,
	key string,
) ([]domain.Submission, error) {
	updated := make([]domain.Submission, 0, len(cards))
	found := false
	for _, card := range cards {
		if card.Key() != key {
			updated = append(updated, card)
			continue
		}
		found = true
		if !card.Acknowledged || !card.Terminal() && !card.Retained {
			return cards, errors.New("acknowledge terminal card before removing it")
		}
	}
	if !found {
		return cards, errors.New("submission not found")
	}
	if err := e.store.Save(updated); err != nil {
		return cards, err
	}
	return updated, nil
}

func (e *Engine) reconcile(
	current []domain.Submission,
	observed []domain.Submission,
) Result {
	byKey := make(map[string]domain.Submission, len(current)+len(observed))
	for _, card := range current {
		byKey[card.Key()] = card
	}
	var changes []Change
	now := e.now()
	for _, next := range observed {
		previous, exists := byKey[next.Key()]
		if !exists {
			next.Acknowledged = true
			if next.LastChangedAt.IsZero() {
				next.LastChangedAt = now
			}
			byKey[next.Key()] = next
			continue
		}
		next.Acknowledged = previous.Acknowledged
		next.LastChangedAt = previous.LastChangedAt
		next.Retained = previous.Retained || next.Terminal()
		if previous.Fingerprint() != next.Fingerprint() {
			next.Acknowledged = false
			next.LastChangedAt = now
			changes = append(changes, Change{Before: previous, After: next})
		}
		byKey[next.Key()] = next
	}
	cards := make([]domain.Submission, 0, len(byKey))
	for _, card := range byKey {
		cards = append(cards, card)
	}
	sort.Slice(cards, func(i, j int) bool {
		if cards[i].SubmittedAt.Equal(cards[j].SubmittedAt) {
			return cards[i].Key() < cards[j].Key()
		}
		return cards[i].SubmittedAt.Before(cards[j].SubmittedAt)
	})
	return Result{Cards: cards, Changes: changes}
}

func clone(cards []domain.Submission) []domain.Submission {
	return append([]domain.Submission(nil), cards...)
}
