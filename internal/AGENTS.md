# AGENTS — Internal packages

Internal packages contain all application behavior.

- Keep package boundaries narrow and dependencies explicit.
- Every blocking operation accepts `context.Context`.
- Never expose raw `asc` output, credentials, live identifiers, or machine paths.
- Tests use synthetic records and fake executables only.
- Keep source files below the repository line limits in `VIBE.yaml`.
