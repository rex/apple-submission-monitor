// Command apple-submission-monitor watches active App Store review submissions.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rex/apple-submission-monitor/internal/asc"
	"github.com/rex/apple-submission-monitor/internal/config"
	"github.com/rex/apple-submission-monitor/internal/monitor"
	"github.com/rex/apple-submission-monitor/internal/openurl"
	"github.com/rex/apple-submission-monitor/internal/speech"
	"github.com/rex/apple-submission-monitor/internal/state"
	"github.com/rex/apple-submission-monitor/internal/tui"
)

var version = "dev"

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()
	os.Exit(run(ctx, os.Args[1:], os.Stdout, os.Stderr))
}

func run(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("apple-submission-monitor", flag.ContinueOnError)
	flags.SetOutput(stderr)
	configPath := flags.String("config", "", "path to local YAML configuration")
	noSpeech := flags.Bool("no-speech", false, "disable spoken announcements")
	showVersion := flags.Bool("version", false, "print version and exit")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *showVersion {
		_, _ = fmt.Fprintln(stdout, version)
		return 0
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, "configuration could not be loaded")
		return 1
	}
	if *noSpeech {
		cfg.SayPath = ""
	}
	speaker, err := speech.New(cfg.SayPath, cfg.Announcements)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, "announcement templates are invalid")
		return 1
	}

	runner := asc.NewRunner(cfg.ASCPath, cfg.CommandTimeout)
	client := asc.NewClient(runner, cfg.DiscoveryWorkers)
	store := state.NewStore(cfg.StatePath)
	engine := monitor.NewEngine(client, store)
	opener := openurl.New(cfg.OpenPath)
	model := tui.New(ctx, engine, speaker, opener, cfg)
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		tea.WithContext(ctx),
		tea.WithOutput(stdout),
	)
	if _, err := program.Run(); err != nil && !errors.Is(err, tea.ErrProgramKilled) {
		_, _ = fmt.Fprintln(stderr, "terminal interface stopped unexpectedly")
		return 1
	}
	return 0
}
