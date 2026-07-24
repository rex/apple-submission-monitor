# TASK_STATE — Initial monitor

> Source of truth for in-flight work. Humans and agents both write here.
> This file is **committed** to the repo. It survives sessions, machines,
> and context compactions.
>
> Spec: `specs/initial-monitor/spec.md` · Plan: `specs/initial-monitor/plan.md`
> Branch: `main` · Last update: 2026-07-24 by Codex

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
| 3 | Interactive TUI | 🟡 in-prog | Responsive cards, gradients, mouse, and animation work |
| 4 | Product hardening | ⏸ pending | Tests, docs, CI, and completion gates pass |

Statuses: `⏸ pending` · `🟡 in-prog` · `✅ done` · `🔴 blocked`

## 2. Slices (vertical, atomic, independently mergeable)

### Slice 1.1 — Scaffold public-safe Go project

- Status: ✅ done
- Owner: Codex
- Files (planned edits): repository policy, build tooling, Go module
- Files (do NOT edit): live asc authentication or runtime data
- Depends on: none
- Acceptance (EARS notation):
  - [ ] Repository contains no personal or machine-specific information.
  - [ ] Go build, lint, test, architecture, and public-safety gates are wired.
  - [ ] Documentation records the confirmed product behavior.

### Slice 2.1 — Implement the monitor

- Status: ✅ done
- Files (planned edits): `cmd/`, `internal/`, `config/`, tests, documentation

## 3. Blockers / open questions

- The `asc` global review-submission command still requires an app ID, so
  discovery must list apps and inspect each app's review submissions.

## 4. Recent decisions (append-only, newest first)

- 2026-07-24 — Use Go, Bubble Tea, and Lip Gloss (project owner).
- 2026-07-24 — Poll with asc only; persist runtime state outside the repo.

## 5. Next actions (ordered)

1. Implement and verify the TUI.
2. Add CI, installation docs, and release packaging.
3. Run the final completion gate and public-history audit.

## 6. Handoff note (fill when ending a session)

Initial implementation is active.
