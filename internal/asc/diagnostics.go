package asc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

func (c *Client) enrichDiagnostics(
	ctx context.Context,
	card *domain.Submission,
) error {
	var errs []error
	if err := c.reviewDiagnostics(ctx, card); err != nil {
		errs = append(errs, err)
	}
	if card.Terminal() {
		if err := c.reviewOutcome(ctx, card); err != nil {
			errs = append(errs, err)
		}
	}
	if err := c.buildDiagnostics(ctx, card); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (c *Client) reviewDiagnostics(
	ctx context.Context,
	card *domain.Submission,
) error {
	args := []string{
		"review", "status",
		"--app", card.AppID,
		"--output", "json",
	}
	if card.VersionID != "" {
		args = append(args, "--version-id", card.VersionID)
	}
	if card.Platform != "" {
		args = append(args, "--platform", card.Platform)
	}
	output, err := c.runner.Run(ctx, args...)
	if err != nil {
		return err
	}
	var response reviewStatusResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return fmt.Errorf("decode asc review status: %w", err)
	}
	card.ReviewKnown = true
	card.ReviewDetails = response.ReviewDetailConfigured
	card.VersionID = firstNonEmpty(response.Version.ID, card.VersionID)
	card.Version = firstNonEmpty(response.Version.Version, card.Version)
	card.Platform = firstNonEmpty(response.Version.Platform, response.LatestSubmission.Platform, card.Platform)
	card.AppStoreState = firstNonEmpty(response.Version.State, card.AppStoreState)
	card.ReviewState = firstNonEmpty(response.ReviewState, response.LatestSubmission.State, card.ReviewState)
	card.NextAction = firstNonEmpty(response.NextAction, card.NextAction)
	card.ID = firstNonEmpty(response.LatestSubmission.ID, card.ID)
	if submitted := parseTime(response.LatestSubmission.SubmittedDate); !submitted.IsZero() {
		card.SubmittedAt = submitted
	}
	if created := parseTime(response.Version.CreatedDate); !created.IsZero() {
		card.CreatedAt = created
	}
	return nil
}

func (c *Client) reviewOutcome(
	ctx context.Context,
	card *domain.Submission,
) error {
	if card.Version == "" {
		return nil
	}
	args := []string{
		"review", "history",
		"--app", card.AppID,
		"--version", card.Version,
		"--paginate",
		"--output", "json",
	}
	if card.Platform != "" {
		args = append(args, "--platform", card.Platform)
	}
	output, err := c.runner.Run(ctx, args...)
	if err != nil {
		return err
	}
	var response reviewHistoryResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return fmt.Errorf("decode asc review history: %w", err)
	}
	for _, entry := range response {
		if card.ID != "" && entry.SubmissionID != "" && entry.SubmissionID != card.ID {
			continue
		}
		if outcome := strings.TrimSpace(entry.Outcome); outcome != "" {
			card.Outcome = strings.ToUpper(outcome)
			return nil
		}
	}
	return nil
}

func (c *Client) buildDiagnostics(
	ctx context.Context,
	card *domain.Submission,
) error {
	if card.VersionID == "" {
		return nil
	}
	output, err := c.runner.Run(
		ctx,
		"versions", "view",
		"--version-id", card.VersionID,
		"--include-build",
		"--output", "json",
	)
	if err != nil {
		return err
	}
	var version versionViewResponse
	if err := json.Unmarshal(output, &version); err != nil {
		return fmt.Errorf("decode asc version build: %w", err)
	}
	if version.BuildID == "" {
		return nil
	}

	output, err = c.runner.Run(
		ctx,
		"builds", "info",
		"--build-id", version.BuildID,
		"--output", "json",
	)
	if err != nil {
		return err
	}
	var build buildInfoResponse
	if err := json.Unmarshal(output, &build); err != nil {
		return fmt.Errorf("decode asc build info: %w", err)
	}
	card.BuildKnown = build.Data.Attributes.ProcessingState != ""
	card.BuildState = build.Data.Attributes.ProcessingState
	return nil
}
