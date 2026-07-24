# AGENTS — Scripts

Scripts provide deterministic local and CI validation.

- Never print environment values, credentials, live app metadata, or home paths.
- Keep checks portable across macOS and Linux unless explicitly macOS-only.
- Destructive actions require explicit targets and confirmation.
- Security checks must inspect tracked content, not ignored runtime data.
