package tui

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

func TestSubmissionAgeTextUsesTotalSubmittedTime(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 27, 16, 31, 0, 0, time.UTC)
	card := domain.Submission{SubmittedAt: now.Add(-8*24*time.Hour - 17*time.Hour - 19*time.Minute)}
	require.Equal(t, "8D 17H", submissionAgeText(card, now))

	card.SubmittedAt = now.Add(-17*time.Hour - 19*time.Minute)
	require.Equal(t, "17H 19M", submissionAgeText(card, now))
}

func TestSubmissionAgePaletteThresholds(t *testing.T) {
	t.Parallel()

	require.Equal(t, ageBandGreen, classifySubmissionAge(48*time.Hour))
	require.Equal(t, ageBandYellow, classifySubmissionAge(48*time.Hour+time.Second))
	require.Equal(t, ageBandYellow, classifySubmissionAge(120*time.Hour-time.Second))
	require.Equal(t, ageBandRed, classifySubmissionAge(120*time.Hour))
}
