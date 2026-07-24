---
status: accepted
date: 2026-07-24
deciders: project-owner
consulted: none
informed: contributors
tags: [cli, tui, app-store-connect]
---

# ADR-0001: Build a local Go TUI around asc

## Context

The monitor must remain responsive for long sessions, react to mouse and
keyboard input, animate status changes, and retrieve all App Store Connect data
through an already authenticated `asc` installation. It must also be safe to
develop in a public repository without fixtures or logs containing live data.

## Decision drivers

- Native, low-overhead behavior during long terminal sessions.
- First-class terminal layout, color, mouse, and animation primitives.
- A single local executable with no server or browser.
- Strict reuse of `asc` authentication and data retrieval.

## Considered options

1. Go with Bubble Tea and Lip Gloss.
2. Python with Textual.
3. A browser dashboard with a local service.

## Decision

Use Go with Bubble Tea and Lip Gloss. Execute `asc` as a typed external adapter;
do not call Apple APIs directly. Store generated state and configuration only
in the platform user directories.

## Consequences

- Positive: low idle resource use, a native binary, and a deterministic
  model-update-view architecture.
- Positive: authentication remains entirely owned by `asc`.
- Negative: rich graphics must adapt to terminal capability and dimensions.
- Negative: discovery requires inspecting apps individually because the
  installed `asc` global submission command still requires an app ID.

## Validation

- Race-enabled tests cover transition and persistence behavior.
- A fake `asc` integration test validates fixed argument vectors and schemas.
- Repository safety checks reject machine paths, likely credentials, and
  personal contact data.
