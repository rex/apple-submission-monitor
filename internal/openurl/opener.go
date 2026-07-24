// Package openurl safely opens App Store Connect links.
package openurl

import (
	"context"
	"errors"
	"net/url"
	"os/exec"
)

// ErrUnavailable indicates that a validated URL could not be opened.
var ErrUnavailable = errors.New("URL opener unavailable")

// Opener launches validated HTTPS links with a configured local command.
type Opener struct {
	path string
	run  func(context.Context, string, ...string) error
}

// New creates a URL opener.
func New(path string) *Opener {
	return &Opener{path: path, run: runCommand}
}

// Open validates and opens one URL.
func (o *Opener) Open(ctx context.Context, value string) error {
	parsed, err := url.ParseRequestURI(value)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
		return errors.New("invalid App Store Connect URL")
	}
	if o.path == "" {
		return ErrUnavailable
	}
	if err := o.run(ctx, o.path, value); err != nil {
		return ErrUnavailable
	}
	return nil
}

func runCommand(ctx context.Context, path string, args ...string) error {
	// #nosec G204 -- path is validated configuration and URL is HTTPS-validated.
	return exec.CommandContext(ctx, path, args...).Run()
}
