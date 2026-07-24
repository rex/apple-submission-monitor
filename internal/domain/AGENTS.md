# AGENTS — Domain

Domain types contain no I/O and no framework dependencies.

- Keep transition identity stable and deterministic.
- Unknown upstream values must degrade safely rather than fail decoding.
- JSON tags are part of the persisted-state contract.
