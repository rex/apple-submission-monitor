# CONVENTIONS

> Stable conventions. Changes require ADR + PR review.
> Deterministic checks are enforced by pre-commit + CI, not by the agent.

## Runtime
- Go 1.24+.
- Dependency manager: Go modules. `go.sum` is committed.
- Install: `make setup`. Run: `make dev`.

## Layout

```
cmd/apple-submission-monitor/  Thin process entry point.
internal/asc/                   Typed asc command adapter.
internal/config/                Local configuration and defaults.
internal/monitor/               Status transitions and orchestration.
internal/state/                 Private, atomic runtime persistence.
internal/speech/                Default-voice macOS announcements.
internal/tui/                   Bubble Tea model and rendering.
```

Dependencies point inward toward typed domain data. The TUI does not execute
commands or write files directly.

## Concurrency
- Every blocking operation accepts `context.Context`.
- Bubble Tea commands own I/O; `Update` and `View` remain deterministic.
- Cap `asc` fan-out and never overlap polling cycles.

## Error handling
- Wrap errors with `%w` at package boundaries.
- UI-visible failures are concise and retain the last successful snapshot.
- Do not panic after startup.

## Logging
- Do not log raw `asc` output, credentials, app names, IDs, or links.
- The TUI is the operational status surface; diagnostics must be sanitized.

## Testing
- Test framework: Go `testing` plus Testify.
- Coverage: see `VIBE.yaml` `quality_gates.tests.coverage.minimum_percentage`.
- Use table tests and synthetic fixtures only.
- Integration tests execute a fake `asc` binary, never a live account.

## Dependency injection
- Define small interfaces at consumers and wire dependencies in `main`.
- No global mutable state.

## What NOT to put in prompts (the linter's job)

- Formatting rules — pre-commit + `make fix`
- Import ordering — linter
- Type nitpicks — `make typecheck`
- Terraform formatting — `terraform fmt -recursive`

## Stack-specific additions

<!-- If you ran bootstrap-greenfield.sh with a --stack flag, the lang-* skill
     appended its conventions here. -->
