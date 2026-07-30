# Design — Initial monitor

## Data flow

```text
asc apps list
      │
      ├─ periodic discovery ──> asc review submissions-list
      │                               │ active app IDs
      ▼                               ▼
persisted cards <── transition engine <── asc status
      ▲                                      │ completed review
      └──── asc review status/history <──────┘
      │                    │
      │                    └─> announcement template ─> /usr/bin/say
      ▼
Bubble Tea model ──> equal-cell grid ──> terminal
```

## Components

- `internal/asc`: Executes fixed argument vectors, enforces timeouts, and decodes
  stable fields. It never invokes a shell.
- `internal/monitor`: Merges discovery results with persisted terminal cards,
  detects meaningful status changes, and produces transition events.
- `internal/state`: Atomically writes versioned JSON with private permissions in
  the platform user-state directory.
- `internal/config`: Loads YAML defaults, validates durations/templates, and
  creates an editable local file on first launch.
- `internal/speech`: Expands safe templates and calls `/usr/bin/say` without
  any voice option.
- `internal/tui`: Owns selection, layout, animation, mouse hit-testing, and views.

## Status identity and transitions

A card is keyed by review-submission ID, not app ID, so consecutive submissions
for the same app do not overwrite history. A transition fingerprint contains
the authoritative user-visible status and terminal identity. Transport flags,
next actions, health-only updates, and other metadata do not produce
same-label announcements.

Persisted cards continue polling after active discovery drops them. Completed
reviews are enriched from review status/history until an outcome is known. Once
a monitored card becomes terminal, it remains until the user acknowledges it
and invokes removal.

## Security and privacy

- Commands use `exec.CommandContext` with argument slices and no shell.
- Runtime state and generated config are mode `0600`.
- Diagnostics exclude raw output and identifiers.
- Tests use synthetic app names, IDs, URLs, and a fake `asc` executable.
- Repository validation scans tracked and unignored candidate content.

## Performance

Active polls default to 30 seconds. Full app discovery defaults to five minutes.
Discovery workers are capped at four; submitted-build checks run only during
discovery, while known terminal cards use bounded review status/history calls.
A shared 400 ms animation tick recolors cached hero masks, while only
unacknowledged cards flash their border and background.
