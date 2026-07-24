package openurl

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenValidatesHTTPSAndUsesOneArgument(t *testing.T) {
	t.Parallel()

	opener := New("open")
	var path string
	var args []string
	opener.run = func(_ context.Context, command string, values ...string) error {
		path = command
		args = append([]string(nil), values...)
		return nil
	}

	require.NoError(t, opener.Open(context.Background(), "https://example.test/review"))
	require.Equal(t, "open", path)
	require.Equal(t, []string{"https://example.test/review"}, args)
	require.Error(t, opener.Open(context.Background(), "file:///tmp/synthetic"))
}

func TestOpenReturnsSanitizedFailure(t *testing.T) {
	t.Parallel()

	opener := New("open")
	opener.run = func(context.Context, string, ...string) error {
		return errors.New("synthetic process detail")
	}

	err := opener.Open(context.Background(), "https://example.test/review")
	require.ErrorIs(t, err, ErrUnavailable)
	require.NotContains(t, err.Error(), "synthetic process detail")
}
