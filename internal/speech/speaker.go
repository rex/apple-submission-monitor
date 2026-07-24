// Package speech renders and speaks configurable status announcements.
package speech

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"text/template"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

// ErrUnavailable indicates that the configured speech command cannot run.
var ErrUnavailable = errors.New("speech command unavailable")

// Runner executes one command with an explicit argument vector.
type Runner interface {
	Run(context.Context, string, ...string) error
}

type execRunner struct{}

func (execRunner) Run(ctx context.Context, path string, args ...string) error {
	// #nosec G204 -- path is validated configuration and only rendered text is passed.
	return exec.CommandContext(ctx, path, args...).Run()
}

// Event is the typed data available to announcement templates.
type Event struct {
	AppName   string
	OldStatus string
	NewStatus string
}

// Speaker renders templates and invokes the operating-system speech command.
type Speaker struct {
	path      string
	runner    Runner
	templates map[string]*template.Template
}

// New creates a speaker and validates all announcement templates.
func New(path string, values map[string]string) (*Speaker, error) {
	return NewWithRunner(path, values, execRunner{})
}

// NewWithRunner creates a speaker with an injectable process runner.
func NewWithRunner(
	path string,
	values map[string]string,
	runner Runner,
) (*Speaker, error) {
	templates := make(map[string]*template.Template, len(values))
	for name, value := range values {
		parsed, err := template.New(name).Option("missingkey=error").Parse(value)
		if err != nil {
			return nil, fmt.Errorf("parse announcement template %q: %w", name, err)
		}
		templates[name] = parsed
	}
	return &Speaker{path: path, runner: runner, templates: templates}, nil
}

// Announce speaks a transition using the default system voice.
func (s *Speaker) Announce(
	ctx context.Context,
	before domain.Submission,
	after domain.Submission,
) error {
	name := announcementName(after)
	selected, ok := s.templates[name]
	if !ok {
		return fmt.Errorf("announcement template %q is missing", name)
	}
	event := Event{
		AppName:   after.AppName,
		OldStatus: before.StatusLabel(),
		NewStatus: after.StatusLabel(),
	}
	var rendered bytes.Buffer
	if err := selected.Execute(&rendered, event); err != nil {
		return fmt.Errorf("render announcement: %w", err)
	}
	text := strings.TrimSpace(rendered.String())
	if text == "" || s.path == "" {
		return nil
	}
	if err := s.runner.Run(ctx, s.path, text); err != nil {
		return ErrUnavailable
	}
	return nil
}

func announcementName(card domain.Submission) string {
	if card.Health == domain.HealthGreen && card.Terminal() {
		return "approved"
	}
	if card.Health == domain.HealthRed ||
		strings.Contains(strings.ToUpper(card.AppStoreState), "REJECT") ||
		card.ReviewState == "UNRESOLVED_ISSUES" {
		return "rejected"
	}
	return "status_changed"
}
