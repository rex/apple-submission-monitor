# Apple Submission Monitor

A polished terminal dashboard for watching active App Store Connect submissions
move through Apple's review process.

## Status

🟡 Early development

## Why this exists

App Store reviews can take hours or days, and repeatedly checking a browser
makes status changes easy to miss. This tool uses the authenticated `asc` CLI
as its only data source, keeps active submissions visible in evenly sized
status-colored cards, and announces transitions with the macOS default voice.

## Planned experience

- Discover only apps with an active review submission.
- Poll status through `asc status` and use `summary.health` for card colors.
- Flash changed cards until they are clicked or selected and acknowledged.
- Announce transitions via `/usr/bin/say` without selecting a voice.
- Retain completed submissions until acknowledgment and explicit removal.
- Store editable announcement templates and runtime state outside the repo.

## Architecture

```
asc CLI → monitor engine → persisted snapshot → Bubble Tea dashboard
```

- Go 1.24 or newer
- `asc` 3.1 or newer, already authenticated
- macOS for spoken announcements; the dashboard remains usable elsewhere

## Development

```bash
make setup
make dev
make validate
```

`VIBE.yaml` is the machine-readable repository policy. Testing is required;
Docker and remote deployment are intentionally out of scope.

## Privacy

This is a public repository. Example data is synthetic. The program never
stores credentials, and generated configuration, state, logs, app names, IDs,
and other live `asc` output must never be committed.
