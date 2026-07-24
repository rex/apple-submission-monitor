package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

func (m Model) handleKey(message tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch message.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "left", "up", "h", "k":
		m.moveSelection(-1)
	case "right", "down", "l", "j":
		m.moveSelection(1)
	case "enter", " ":
		return m.acknowledgeSelected()
	case "d", "D", "backspace", "delete":
		return m.removeSelected()
	case "o":
		return m.openSelected()
	case "r":
		if !m.loading && !m.actionPending {
			m.loading = true
			return m, refreshCmd(m.ctx, m.engine, m.cards, true)
		}
	case "?":
		m.help = !m.help
	}
	return m, nil
}

func (m Model) handleMouse(message tea.MouseMsg) (tea.Model, tea.Cmd) {
	event := tea.MouseEvent(message)
	if event.Action != tea.MouseActionPress || event.Button != tea.MouseButtonLeft {
		return m, nil
	}
	index := calculateGrid(m.width, m.height, len(m.cards)).hit(event.X, event.Y)
	if index < 0 {
		return m, nil
	}
	m.selected = index
	return m.acknowledgeSelected()
}

func (m Model) acknowledgeSelected() (tea.Model, tea.Cmd) {
	if len(m.cards) == 0 || m.actionPending {
		return m, nil
	}
	card := m.cards[m.selected]
	if card.Acknowledged {
		return m, nil
	}
	m.cards = append([]domain.Submission(nil), m.cards...)
	m.cards[m.selected].Acknowledged = true
	m.actionPending = true
	return m, acknowledgeCmd(m.engine, m.cards, card.Key())
}

func (m Model) removeSelected() (tea.Model, tea.Cmd) {
	if len(m.cards) == 0 || m.actionPending {
		return m, nil
	}
	card := m.cards[m.selected]
	m.actionPending = true
	return m, removeCmd(m.engine, m.cards, card.Key())
}

func (m Model) openSelected() (tea.Model, tea.Cmd) {
	if len(m.cards) == 0 {
		return m, nil
	}
	card := m.cards[m.selected]
	value := card.ReviewURL
	if value == "" {
		value = card.AppStoreURL
	}
	if value == "" {
		m.notice = "No App Store Connect link is available."
		return m, nil
	}
	return m, openCmd(m.ctx, m.opener, value)
}

func (m *Model) moveSelection(delta int) {
	if len(m.cards) == 0 {
		m.selected = 0
		return
	}
	m.selected = (m.selected + delta + len(m.cards)) % len(m.cards)
}

func (m *Model) clampSelection() {
	if len(m.cards) == 0 {
		m.selected = 0
		return
	}
	m.selected = min(m.selected, len(m.cards)-1)
}

func (m Model) hasUnacknowledged() bool {
	for _, card := range m.cards {
		if !card.Acknowledged {
			return true
		}
	}
	return false
}
