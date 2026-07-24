# PROGRESS

<!-- ≤50 lines. Read this first when a fresh agent session starts.
     Points to TASK_STATE.md for the details. -->

- **Project**: apple-submission-monitor (Go CLI)
- **Active branch**: `main`
- **Active feature spec**: `specs/initial-monitor/`
- **Active TASK_STATE**: `TASK_STATE.md` (Phase 1 / Slice 1.1)
- **Last session**: 2026-07-24 (Codex, initial implementation)

## Last three decisions

- 2026-07-24 Use Go with Bubble Tea and Lip Gloss.
- 2026-07-24 Keep the application local-only and Docker-free.
- 2026-07-24 Keep all runtime data outside the public repository.

## Open blockers

- None.

## How to resume (for a fresh agent)

1. Read `AGENTS.md` then `TASK_STATE.md` §0 and current slice.
2. Skim `specs/<slug>/plan.md` for the active phase.
3. Do NOT re-plan if the plan is frozen — follow it.
4. Run `make validate` to see current gate status.

## Do NOT

- Commit live App Store Connect output or local configuration.
- Add direct Apple API clients; all retrieval goes through `asc`.
