package asc

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

// ErrCommandFailed is returned without leaking stderr or command arguments.
var ErrCommandFailed = errors.New("asc command failed")

// Runner executes an asc argument vector and returns stdout.
type Runner interface {
	Run(context.Context, ...string) ([]byte, error)
}

type commandRunner struct {
	path    string
	timeout time.Duration
}

// NewRunner creates a cancellable direct-process asc runner.
func NewRunner(path string, timeout time.Duration) Runner {
	return &commandRunner{path: path, timeout: timeout}
}

func (r *commandRunner) Run(ctx context.Context, args ...string) ([]byte, error) {
	commandCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	// #nosec G204 -- path is validated configuration and args are fixed by this package.
	output, err := exec.CommandContext(commandCtx, r.path, args...).Output()
	if commandCtx.Err() != nil {
		return nil, fmt.Errorf("asc timed out: %w", commandCtx.Err())
	}
	if err != nil {
		return nil, fmt.Errorf("asc %s: %w", operation(args), ErrCommandFailed)
	}
	return output, nil
}

func operation(args []string) string {
	if len(args) == 0 {
		return "execution"
	}
	if len(args) == 1 {
		return args[0]
	}
	return args[0] + " " + args[1]
}
