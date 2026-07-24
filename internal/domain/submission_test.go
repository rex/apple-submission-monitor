package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

func TestSubmissionStatusAndTerminal(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		card     domain.Submission
		label    string
		terminal bool
	}{
		"waiting": {
			card: domain.Submission{
				Health:        domain.HealthYellow,
				AppStoreState: "WAITING_FOR_REVIEW",
				InFlight:      true,
			},
			label: "Waiting For Review",
		},
		"complete": {
			card: domain.Submission{
				Health:      domain.HealthGreen,
				ReviewState: "COMPLETE",
			},
			label:    "Complete",
			terminal: true,
		},
		"rejected": {
			card: domain.Submission{
				Health:        domain.HealthRed,
				AppStoreState: "REJECTED",
			},
			label:    "Rejected",
			terminal: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.label, test.card.StatusLabel())
			require.Equal(t, test.terminal, test.card.Terminal())
		})
	}
}

func TestNormalizeHealth(t *testing.T) {
	t.Parallel()

	require.Equal(t, domain.HealthGreen, domain.NormalizeHealth(" GREEN "))
	require.Equal(t, domain.HealthGray, domain.NormalizeHealth("chartreuse"))
}

func TestFingerprintIgnoresMetadataRefresh(t *testing.T) {
	t.Parallel()

	before := domain.Submission{
		ID:            "review-alpha",
		Health:        domain.HealthYellow,
		AppStoreState: "WAITING_FOR_REVIEW",
		InFlight:      true,
	}
	after := before
	after.AppName = "Synthetic Alpha"
	after.Version = "2.0"

	require.Equal(t, before.Fingerprint(), after.Fingerprint())
}
