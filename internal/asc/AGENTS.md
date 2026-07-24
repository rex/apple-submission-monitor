# AGENTS — asc adapter

This package is the exclusive App Store Connect retrieval boundary.

- Execute `asc` directly with fixed argument slices; never invoke a shell.
- Enforce context cancellation, timeouts, and bounded discovery concurrency.
- Decode only required fields and ignore compatible upstream additions.
- Return sanitized errors without raw output, app identifiers, or local paths.
