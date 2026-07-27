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









## [0.6.0] — 2026-07-27 — Agent: Codex
### Added
- Pair each app-name hero with an equally large, right-aligned total submission
  age using a shared FIGlet font and independently animated gradient.
- Classify submission age as green through 48 hours, gold until 120 hours, and
  deteriorated blood-red from 120 hours onward.
- Show App ID and ASC-derived submitted-build, blocker, review-detail, and
  in-flight checks across bold full-width metadata rails.
### Changed
- Enrich diagnostic metadata only during five-minute discovery cycles while
  preserving lightweight 30-second status polling.
- Treat diagnostic refreshes as metadata so they never create false
  status-change alerts or speech announcements.

## [0.5.0] — 2026-07-27 — Agent: Codex
### Added
- Pin a live current-status timer to the far-right edge of every card's status
  rail with an independent reverse-flow mint, coral, and rose gradient.
### Changed
- Move version and platform metadata to the review row so status timing remains
  consistently glanceable without competing with the FIGlet hero.

## [0.4.0] — 2026-07-27 — Agent: Codex
### Added
- Render app names with large, filled FIGlet typography, multicolor animated
  gradients, and a down-right dimensional shadow.
- Give rejected or otherwise red-health cards a distorted poison-style banner
  with deterministic blood drips.
### Changed
- Confine ambient motion to the banner glyphs so acknowledged waiting cards no
  longer strobe their entire background.
- Cache prepared FIGlet masks across animation frames for stable long-session
  performance.

## [0.3.4] — 2026-07-27 — Agent: Codex
### Changed
- Stack full-width review cards vertically and replace the block font with a
  readable, layered gradient wordmark.
- Render the elapsed time since each card's last meaningful status change.
- Give waiting reviews a calm, rich blue palette with a slow ambient drift;
  unacknowledged changes still flash until acknowledged.

## [0.3.3] — 2026-07-24 — Agent: Codex
### Fixed
- Remove the unused `.env` scaffold from `make setup`, eliminating its shell
  parse failure.
- Execute the public setup entrypoint from the validation gate to prevent
  onboarding regressions.

## [0.3.2] — 2026-07-24 — Agent: Codex
### Changed
- Complete installation, controls, configuration, privacy, and long-session
  operating documentation.
- Add local Markdown-link validation, module verification, and pinned
  reachable-vulnerability scanning to the completion gate.
- Add a versioned `make install-cli` workflow.
### Security
- Strip terminal control characters from metadata and validate HTTPS links
  before emitting terminal hyperlink sequences.
### Fixed
- Let the final verification gate validate a clean committed release while
  still requiring pending changes to bump their version.

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
