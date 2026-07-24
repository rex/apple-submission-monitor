# Repository map

## Domains

| Domain | Purpose | Entry point |
|---|---|---|
| Process | Parse configuration and wire the application | `cmd/apple-submission-monitor/main.go` |
| ASC adapter | Discover active reviews and decode status | `internal/asc/client.go` |
| Domain | Persisted submission and health contracts | `internal/domain/submission.go` |
| Monitor | Reconcile snapshots and user acknowledgment | `internal/monitor/engine.go` |
| Configuration | Create and validate editable YAML | `internal/config/config.go` |
| State | Atomically persist cards outside the repo | `internal/state/store.go` |
| Speech | Render templates and invoke default-voice speech | `internal/speech/speaker.go` |
| TUI | Adaptive layout, input, and animation | `internal/tui/` |

## External boundaries

| System | Access | Failure behavior |
|---|---|---|
| `asc` | Direct child process with fixed arguments | Retain last snapshot and show sanitized warning |
| `say` | Announcement text as the sole argument | Show non-fatal warning and continue |
| `open` | Selected App Store Connect HTTPS link | Show non-fatal warning |
| Local state | Versioned private JSON and YAML | Fail safely without leaking paths |

## Invariants

- `asc` is the only App Store Connect data source.
- Initial discovery includes only submitted, non-complete reviews.
- A meaningful transition becomes unacknowledged and persists immediately.
- Terminal cards survive upstream disappearance until acknowledged and removed.
- Repository tests and examples contain synthetic data only.

## Read order

1. `README.md`
2. `AGENTS.md` and `CONVENTIONS.md`
3. `specs/initial-monitor/spec.md`
4. `internal/domain/submission.go`
5. `internal/monitor/engine.go`
6. `internal/tui/`
