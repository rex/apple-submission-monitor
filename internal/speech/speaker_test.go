package speech_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rex/apple-submission-monitor/internal/domain"
	"github.com/rex/apple-submission-monitor/internal/speech"
)

type recordingRunner struct {
	path string
	args []string
	err  error
}

func (r *recordingRunner) Run(_ context.Context, path string, args ...string) error {
	r.path = path
	r.args = append([]string(nil), args...)
	return r.err
}

func TestAnnounceUsesDefaultVoice(t *testing.T) {
	t.Parallel()

	runner := &recordingRunner{}
	speaker, err := speech.NewWithRunner("say", templates(), runner)
	require.NoError(t, err)
	before := waitingCard()
	after := before
	after.Health = domain.HealthGreen
	after.ReviewState = "COMPLETE"
	after.AppStoreState = "READY_FOR_DISTRIBUTION"
	after.Outcome = "APPROVED"
	after.InFlight = false

	require.NoError(t, speaker.Announce(context.Background(), before, after))
	require.Equal(t, "say", runner.path)
	require.Equal(t, []string{"Approved Synthetic Alpha"}, runner.args)
	require.NotContains(t, runner.args, "-v")
}

func TestAnnounceFailureIsSanitized(t *testing.T) {
	t.Parallel()

	runner := &recordingRunner{err: errors.New("synthetic failure")}
	speaker, err := speech.NewWithRunner("say", templates(), runner)
	require.NoError(t, err)

	err = speaker.Announce(context.Background(), waitingCard(), waitingCard())
	require.ErrorIs(t, err, speech.ErrUnavailable)
	require.NotContains(t, err.Error(), "synthetic failure")
}

func templates() map[string]string {
	return map[string]string{
		"approved":       "Approved {{.AppName}}",
		"rejected":       "Rejected {{.AppName}}",
		"status_changed": "{{.AppName}} {{.OldStatus}} {{.NewStatus}}",
	}
}

func waitingCard() domain.Submission {
	return domain.Submission{
		AppName:       "Synthetic Alpha",
		Health:        domain.HealthYellow,
		ReviewState:   "WAITING_FOR_REVIEW",
		AppStoreState: "WAITING_FOR_REVIEW",
		InFlight:      true,
	}
}
