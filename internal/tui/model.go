// Package tui implements the interactive Apple submission dashboard.
package tui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rex/apple-submission-monitor/internal/config"
	"github.com/rex/apple-submission-monitor/internal/domain"
	"github.com/rex/apple-submission-monitor/internal/monitor"
)

// Monitor is the stateful application boundary consumed by the TUI.
type Monitor interface {
	Bootstrap(context.Context) (monitor.Result, error)
	Refresh(context.Context, []domain.Submission, bool) (monitor.Result, error)
	Acknowledge([]domain.Submission, string) ([]domain.Submission, error)
	Remove([]domain.Submission, string) ([]domain.Submission, error)
}

// Speaker announces meaningful status changes.
type Speaker interface {
	Announce(context.Context, domain.Submission, domain.Submission) error
}

// Opener opens a validated App Store Connect link.
type Opener interface {
	Open(context.Context, string) error
}

// Model is the Bubble Tea application state.
type Model struct {
	ctx               context.Context
	engine            Monitor
	speaker           Speaker
	opener            Opener
	cards             []domain.Submission
	width             int
	height            int
	selected          int
	loading           bool
	loaded            bool
	actionPending     bool
	help              bool
	pulse             bool
	animating         bool
	notice            string
	configPath        string
	pollInterval      time.Duration
	discoveryInterval time.Duration
	animationInterval time.Duration
	lastUpdated       time.Time
	lastDiscovery     time.Time
}

// New creates a fully wired TUI model.
func New(
	ctx context.Context,
	engine Monitor,
	speaker Speaker,
	opener Opener,
	cfg config.Config,
) Model {
	return Model{
		ctx:               ctx,
		engine:            engine,
		speaker:           speaker,
		opener:            opener,
		loading:           true,
		configPath:        cfg.Path,
		pollInterval:      cfg.PollInterval,
		discoveryInterval: cfg.DiscoveryInterval,
		animationInterval: cfg.AnimationInterval,
	}
}

// Init starts bootstrap and the shared poll timer.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		bootstrapCmd(m.ctx, m.engine),
		pollTickCmd(m.pollInterval),
	)
}

// Update handles terminal, timer, and I/O-result messages.
func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = message.Width, message.Height
	case tea.KeyMsg:
		return m.handleKey(message)
	case tea.MouseMsg:
		return m.handleMouse(message)
	case bootstrapMsg:
		m.loading = false
		m.loaded = true
		if message.err != nil {
			m.notice = "Could not load local monitor state."
			return m, nil
		}
		return m.applyResult(message.result, true)
	case refreshMsg:
		m.loading = false
		if message.err != nil {
			m.notice = "Refresh failed; showing last successful data."
			return m, nil
		}
		if message.discovered {
			m.lastDiscovery = time.Now()
		}
		return m.applyResult(message.result, false)
	case stateMsg:
		m.actionPending = false
		if message.err != nil {
			m.notice = "Could not save that action."
			return m, nil
		}
		m.cards = message.cards
		m.clampSelection()
	case speechMsg:
		if message.err != nil {
			m.notice = "Speech announcement unavailable; monitoring continues."
		}
	case openMsg:
		if message.err != nil {
			m.notice = "Could not open the App Store Connect link."
		}
	case pollTickMsg:
		commands := []tea.Cmd{pollTickCmd(m.pollInterval)}
		if !m.loading && !m.actionPending {
			discover := m.lastDiscovery.IsZero() ||
				time.Since(m.lastDiscovery) >= m.discoveryInterval
			m.loading = true
			commands = append(commands, refreshCmd(m.ctx, m.engine, m.cards, discover))
		}
		return m, tea.Batch(commands...)
	case animationTickMsg:
		if !m.shouldAnimate() {
			m.animating = false
			m.pulse = false
			return m, nil
		}
		m.pulse = !m.pulse
		return m, animationTickCmd(m.nextAnimationInterval())
	}
	return m, nil
}

func (m Model) applyResult(result monitor.Result, bootstrap bool) (tea.Model, tea.Cmd) {
	m.cards = result.Cards
	m.lastUpdated = time.Now()
	if bootstrap {
		m.lastDiscovery = m.lastUpdated
	}
	m.clampSelection()
	m.notice = ""
	if result.Err != nil {
		m.notice = "Some asc requests failed; showing available data."
	}
	var commands []tea.Cmd
	if len(result.Changes) > 0 {
		commands = append(commands, announceCmd(m.ctx, m.speaker, result.Changes))
	}
	if m.shouldAnimate() && !m.animating {
		m.animating = true
		commands = append(commands, animationTickCmd(m.nextAnimationInterval()))
	}
	return m, tea.Batch(commands...)
}
