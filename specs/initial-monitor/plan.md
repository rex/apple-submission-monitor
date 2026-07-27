# Plan — Initial monitor

## Phase 1 — Public-safe scaffold

- [x] Configure Go module, quality gates, repository policy, and synthetic-only fixtures.
- [x] Record the local `asc` architecture decision.

## Phase 2 — Monitoring engine

- [x] Add typed domain records and tolerant `asc` JSON decoders.
- [x] Discover apps, filter active submissions, and poll `asc status`.
- [x] Persist cards atomically outside the repository.
- [x] Detect transitions, retain terminal cards, and render announcement templates.

## Phase 3 — Terminal experience

- [x] Render an adaptive equal-cell grid with large gradient app names.
- [x] Pair app names with same-scale, age-classified submission-time heroes.
- [x] Distribute App ID and ASC-derived health checks across wide metadata rails.
- [x] Use health colors for borders, badges, and animation.
- [x] Support mouse acknowledgment plus keyboard navigation, links, refresh, and removal.
- [x] Surface polling freshness, shortcuts, and sanitized degraded states.

## Phase 4 — Hardening

- [x] Cover unit and fake-binary integration scenarios.
- [x] Add installation documentation, example configuration, and safety scanning.
- [x] Run the completion gate, sign commits, push, and verify the public remote.

## Risks

| Risk | Mitigation |
|---|---|
| `asc` aggregate JSON evolves | Decode only stable fields and tolerate unknown fields. |
| Many apps create expensive scans | Limit concurrency and scan all apps less often than active polls. |
| Terminal animation wastes resources | Cache both FIGlet masks and recolor them on one low-frequency shared tick. |
| Terminal submissions disappear upstream | Merge persisted cards before pruning and require explicit removal. |
| Runtime data leaks into git | Store it in the user state directory with mode `0600` and scan candidates. |

## Frozen decisions

- Go 1.24+, Bubble Tea, and Lip Gloss.
- Local-only and Docker-free.
- `asc` is the exclusive App Store Connect data source.
- Tests and public-safety validation are required.
