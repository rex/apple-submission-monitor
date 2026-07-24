// Package config loads validated local monitor settings.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is the validated runtime configuration.
type Config struct {
	Path              string
	StatePath         string
	ASCPath           string
	SayPath           string
	OpenPath          string
	PollInterval      time.Duration
	DiscoveryInterval time.Duration
	CommandTimeout    time.Duration
	AnimationInterval time.Duration
	DiscoveryWorkers  int
	Announcements     map[string]string
}

type fileConfig struct {
	PollInterval      string            `yaml:"poll_interval"`
	DiscoveryInterval string            `yaml:"discovery_interval"`
	CommandTimeout    string            `yaml:"command_timeout"`
	AnimationInterval string            `yaml:"animation_interval"`
	DiscoveryWorkers  int               `yaml:"discovery_workers"`
	ASCPath           string            `yaml:"asc_path"`
	SayPath           string            `yaml:"say_path"`
	OpenPath          string            `yaml:"open_path"`
	Announcements     map[string]string `yaml:"announcements"`
}

// Load creates a default file when needed, then parses and validates it.
func Load(path string) (Config, error) {
	resolvedPath, statePath, err := resolvePaths(path)
	if err != nil {
		return Config{}, err
	}
	if err := ensureDefaultFile(resolvedPath); err != nil {
		return Config{}, err
	}
	// #nosec G304 -- reading an explicitly selected local configuration is intended.
	content, err := os.ReadFile(resolvedPath)
	if err != nil {
		return Config{}, fmt.Errorf("read configuration: %w", err)
	}
	var raw fileConfig
	if err := yaml.Unmarshal(content, &raw); err != nil {
		return Config{}, fmt.Errorf("parse configuration: %w", err)
	}
	cfg, err := raw.validate()
	if err != nil {
		return Config{}, err
	}
	cfg.Path = resolvedPath
	cfg.StatePath = statePath
	if override := os.Getenv("ASM_POLL_INTERVAL"); override != "" {
		cfg.PollInterval, err = time.ParseDuration(override)
		if err != nil {
			return Config{}, fmt.Errorf("parse ASM_POLL_INTERVAL: %w", err)
		}
	}
	return cfg, nil
}

func (raw fileConfig) validate() (Config, error) {
	poll, err := duration("poll_interval", raw.PollInterval, 5*time.Second)
	if err != nil {
		return Config{}, err
	}
	discovery, err := duration("discovery_interval", raw.DiscoveryInterval, poll)
	if err != nil {
		return Config{}, err
	}
	timeout, err := duration("command_timeout", raw.CommandTimeout, time.Second)
	if err != nil {
		return Config{}, err
	}
	animation, err := duration("animation_interval", raw.AnimationInterval, 100*time.Millisecond)
	if err != nil {
		return Config{}, err
	}
	if raw.DiscoveryWorkers < 1 || raw.DiscoveryWorkers > 16 {
		return Config{}, errors.New("discovery_workers must be between 1 and 16")
	}
	if raw.ASCPath == "" {
		return Config{}, errors.New("asc_path is required")
	}
	requiredTemplates := []string{"approved", "rejected", "status_changed"}
	for _, name := range requiredTemplates {
		if raw.Announcements[name] == "" {
			return Config{}, fmt.Errorf("announcement template %q is required", name)
		}
	}
	return Config{
		ASCPath:           raw.ASCPath,
		SayPath:           raw.SayPath,
		OpenPath:          raw.OpenPath,
		PollInterval:      poll,
		DiscoveryInterval: discovery,
		CommandTimeout:    timeout,
		AnimationInterval: animation,
		DiscoveryWorkers:  raw.DiscoveryWorkers,
		Announcements:     raw.Announcements,
	}, nil
}

func duration(name string, value string, minimum time.Duration) (time.Duration, error) {
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", name, err)
	}
	if parsed < minimum {
		return 0, fmt.Errorf("%s must be at least %s", name, minimum)
	}
	return parsed, nil
}

func resolvePaths(path string) (string, string, error) {
	if path == "" {
		path = os.Getenv("ASM_CONFIG")
	}
	if path != "" {
		return path, filepath.Join(filepath.Dir(path), "state.json"), nil
	}
	base, err := os.UserConfigDir()
	if err != nil {
		return "", "", fmt.Errorf("resolve user config directory: %w", err)
	}
	directory := filepath.Join(base, "apple-submission-monitor")
	return filepath.Join(directory, "config.yaml"), filepath.Join(directory, "state.json"), nil
}

func ensureDefaultFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect configuration: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create configuration directory: %w", err)
	}
	if err := os.WriteFile(path, []byte(DefaultYAML), 0o600); err != nil {
		return fmt.Errorf("write default configuration: %w", err)
	}
	return nil
}

// DefaultYAML is the credential-free configuration created on first launch.
const DefaultYAML = `poll_interval: 30s
discovery_interval: 5m
command_timeout: 25s
animation_interval: 400ms
discovery_workers: 4
asc_path: asc
say_path: /usr/bin/say
open_path: /usr/bin/open
announcements:
  approved: "Great news. {{.AppName}} has been approved."
  rejected: "Attention. {{.AppName}} was rejected. Its status is {{.NewStatus}}."
  status_changed: "{{.AppName}} changed from {{.OldStatus}} to {{.NewStatus}}."
`

// String returns a public-safe description for diagnostics.
func (c Config) String() string {
	return "poll=" + c.PollInterval.String() +
		", discovery=" + c.DiscoveryInterval.String() +
		", workers=" + strconv.Itoa(c.DiscoveryWorkers)
}
