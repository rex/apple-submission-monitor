# AGENTS.md

<!-- Single source of truth for coding agents and contributors.
     CLAUDE.md and GEMINI.md are symlinks to this file. Keep under 150 lines. -->

## 1. Project snapshot
- **What**: A long-running terminal dashboard for active App Store Connect reviews.
- **Runtime**: Go 1.24+, Bubble Tea, Lip Gloss, and the external `asc` CLI.
- **Infra**: Local-only native CLI; no server, database, cloud resources, or Docker.
- **Privacy**: This is a public repository. Never commit live app data, credentials,
  usernames, home-directory paths, hostnames, or generated runtime state.
- **Non-goals**: Replacing `asc`, mutating App Store Connect, or monitoring every app.

## 2. Setup

```bash
# Replace per stack (lang-* skill fills these in at bootstrap time)
make setup
```

## 3. Commands the agent MUST run before declaring done

- `make lint`
- `make typecheck`
- `make test`  (if testing is enabled in `VIBE.yaml`)
- `make check-architecture`  (if declared)
- If `infra/**` changed: `terraform fmt -recursive infra/ && terraform validate`
- If `ansible/**` changed: `ansible-lint ansible/ && ansible-playbook --syntax-check`

## 4. Repo layout

```
cmd/           Thin CLI entry point
internal/      Private packages for ASC access, state, config, and the TUI
config/        Public-safe example configuration
scripts/       Repository validation and release helpers
specs/         Spec-driven dev artifacts — read specs/<active>/ first
docs/adr/      Architectural decisions (MADR) — read README.md index
.claude/       Subagents, slash commands, rules, hooks, MCP config
```

## 5. Code style (non-negotiable)

See `CONVENTIONS.md`. Linters enforce formatting and import order — **do not
write style rules here that the linter already checks.**

## 6. Testing policy

See `VIBE.yaml` (`quality_gates.tests`). If testing is `required`, every code
change gets a corresponding test update.

## 7. Security (hard stops)

- No secrets, PII, machine paths, live app metadata, or runtime state committed.
- All App Store Connect reads must execute `asc`; never call Apple's API directly.
- Never log or persist `asc` authentication material or raw command output.
- See `.claude/rules/security.md` for full checklist.

## 8. Architectural decisions

- Read `docs/adr/README.md` index before proposing layering / DB / auth / deploy changes.
- New decisions: create ADR (`docs/adr/template.md`), merge, THEN implement.

## 9. Things agents get wrong here

<!-- Update whenever an agent makes the same mistake twice. Start empty. -->

- (none yet)

## 10. Workflow

1. Read `PROGRESS.md` (session orientation) then `TASK_STATE.md` §0 and current slice.
2. Read `CONVENTIONS.md` before editing code.
3. Use `MAP.md` for exploration; `agent_docs/` for deep-dive detail.
4. Use Serena MCP symbolic tools over raw grep for cross-file work.
5. Use Context7 MCP for up-to-date library docs; don't rely on training data.

## 11. When ending a session

- Update `TASK_STATE.md` §6 (Handoff note).
- Update `PROGRESS.md` "Last session".
- Propose AGENTS.md updates for durable new facts — don't accumulate tribal
  knowledge in auto-memory.

## 12. Subdirectory AGENTS.md (precedence: nearest wins)

- `cmd/AGENTS.md` — entry-point constraints
- `internal/AGENTS.md` — package boundaries and privacy invariants
- `scripts/AGENTS.md` — validation script constraints

## 13. Composition with skills

This repo was bootstrapped with:
- `agentic-skeleton` (collaboration container — this file's shape)
- `lang-go` (Go code and test patterns)

Add new durable rules to the appropriate policy file. Transient context
goes in `TASK_STATE.md`; never here.
