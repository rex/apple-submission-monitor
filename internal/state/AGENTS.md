# AGENTS — State

- Persist versioned JSON atomically with mode `0600`.
- State lives in the platform user-config directory, never in the repository.
- Do not store credentials or raw `asc` responses.
- Treat unknown future file versions as errors.
