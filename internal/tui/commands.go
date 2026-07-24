package tui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rex/apple-submission-monitor/internal/domain"
	"github.com/rex/apple-submission-monitor/internal/monitor"
)

type bootstrapMsg struct {
	result monitor.Result
	err    error
}

type refreshMsg struct {
	result     monitor.Result
	err        error
	discovered bool
}

type stateMsg struct {
	cards []domain.Submission
	err   error
}

type speechMsg struct {
	err error
}

type openMsg struct {
	err error
}

type pollTickMsg time.Time
type animationTickMsg time.Time

func bootstrapCmd(ctx context.Context, engine Monitor) tea.Cmd {
	return func() tea.Msg {
		result, err := engine.Bootstrap(ctx)
		return bootstrapMsg{result: result, err: err}
	}
}

func refreshCmd(
	ctx context.Context,
	engine Monitor,
	cards []domain.Submission,
	discover bool,
) tea.Cmd {
	snapshot := append([]domain.Submission(nil), cards...)
	return func() tea.Msg {
		result, err := engine.Refresh(ctx, snapshot, discover)
		return refreshMsg{result: result, err: err, discovered: discover}
	}
}

func acknowledgeCmd(
	engine Monitor,
	cards []domain.Submission,
	key string,
) tea.Cmd {
	snapshot := append([]domain.Submission(nil), cards...)
	return func() tea.Msg {
		updated, err := engine.Acknowledge(snapshot, key)
		return stateMsg{cards: updated, err: err}
	}
}

func removeCmd(engine Monitor, cards []domain.Submission, key string) tea.Cmd {
	snapshot := append([]domain.Submission(nil), cards...)
	return func() tea.Msg {
		updated, err := engine.Remove(snapshot, key)
		return stateMsg{cards: updated, err: err}
	}
}

func announceCmd(
	ctx context.Context,
	speaker Speaker,
	changes []monitor.Change,
) tea.Cmd {
	return func() tea.Msg {
		for _, change := range changes {
			if err := speaker.Announce(ctx, change.Before, change.After); err != nil {
				return speechMsg{err: err}
			}
		}
		return speechMsg{}
	}
}

func openCmd(ctx context.Context, opener Opener, value string) tea.Cmd {
	return func() tea.Msg {
		return openMsg{err: opener.Open(ctx, value)}
	}
}

func pollTickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(value time.Time) tea.Msg {
		return pollTickMsg(value)
	})
}

func animationTickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(value time.Time) tea.Msg {
		return animationTickMsg(value)
	})
}
