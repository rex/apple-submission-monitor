# Changelog

All notable changes to this project are documented here. This project
follows [Semantic Versioning](https://semver.org/) and
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

Format:

```markdown
## [X.Y.Z] — YYYY-MM-DD — Agent: <name>
### Added | Changed | Fixed | Removed | Deprecated | Security
- <what changed, in imperative voice>
```

**Every commit requires a version bump and a matching entry here.** The
`scripts/check_version_bumped.py` gate enforces this; `auto-commit.sh`
calls `scripts/bump_version.py <level>` before commit.

**Bump level guidance** (agent decides per slice):
- `patch` — bug fix, documentation change, refactor with no behavior change
- `minor` — new feature, new public API, any backward-compatible addition
- `major` — breaking change, removal, incompatible behavior change

Append new entries at the top. One entry per commit (same cadence as
version bumps).

---



## [0.3.1] — 2026-07-24 — Agent: Codex
### Fixed
- Fix command-source packaging and ignored-source detection.

## [0.3.0] — 2026-07-24 — Agent: Codex
### Added
- Adaptive equal-cell Bubble Tea dashboard with health-colored cards, large
  gradient app names, metadata, timing, and App Store Connect links.
- Mouse and keyboard acknowledgment, shared pulse animation, manual refresh,
  retained-outcome removal, and resilient long-lived polling.
- Native CLI wiring, safe HTTPS opening, versioned builds, and TUI coverage.
### Fixed
- Version bump tooling now inserts entries after the changelog preamble instead
  of modifying its fenced format example.

## [0.2.0] — 2026-07-24 — Agent: Codex
### Added
- Typed, bounded `asc` discovery and status retrieval with sanitized failures.
- Atomic private state, configurable announcement templates, status transition
  detection, acknowledgment, retained outcomes, and explicit removal.
- Synthetic unit and fake-process integration coverage for the monitoring core.

## [0.1.0] — 2026-07-24 — Agent: Codex
### Added
- Public-safe Go CLI scaffold with durable policy and quality gates.
- Agent collaboration files, semantic versioning, and project documentation.
