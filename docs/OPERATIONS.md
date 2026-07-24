# Operations

Apple Submission Monitor is designed to remain open in a dedicated terminal
window for days at a time. It has no server, daemon, telemetry, or database.

## Startup

Validate the external dependency before the first run:

```bash
asc version
asc auth status --validate
apple-submission-monitor
```

On first launch, the app creates `config.yaml` and `state.json` in the
platform user-config directory. The exact configuration path appears in the
footer. Both files and their directory use private permissions.

## Polling model

- Active cards refresh every 30 seconds by default.
- Full app discovery runs every five minutes by default.
- A full discovery lists apps through `asc`, checks each app's review
  submissions with bounded concurrency, and renders only submitted,
  non-complete reviews.
- Polling cycles never overlap. Manual refresh waits for any local state action
  to finish.
- A shared 400 ms animation timer runs only while a card is unacknowledged.

Tune the values in the generated YAML:

```yaml
poll_interval: 30s
discovery_interval: 5m
command_timeout: 25s
animation_interval: 400ms
discovery_workers: 4
```

The poll interval cannot be lower than five seconds. Discovery cannot run more
frequently than polling, and worker count must remain between 1 and 16.

## Status changes

A transition occurs when health, review state, App Store state, next action, or
the in-flight flag changes. App names, versions, links, and other metadata may
refresh without generating duplicate alerts.

For each transition the app:

1. Persists the new card immediately.
2. Marks it unacknowledged.
3. Starts the shared pulse animation.
4. Renders the configured announcement.
5. Executes `say` with only that sentence.

Click the card or press `Enter` to acknowledge it. Approved, rejected, and
otherwise terminal cards remain persisted even after `asc` stops listing them.
After acknowledgment, press `D` to remove one.

## Degraded operation

- If a poll fails, the dashboard retains its last successful snapshot and
  displays a sanitized warning.
- If speech fails, visual monitoring continues.
- If a link cannot open, the card remains available and monitoring continues.
- If the terminal becomes too small, resize it to at least 36 columns by nine
  rows; selection and state are preserved.
- `q`, `Ctrl-C`, SIGINT, and SIGTERM restore the terminal before exit.

## Privacy

The application does not log raw `asc` output and does not persist credentials.
Runtime state necessarily contains the app metadata displayed in the dashboard,
so it stays outside the repository with mode `0600`. Repository tests use a
fake executable and synthetic records.

Before publishing any change, run:

```bash
make check-public-safety
make validate
```
