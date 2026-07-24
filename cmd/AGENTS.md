# AGENTS — Commands

Command packages are thin composition roots.

- Parse flags, load configuration, wire dependencies, and start the TUI.
- Keep monitoring, persistence, speech, and rendering logic in `internal/`.
- Exit startup failures with concise messages that contain no credentials,
  app metadata, machine paths, or raw `asc` output.
