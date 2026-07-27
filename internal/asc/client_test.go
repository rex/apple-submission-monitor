package asc

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type fakeRunner struct {
	mu    sync.Mutex
	calls [][]string
}

func (f *fakeRunner) Run(_ context.Context, args ...string) ([]byte, error) {
	f.mu.Lock()
	f.calls = append(f.calls, append([]string(nil), args...))
	f.mu.Unlock()

	joined := strings.Join(args, " ")
	switch {
	case strings.HasPrefix(joined, "apps list"):
		return []byte(`{"data":[
			{"id":"app-alpha","attributes":{"name":"Synthetic Alpha","bundleId":"test.example.alpha"}},
			{"id":"app-beta","attributes":{"name":"Synthetic Beta","bundleId":"test.example.beta"}}
		]}`), nil
	case strings.Contains(joined, "submissions-list --app app-alpha"):
		return []byte(`{"data":[
			{"id":"review-alpha","attributes":{"platform":"IOS","state":"WAITING_FOR_REVIEW","submittedDate":"2026-01-02T03:04:05Z"}},
			{"id":"review-old","attributes":{"platform":"IOS","state":"COMPLETE","submittedDate":"2025-01-02T03:04:05Z"}}
		]}`), nil
	case strings.Contains(joined, "submissions-list --app app-beta"):
		return []byte(`{"data":[
			{"id":"review-beta","attributes":{"platform":"MAC_OS","state":"READY_FOR_REVIEW","submittedDate":""}}
		]}`), nil
	case strings.HasPrefix(joined, "status --app app-alpha"):
		return []byte(statusFixture), nil
	case strings.HasPrefix(joined, "review status --app app-alpha"):
		return []byte(`{"reviewDetailConfigured":true}`), nil
	case strings.HasPrefix(joined, "versions view --version-id version-alpha"):
		return []byte(`{"buildId":"build-alpha"}`), nil
	case strings.HasPrefix(joined, "builds info --build-id build-alpha"):
		return []byte(`{"data":{"attributes":{"processingState":"VALID"}}}`), nil
	default:
		return nil, errors.New("unexpected synthetic command")
	}
}

func TestDiscoverFiltersNonSubmittedAndCompleteReviews(t *testing.T) {
	t.Parallel()

	runner := &fakeRunner{}
	client := NewClient(runner, 2)
	client.now = func() time.Time {
		return time.Date(2026, 1, 2, 4, 0, 0, 0, time.UTC)
	}

	cards, err := client.Discover(context.Background())
	require.NoError(t, err)
	require.Len(t, cards, 1)
	require.Equal(t, "review-alpha", cards[0].ID)
	require.Equal(t, "Synthetic Alpha", cards[0].AppName)
	require.Equal(t, "1.2.3", cards[0].Version)
	require.Equal(t, "yellow", string(cards[0].Health))
	require.True(t, cards[0].BuildValid())
	require.True(t, cards[0].ReviewKnown)
	require.True(t, cards[0].ReviewDetails)
	require.True(t, cards[0].BlockersKnown)
	require.Zero(t, cards[0].BlockerCount)

	runner.mu.Lock()
	defer runner.mu.Unlock()
	require.True(t, hasPlatformFlag(runner.calls, "IOS"))
}

func TestDiscoverRejectsMalformedAppJSON(t *testing.T) {
	t.Parallel()

	client := NewClient(runnerFunc(func(context.Context, ...string) ([]byte, error) {
		return []byte(`{"data":`), nil
	}), 1)

	_, err := client.Discover(context.Background())
	require.ErrorContains(t, err, "decode asc apps list")
}

func TestRefreshPreservesDiscoveryDiagnostics(t *testing.T) {
	t.Parallel()

	client := NewClient(&fakeRunner{}, 1)
	cards, err := client.Discover(context.Background())
	require.NoError(t, err)

	refreshed, err := client.Refresh(context.Background(), cards)
	require.NoError(t, err)
	require.Len(t, refreshed, 1)
	require.True(t, refreshed[0].BuildValid())
	require.True(t, refreshed[0].ReviewKnown)
	require.True(t, refreshed[0].ReviewDetails)
}

type runnerFunc func(context.Context, ...string) ([]byte, error)

func (function runnerFunc) Run(ctx context.Context, args ...string) ([]byte, error) {
	return function(ctx, args...)
}

func hasPlatformFlag(calls [][]string, platform string) bool {
	for _, call := range calls {
		for index := 0; index+1 < len(call); index++ {
			if call[index] == "--platform" && call[index+1] == platform {
				return true
			}
		}
	}
	return false
}

const statusFixture = `{
  "app":{"id":"app-alpha","name":"Synthetic Alpha","bundleId":"test.example.alpha"},
  "summary":{"health":"yellow","nextAction":"Wait for review.","blockers":[]},
  "appstore":{"createdDate":"2026-01-01T01:02:03Z","platform":"IOS","state":"WAITING_FOR_REVIEW","version":"1.2.3","versionId":"version-alpha"},
  "submission":{"inFlight":true,"blockingIssues":[]},
  "review":{"latestSubmissionId":"review-alpha","platform":"IOS","state":"WAITING_FOR_REVIEW","submittedDate":"2026-01-02T03:04:05Z"},
  "links":{"appStoreConnect":"https://example.test/apps/alpha","review":"https://example.test/reviews/alpha","testFlight":"https://example.test/testflight/alpha"}
}`
