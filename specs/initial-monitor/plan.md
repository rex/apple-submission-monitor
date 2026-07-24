# Plan — Initial monitor

## Phase 1 — Public-safe scaffold

- [ ] Configure Go module, quality gates, repository policy, and synthetic-only fixtures.
- [ ] Record the local `asc` architecture decision.

## Phase 2 — Monitoring engine

- [ ] Add typed domain records and tolerant `asc` JSON decoders.
- [ ] Discover apps, filter active submissions, and poll `asc status`.
- [ ] Persist cards atomically outside the repository.
- [ ] Detect transitions, retain terminal cards, and render announcement templates.

## Phase 3 — Terminal experience

- [ ] Render an adaptive equal-cell grid with large gradient app names.
- [ ] Use health colors for borders, badges, and animation.
- [ ] Support mouse acknowledgment plus keyboard navigation, links, refresh, and removal.
- [ ] Surface polling freshness, shortcuts, and sanitized degraded states.

## Phase 4 — Hardening

- [ ] Cover unit and fake-binary integration scenarios.
- [ ] Add CI, installation documentation, example configuration, and safety scanning.
- [ ] Run the completion gate, sign commits, push, and verify the public remote.

## Risks

| Risk | Mitigation |
|---|---|
| `asc` aggregate JSON evolves | Decode only stable fields and tolerate unknown fields. |
| Many apps create expensive scans | Limit concurrency and scan all apps less often than active polls. |
| Terminal animation wastes resources | Use a low-frequency shared tick and animate only unacknowledged cards. |
| Terminal submissions disappear upstream | Merge persisted cards before pruning and require explicit removal. |
| Runtime data leaks into git | Store it in the user state directory with mode `0600` and scan candidates. |

## Frozen decisions

- Go 1.24+, Bubble Tea, and Lip Gloss.
- Local-only and Docker-free.
- `asc` is the exclusive App Store Connect data source.
- Tests and public-safety validation are required.
