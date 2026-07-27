# Spec — Initial monitor

## Summary

Build a local, long-running terminal dashboard for active App Store Connect
review submissions. All remote data comes from the installed `asc` CLI.

## Goals

- Make every active submission readable at a glance in an evenly divided grid.
- Announce and animate real status transitions until the user acknowledges them.
- Retain terminal submissions until the user explicitly removes them.
- Remain stable during long sessions, terminal resizes, transient failures, and restarts.

## Non-goals

- Mutating App Store Connect data or replacing `asc` authentication.
- Displaying apps without active review submissions.
- Remote hosting, Docker, telemetry, or account analytics.

## Acceptance criteria

### Always

- The system shall retrieve App Store Connect data only by executing `asc`.
- The system shall keep credentials, live app data, and machine identifiers out
  of repository files, fixtures, diagnostics, and command arguments.
- The system shall color each card with the normalized `summary.health` value
  returned by `asc status`.
- The system shall render total time since submission in the same large FIGlet
  scale as the app name and classify it green through two days, yellow from two
  to five days, and deteriorated blood-red at five days or later.
- The system shall show App ID and trustworthy build, blocker, review-detail,
  and in-flight checks retrieved through `asc`.

### Events

- When a status changes, the system shall persist the transition, speak the
  configured announcement with `/usr/bin/say` and no voice flag, and animate
  the card until acknowledgment.
- When an animated card is clicked or selected and acknowledged, the system
  shall stop its alert flashing and retain its stable health color while
  ambient hero gradients continue.
- When an acknowledged terminal card is removed, the system shall delete only
  that persisted card.

### State

- While a submission remains active, the system shall poll it on a configurable
  interval without overlapping polling cycles.
- While `asc` is unavailable or a poll fails, the system shall retain the last
  successful data and display a sanitized error.
- While the terminal is resized, the system shall preserve evenly sized cards,
  selection, acknowledgment state, and readable fallback content.

### Failure behavior

- If no submissions are active, the system shall render a helpful empty state.
- If an approved or otherwise terminal submission disappears from `asc`, the
  system shall retain it until acknowledgment and explicit removal.
- If speech is unavailable, the system shall continue monitoring and show a
  non-fatal warning.

## Success measures

- Race-enabled tests remain green across polling, persistence, and transition logic.
- A synthetic end-to-end test proves discovery, change detection, acknowledgment,
  retention, removal, and no-voice speech invocation.
- The public-safety gate finds no secrets, personal data, or machine paths.
