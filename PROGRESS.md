# PROGRESS

<!-- ≤50 lines. Read this first when a fresh agent session starts.
     Points to TASK_STATE.md for the details. -->

- **Project**: apple-submission-monitor (Go CLI)
- **Active branch**: `main`
- **Active feature spec**: `specs/initial-monitor/`
- **Active TASK_STATE**: `TASK_STATE.md` (Phase 4 complete)
- **Last session**: 2026-07-27 (Codex, dual giant hero and ASC check rails)

## Last three decisions

- 2026-07-27 Use cached filled FIGlet banners with animated multicolor
  gradients and a blood-dripping red-health treatment.
- 2026-07-27 Use total submission age as a same-scale right-hand FIGlet hero,
  with green, gold, and bleeding-red urgency bands.
- 2026-07-27 Enrich wide metadata rails with App ID and truthful ASC checks only
  during full discovery.

## Open blockers

- None.

## How to resume (for a fresh agent)

1. Read `AGENTS.md` then `TASK_STATE.md` §0 and current slice.
2. Skim `specs/<slug>/plan.md` for the active phase.
3. Do NOT re-plan if the plan is frozen — follow it.
4. Run `make validate` to verify the repository.

## Do NOT

- Commit live App Store Connect output or local configuration.
- Add direct Apple API clients; all retrieval goes through `asc`.
