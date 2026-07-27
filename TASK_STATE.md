# TASK_STATE — Initial monitor

> Source of truth for in-flight work. Humans and agents both write here.
> This file is **committed** to the repo. It survives sessions, machines,
> and context compactions.
>
> Spec: `specs/initial-monitor/spec.md` · Plan: `specs/initial-monitor/plan.md`
> Branch: `main` · Last update: 2026-07-27 by Codex

## 0. TL;DR for a fresh agent session

Build a public-safe Go TUI that discovers only active App Store review
submissions through `asc`, polls their status, announces transitions with the
default macOS `say` voice, and retains terminal submissions until acknowledged
and manually removed. No live app metadata or machine-specific values may be
committed.

## Standing user directives

- Continue until blocked; do not pause between implementation stages.
- Use Go, include unit and integration tests, and keep deployment local-only.
- Treat the repository as public: commit no personal information, secrets,
  machine identifiers, or live App Store Connect data.
- Use `asc` for every App Store Connect data retrieval operation.
- Invoke `say` without a voice option so macOS uses its configured default.

## 1. Phases

| # | Phase | Status | Exit criteria |
|---|---|---|---|
| 1 | Repository scaffold | ✅ done | Public-safe Go policy and gates committed |
| 2 | Data and state engine | ✅ done | Active submissions poll reliably through asc |
| 3 | Interactive TUI | ✅ done | Responsive cards, gradients, mouse, and animation work |
| 4 | Product hardening | ✅ done | Tests, docs, packaging, and completion gates pass |

Statuses: `⏸ pending` · `🟡 in-prog` · `✅ done` · `🔴 blocked`

## 2. Slices (vertical, atomic, independently mergeable)

### Slice 1.1 — Scaffold public-safe Go project

- Status: ✅ done
- Owner: Codex
- Files (planned edits): repository policy, build tooling, Go module
- Files (do NOT edit): live asc authentication or runtime data
- Depends on: none
- Acceptance (EARS notation):
  - [x] Repository contains no personal or machine-specific information.
  - [x] Go build, lint, test, architecture, and public-safety gates are wired.
  - [x] Documentation records the confirmed product behavior.

### Slice 2.1 — Implement the monitor

- Status: ✅ done
- Files (planned edits): `cmd/`, `internal/`, `config/`, tests, documentation

### Slice 4.1 — Harden and document the public release

- Status: ✅ done
- Files: `README.md`, `docs/OPERATIONS.md`, validation scripts, `Makefile`
- Acceptance:
  - [x] Installation, configuration, controls, and long-session operation are documented.
  - [x] Reachable vulnerability and module-integrity checks pass.
  - [x] Public history and tracked-source audits contain no private data.

## 3. Blockers / open questions

- The `asc` global review-submission command still requires an app ID, so
  discovery must list apps and inspect each app's review submissions.

## 4. Recent decisions (append-only, newest first)

- 2026-07-24 — Use Go, Bubble Tea, and Lip Gloss (project owner).
- 2026-07-24 — Poll with asc only; persist runtime state outside the repo.

## 5. Next actions (ordered)

1. None; initial release objective is complete.

## 6. Handoff note (fill when ending a session)

Version 0.6.0 pairs every app-name hero with an equally large total submission
age in a shared FIGlet font. Age art is green through 48 hours, gold until 120
hours, then deteriorated and blood-red. Three bold metadata rails show App ID
plus ASC-derived build, blocker, review-detail, and in-flight checks; diagnostic
refreshes never trigger false transition alerts. Resume from a new objective.
