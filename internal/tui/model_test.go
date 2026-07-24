package tui

import (
	"context"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/require"

	"github.com/rex/apple-submission-monitor/internal/domain"
	"github.com/rex/apple-submission-monitor/internal/monitor"
)

func TestMouseAcknowledgesChangedCard(t *testing.T) {
	t.Parallel()

	engine := &fakeMonitor{}
	model := testModel(engine)
	model.cards[0].Acknowledged = false
	area := calculateGrid(model.width, model.height, len(model.cards)).rects[0]
	message := tea.MouseMsg(tea.MouseEvent{
		X: area.x + 1, Y: area.y + 1,
		Action: tea.MouseActionPress, Button: tea.MouseButtonLeft,
	})

	updatedModel, command := model.Update(message)
	updated := updatedModel.(Model)
	require.True(t, updated.cards[0].Acknowledged)
	require.True(t, updated.actionPending)
	require.NotNil(t, command)

	result := command()
	updatedModel, _ = updated.Update(result)
	updated = updatedModel.(Model)
	require.False(t, updated.actionPending)
	require.Equal(t, "review-alpha", engine.acknowledged)
}

func TestKeyboardActionsNavigateOpenAndRemove(t *testing.T) {
	t.Parallel()

	engine := &fakeMonitor{}
	opener := &fakeOpener{}
	model := testModel(engine)
	model.opener = opener
	model.cards = append(model.cards, terminalCard())

	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	updated := updatedModel.(Model)
	require.Equal(t, 1, updated.selected)

	updatedModel, command := updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	require.NotNil(t, command)
	_, _ = updatedModel.(Model).Update(command())
	require.Equal(t, "https://example.test/review", opener.value)

	updatedModel, command = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'D'}})
	require.NotNil(t, command)
	_, _ = updatedModel.(Model).Update(command())
	require.Equal(t, "review-terminal", engine.removed)
}

func TestApplyResultAnimatesAndAnnouncesChanges(t *testing.T) {
	t.Parallel()

	speaker := &fakeSpeaker{}
	model := testModel(&fakeMonitor{})
	model.speaker = speaker
	before := model.cards[0]
	after := before
	after.Acknowledged = false
	after.Health = domain.HealthGreen

	updatedModel, command := model.applyResult(monitor.Result{
		Cards:   []domain.Submission{after},
		Changes: []monitor.Change{{Before: before, After: after}},
	}, false)
	updated := updatedModel.(Model)
	require.True(t, updated.animating)
	require.NotNil(t, command)

	batch := command().(tea.BatchMsg)
	require.Len(t, batch, 2)
	for _, child := range batch {
		message := child()
		if _, ok := message.(speechMsg); ok {
			_, _ = updated.Update(message)
		}
	}
	require.Equal(t, 1, speaker.count)
}

func TestViewFillsTerminalAndShowsFallbacks(t *testing.T) {
	t.Parallel()

	model := testModel(&fakeMonitor{})
	view := model.View()
	require.Equal(t, model.height, lipgloss.Height(view))
	require.Contains(t, view, "WAITING FOR REVIEW")
	require.Contains(t, view, "App Store Connect")

	model.width = 20
	require.Contains(t, model.View(), "Terminal too small")
	require.Equal(t, "link", hyperlink("javascript:alert(1)", "link"))
}

func TestLifecycleMessagesRemainNonFatal(t *testing.T) {
	t.Parallel()

	model := testModel(&fakeMonitor{})
	model.loaded = false
	model.loading = true
	updatedModel, _ := model.Update(bootstrapMsg{
		result: monitor.Result{Cards: model.cards},
	})
	updated := updatedModel.(Model)
	require.True(t, updated.loaded)
	require.False(t, updated.loading)

	updatedModel, _ = updated.Update(refreshMsg{err: context.Canceled})
	updated = updatedModel.(Model)
	require.Contains(t, updated.notice, "last successful")

	updatedModel, _ = updated.Update(speechMsg{err: context.Canceled})
	updated = updatedModel.(Model)
	require.Contains(t, updated.notice, "Speech")

	updatedModel, _ = updated.Update(openMsg{err: context.Canceled})
	updated = updatedModel.(Model)
	require.Contains(t, updated.notice, "open")

	updated.cards[0].Acknowledged = true
	updated.animating = true
	updatedModel, command := updated.Update(animationTickMsg(time.Now()))
	updated = updatedModel.(Model)
	require.False(t, updated.animating)
	require.Nil(t, command)
}

func TestManualRefreshAndHelp(t *testing.T) {
	t.Parallel()

	model := testModel(&fakeMonitor{})
	updatedModel, command := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	updated := updatedModel.(Model)
	require.True(t, updated.loading)
	require.NotNil(t, command)

	updated.loading = false
	updatedModel, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	require.True(t, updatedModel.(Model).help)

	_, command = updated.Update(pollTickMsg(time.Now()))
	require.NotNil(t, command)
}
