package tui

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rex/apple-submission-monitor/internal/domain"
	"github.com/rex/apple-submission-monitor/internal/monitor"
)

func TestMonitorCommandsReturnTypedMessages(t *testing.T) {
	t.Parallel()

	engine := &fakeMonitor{}
	card := testModel(engine).cards[0]
	engine.bootResult = monitor.Result{Cards: []domain.Submission{card}}
	engine.refreshResult = monitor.Result{Cards: []domain.Submission{card}}

	boot := bootstrapCmd(context.Background(), engine)().(bootstrapMsg)
	require.NoError(t, boot.err)
	require.Len(t, boot.result.Cards, 1)

	refresh := refreshCmd(
		context.Background(),
		engine,
		[]domain.Submission{card},
		true,
	)().(refreshMsg)
	require.NoError(t, refresh.err)
	require.True(t, refresh.discovered)

	state := acknowledgeCmd(engine, []domain.Submission{card}, card.Key())().(stateMsg)
	require.NoError(t, state.err)
	require.Equal(t, card.Key(), engine.acknowledged)
}
