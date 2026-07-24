package asc

import (
	"context"
	"errors"
	"sort"
	"sync"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

// Discover returns only apps with submitted, non-complete review submissions.
func (c *Client) Discover(ctx context.Context) ([]domain.Submission, error) {
	apps, err := c.listApps(ctx)
	if err != nil {
		return nil, err
	}

	jobs := make(chan appRecord)
	results := make(chan discoveryResult, len(apps))
	var group sync.WaitGroup
	for range min(c.workers, len(apps)) {
		group.Add(1)
		go func() {
			defer group.Done()
			for app := range jobs {
				cards, discoverErr := c.discoverApp(ctx, app)
				results <- discoveryResult{cards: cards, err: discoverErr}
			}
		}()
	}
	go sendApps(ctx, apps, jobs)
	go func() {
		group.Wait()
		close(results)
	}()

	var cards []domain.Submission
	var errs []error
	for result := range results {
		cards = append(cards, result.cards...)
		if result.err != nil {
			errs = append(errs, result.err)
		}
	}
	sortCards(cards)
	return cards, errors.Join(errs...)
}

// Refresh retrieves current status for each non-terminal persisted card.
func (c *Client) Refresh(
	ctx context.Context,
	cards []domain.Submission,
) ([]domain.Submission, error) {
	jobs := make(chan domain.Submission)
	results := make(chan refreshResult, len(cards))
	var group sync.WaitGroup
	for range min(c.workers, len(cards)) {
		group.Add(1)
		go func() {
			defer group.Done()
			for card := range jobs {
				current, err := c.status(ctx, card.AppID, card.Platform, card.ID)
				results <- refreshResult{card: current, err: err}
			}
		}()
	}
	go sendCards(ctx, cards, jobs)
	go func() {
		group.Wait()
		close(results)
	}()

	var refreshed []domain.Submission
	var errs []error
	for result := range results {
		if result.err != nil {
			errs = append(errs, result.err)
			continue
		}
		refreshed = append(refreshed, result.card)
	}
	sortCards(refreshed)
	return refreshed, errors.Join(errs...)
}

func sendApps(ctx context.Context, apps []appRecord, jobs chan<- appRecord) {
	defer close(jobs)
	for _, app := range apps {
		select {
		case jobs <- app:
		case <-ctx.Done():
			return
		}
	}
}

func sendCards(
	ctx context.Context,
	cards []domain.Submission,
	jobs chan<- domain.Submission,
) {
	defer close(jobs)
	for _, card := range cards {
		if card.Retained && card.Terminal() {
			continue
		}
		select {
		case jobs <- card:
		case <-ctx.Done():
			return
		}
	}
}

func sortCards(cards []domain.Submission) {
	sort.Slice(cards, func(i, j int) bool {
		if cards[i].SubmittedAt.Equal(cards[j].SubmittedAt) {
			return cards[i].Key() < cards[j].Key()
		}
		return cards[i].SubmittedAt.Before(cards[j].SubmittedAt)
	})
}

type refreshResult struct {
	card domain.Submission
	err  error
}
