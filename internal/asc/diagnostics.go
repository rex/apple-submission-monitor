package asc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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
