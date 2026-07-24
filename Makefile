# ╔══════════════════════════════════════════════════════════════════════╗
# ║          Apple Submission Monitor — Makefile                         ║
# ║          Go CLI + Bubble Tea                                         ║
# ╚══════════════════════════════════════════════════════════════════════╝
#
# Usage: make <target>
# Run `make help` for a full list of available targets.

.PHONY: help install install-cli setup check-setup validate update info \
        dev build start lint typecheck check-architecture fix test coverage \
        check-public-safety check-docs security \
        bump-patch bump-minor bump-major check-version-bumped version \
        clean clean-all \
        check-if-the-agent-can-consider-this-task-completed

# ─── Configuration ────────────────────────────────────────────────────
# Prefer Homebrew zsh on macOS; fall back to whichever zsh is on PATH.
# We require zsh (recipes use zsh-isms). If neither path resolves, Make
# will surface a clear error on first recipe execution.
SHELL       := $(or $(wildcard /opt/homebrew/bin/zsh),$(shell command -v zsh))
APP_NAME    ?= apple-submission-monitor
GOLANGCI_VERSION ?= v2.12.2
GOVULNCHECK_VERSION ?= v1.6.0
COVERAGE_MINIMUM ?= 70
INSTALL_DIR ?= $(shell go env GOPATH)/bin

# Colors for output
CYAN   := $(shell printf '\033[36m')
GREEN  := $(shell printf '\033[32m')
YELLOW := $(shell printf '\033[33m')
RED    := $(shell printf '\033[31m')
RESET  := $(shell printf '\033[0m')
BOLD   := $(shell printf '\033[1m')

# ─── Help ─────────────────────────────────────────────────────────────

## help: Display this help message with all available targets
help:
	@echo ""
	@echo "$(BOLD)$(CYAN)$(APP_NAME)$(RESET)"
	@echo "$(CYAN)════════════════════════════════════════════════════$(RESET)"
	@echo ""
	@echo "$(BOLD)Setup & Installation$(RESET)"
	@echo "  $(GREEN)make install$(RESET)              Install dependencies"
	@echo "  $(GREEN)make install-cli$(RESET)          Install the versioned binary"
	@echo "  $(GREEN)make setup$(RESET)                Full setup: install + typecheck"
	@echo ""
	@echo "$(BOLD)Development$(RESET)"
	@echo "  $(GREEN)make dev$(RESET)                  Run the monitor from source"
	@echo "  $(GREEN)make build$(RESET)                Build the native binary"
	@echo "  $(GREEN)make start$(RESET)                Run the built binary"
	@echo "  $(GREEN)make lint$(RESET)                 Run linter"
	@echo "  $(GREEN)make typecheck$(RESET)            Run type checker"
	@echo "  $(GREEN)make check-architecture$(RESET)   Run repo-native architecture checks"
	@echo "  $(GREEN)make fix$(RESET)                  Auto-fix lint issues"
	@echo "  $(GREEN)make test$(RESET)                 Run race-enabled tests"
	@echo "  $(GREEN)make coverage$(RESET)             Enforce the coverage threshold"
	@echo "  $(GREEN)make check-public-safety$(RESET)  Scan tracked files for private data"
	@echo "  $(GREEN)make check-docs$(RESET)           Validate local Markdown links"
	@echo "  $(GREEN)make security$(RESET)             Verify modules and reachable vulnerabilities"
	@echo "  $(GREEN)make validate$(RESET)             Run every required quality gate"
	@echo ""
	@echo "$(BOLD)Versioning (required before every commit)$(RESET)"
	@echo "  $(GREEN)make version$(RESET)              Print current VERSION"
	@echo "  $(GREEN)make bump-patch$(RESET)           Bump patch (x.y.Z+1) — bug fixes / doc / refactor"
	@echo "  $(GREEN)make bump-minor$(RESET)           Bump minor (x.Y+1.0) — additive feature, backward-compat"
	@echo "  $(GREEN)make bump-major$(RESET)           Bump major (X+1.0.0) — breaking change"
	@echo ""
	@echo "$(BOLD)Maintenance$(RESET)"
	@echo "  $(GREEN)make clean$(RESET)                Remove build cache"
	@echo "  $(GREEN)make clean-all$(RESET)            Remove build cache + deps (destructive!)"
	@echo "  $(GREEN)make update$(RESET)               Update dependencies"
	@echo "  $(GREEN)make info$(RESET)                 Show project info"
	@echo ""
	@echo "$(BOLD)Completion$(RESET)"
	@echo "  $(GREEN)make check-if-the-agent-can-consider-this-task-completed$(RESET)"
	@echo "    Final verification gate (required before declaring a task complete)"
	@echo ""

# ─── Setup & Installation ────────────────────────────────────────────

## install: Download Go module dependencies
install:
	@echo "$(CYAN)Installing dependencies...$(RESET)"
	@go mod download
	@echo "$(GREEN)Done.$(RESET)"

## install-cli: Build and install the versioned binary into Go's bin directory
install-cli: build
	@mkdir -p "$(INSTALL_DIR)"
	@install -m 0755 "bin/$(APP_NAME)" "$(INSTALL_DIR)/$(APP_NAME)"
	@echo "$(GREEN)Installed $(INSTALL_DIR)/$(APP_NAME)$(RESET)"

## setup: Full project setup
setup: install typecheck
	@echo ""
	@echo "$(GREEN)$(BOLD)Setup complete!$(RESET)"
	@echo "  Run $(CYAN)make dev$(RESET) to start developing."

## check-setup: Execute the public setup entrypoint as a regression test
check-setup:
	@$(MAKE) --no-print-directory setup >/dev/null
	@echo "$(GREEN)Setup entrypoint passed.$(RESET)"

## validate: Run the repo's aggregate validation flow
validate: lint typecheck test coverage check-architecture check-public-safety check-docs security check-setup check-version-bumped
	@echo "$(GREEN)Validation complete.$(RESET)"

# ─── Versioning (non-negotiable: every commit gets a bump) ───────────

## bump-patch: Increment patch (x.y.Z+1) — bug fixes, docs, non-behavior changes
bump-patch:
	@scripts/bump_version.py patch

## bump-minor: Increment minor (x.Y+1.0) — additive features, backward-compatible
bump-minor:
	@scripts/bump_version.py minor

## bump-major: Increment major (X+1.0.0) — breaking change, removal, incompat behavior
bump-major:
	@scripts/bump_version.py major

## check-version-bumped: Fail if VERSION == HEAD's VERSION or CHANGELOG lacks matching entry
check-version-bumped:
	@if [ -x scripts/check_version_bumped.py ]; then \
		scripts/check_version_bumped.py; \
	else \
		echo "$(YELLOW)scripts/check_version_bumped.py missing — skipping gate (bootstrap mode?).$(RESET)"; \
	fi

## version: Print current VERSION
version:
	@cat VERSION 2>/dev/null || echo "0.1.0 (VERSION file missing)"

# ─── Development ──────────────────────────────────────────────────────

## dev: Run the monitor from source
dev:
	@go run ./cmd/apple-submission-monitor

## build: Build the native binary
build:
	@echo "$(CYAN)Building...$(RESET)"
	@mkdir -p bin
	@CGO_ENABLED=0 go build -trimpath \
		-ldflags="-s -w -X main.version=$$(cat VERSION)" \
		-o bin/$(APP_NAME) ./cmd/apple-submission-monitor

## start: Run the built binary
start:
	@./bin/$(APP_NAME)

## lint: Run formatting and focused static analysis
lint:
	@echo "$(CYAN)Running linter...$(RESET)"
	@test -z "$$(gofmt -l $$(find cmd internal -name '*.go' 2>/dev/null))"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION) run; \
	fi

## typecheck: Compile and run Go's static checks
typecheck:
	@go vet ./...

## check-architecture: Run repo-native architecture checks
check-architecture:
	@echo "$(CYAN)Checking architecture (VIBE.yaml limits)...$(RESET)"
	@python3 scripts/check-architecture.py

## fix: Auto-fix Go formatting
fix:
	@gofmt -w $$(find cmd internal -name '*.go' 2>/dev/null)

## test: Run all tests with the race detector
test:
	@echo "$(CYAN)Running tests...$(RESET)"
	@go test -race -count=1 ./...

## coverage: Enforce the configured line coverage threshold
coverage:
	@go test -coverprofile=coverage.out ./...
	@scripts/check-coverage.sh coverage.out $(COVERAGE_MINIMUM)

## check-public-safety: Scan tracked content for private or machine-specific data
check-public-safety:
	@scripts/check-public-safety.sh

## check-docs: Validate local Markdown links
check-docs:
	@python3 scripts/check-doc-links.py

## security: Verify modules and scan reachable code for known vulnerabilities
security:
	@go mod verify
	@go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) ./...

# ─── Maintenance ──────────────────────────────────────────────────────

## clean: Remove build cache
clean:
	@echo "$(CYAN)Cleaning build cache...$(RESET)"
	@go clean -cache -testcache
	@find bin -type f -delete 2>/dev/null || true
	@echo "$(GREEN)Clean.$(RESET)"

## clean-all: Remove build cache + dependencies (destructive — requires confirmation)
clean-all:
	@echo "$(YELLOW)WARNING: This will remove the Go module and build caches.$(RESET)"
	@read -p "Are you sure? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		go clean -cache -testcache -modcache; \
		find bin -type f -delete 2>/dev/null || true; \
		echo "$(GREEN)Deep clean complete.$(RESET)"; \
	else \
		echo "$(YELLOW)Cancelled.$(RESET)"; \
	fi

## update: Update Go dependencies
update:
	@echo "$(CYAN)Updating dependencies...$(RESET)"
	@go get -u ./...
	@go mod tidy

## info: Show project state
info:
	@echo "$(BOLD)$(CYAN)Project Info$(RESET)"
	@echo "──────────────────────────────"
	@echo "  Project: $(APP_NAME)"
	@echo "  Branch:  $$(git branch --show-current 2>/dev/null || echo 'N/A')"
	@echo "  Commit:  $$(git rev-parse --short HEAD 2>/dev/null || echo 'N/A')"
	@echo "  Tree:    $$(git status --porcelain | wc -l | tr -d ' ') uncommitted changes"

# ─── Completion Gate ──────────────────────────────────────────────────

## check-if-the-agent-can-consider-this-task-completed: Final verification gate
check-if-the-agent-can-consider-this-task-completed: validate
	@echo ""
	@echo "$(BOLD)$(GREEN)✓ All gates passed. Task may be declared complete.$(RESET)"
	@echo ""
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "$(YELLOW)NOTE: working tree is dirty:$(RESET)"; \
		git status --short; \
		echo ""; \
		echo "$(YELLOW)The gates passed, but VIBE.yaml clean_worktree_required_on_completion$(RESET)"; \
		echo "$(YELLOW)may still apply. Commit or stash before declaring done.$(RESET)"; \
	fi

.DEFAULT_GOAL := help
