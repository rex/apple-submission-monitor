// Package asc retrieves App Store Connect data exclusively through the asc CLI.
package asc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

// Client decodes stable fields from asc aggregate JSON.
type Client struct {
	runner  Runner
	workers int
	now     func() time.Time
}

// NewClient creates an asc client with bounded discovery concurrency.
func NewClient(runner Runner, workers int) *Client {
	if workers < 1 {
		workers = 1
	}
	return &Client{runner: runner, workers: workers, now: time.Now}
}

func (c *Client) listApps(ctx context.Context) ([]appRecord, error) {
	output, err := c.runner.Run(ctx, "apps", "list", "--paginate", "--output", "json")
	if err != nil {
		return nil, err
	}
	var response appListResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("decode asc apps list: %w", err)
	}
	return response.Data, nil
}

func (c *Client) discoverApp(
	ctx context.Context,
	app appRecord,
) ([]domain.Submission, error) {
	output, err := c.runner.Run(
		ctx,
		"review", "submissions-list",
		"--app", app.ID,
		"--limit", "200",
		"--paginate",
		"--output", "json",
	)
	if err != nil {
		return nil, err
	}
	var response reviewListResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("decode asc review submissions: %w", err)
	}

	var cards []domain.Submission
	var errs []error
	for _, review := range response.Data {
		if !submittedAndActive(review.Attributes) {
			continue
		}
		card, statusErr := c.status(ctx, app.ID, review.Attributes.Platform, review.ID)
		if statusErr != nil {
			errs = append(errs, statusErr)
			continue
		}
		if card.AppName == "" {
			card.AppName = app.Attributes.Name
		}
		if card.BundleID == "" {
			card.BundleID = app.Attributes.BundleID
		}
		cards = append(cards, card)
	}
	return cards, errors.Join(errs...)
}

func (c *Client) status(
	ctx context.Context,
	appID string,
	platform string,
	reviewID string,
) (domain.Submission, error) {
	args := []string{
		"status", "--app", appID,
		"--include", "app,appstore,submission,review,links",
		"--output", "json",
	}
	if platform != "" {
		args = append(args, "--platform", platform)
	}
	output, err := c.runner.Run(ctx, args...)
	if err != nil {
		return domain.Submission{}, err
	}
	var response statusResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return domain.Submission{}, fmt.Errorf("decode asc status: %w", err)
	}
	if response.Review.LatestSubmissionID != "" {
		reviewID = response.Review.LatestSubmissionID
	}
	now := c.now()
	card := domain.Submission{
		ID:            reviewID,
		AppID:         response.App.ID,
		AppName:       response.App.Name,
		BundleID:      response.App.BundleID,
		Platform:      firstNonEmpty(response.AppStore.Platform, response.Review.Platform, platform),
		Version:       response.AppStore.Version,
		VersionID:     response.AppStore.VersionID,
		Health:        domain.NormalizeHealth(response.Summary.Health),
		ReviewState:   response.Review.State,
		AppStoreState: response.AppStore.State,
		NextAction:    response.Summary.NextAction,
		AppStoreURL:   response.Links.AppStoreConnect,
		ReviewURL:     response.Links.Review,
		TestFlightURL: response.Links.TestFlight,
		SubmittedAt:   parseTime(response.Review.SubmittedDate),
		CreatedAt:     parseTime(response.AppStore.CreatedDate),
		LastSeenAt:    now,
		LastChangedAt: now,
		InFlight:      response.Submission.InFlight,
	}
	card.Retained = card.Terminal()
	return card, nil
}

func submittedAndActive(attributes reviewAttributes) bool {
	return attributes.SubmittedDate != "" &&
		attributes.State != "" &&
		attributes.State != "COMPLETE" &&
		attributes.State != "READY_FOR_REVIEW"
}

func parseTime(value string) time.Time {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

type discoveryResult struct {
	cards []domain.Submission
	err   error
}
