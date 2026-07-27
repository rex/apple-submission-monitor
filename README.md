# Apple Submission Monitor

A polished terminal dashboard for watching active App Store Connect submissions
move through Apple's review process.

## Status

🟢 Usable beta · macOS-first · local-only

## Why this exists

App Store reviews can take hours or days, and repeatedly checking a browser
makes status changes easy to miss. This tool uses the authenticated `asc` CLI
as its only data source, keeps active submissions visible in evenly sized
status-colored cards, and announces transitions with the macOS default voice.

## What it does

- Discover only apps with an active review submission.
- Poll status through `asc status` and use `summary.health` for card colors.
- Stack full-width review cards vertically, pairing each large filled FIGlet app
  name with an equally large, right-aligned total submission age.
- Distort red-health/rejected banners into a blood-dripping horror treatment.
- Color the animated submission-age hero green through two days, gold from two
  to five days, then deteriorated and blood-red from five days onward.
- Distribute status, version, platform, App ID, bundle ID, timestamps, links,
  and explicit build/blocker/review checks across three bold metadata rails.
- Flash changed cards until they are clicked or selected and acknowledged.
- Announce transitions via `/usr/bin/say` without selecting a voice.
- Retain completed submissions until acknowledgment and explicit removal.
- Store editable announcement templates and runtime state outside the repo.

## Quick start

Prerequisites:

- macOS
- Go 1.24 or newer
- [`asc`](https://github.com/rorkai/App-Store-Connect-CLI) 3.1 or newer,
  authenticated with `asc auth status --validate`

```bash
git clone https://github.com/rex/apple-submission-monitor.git
cd apple-submission-monitor
make setup
make build
./bin/apple-submission-monitor
```

Install the versioned binary into the Go bin directory:

```bash
make install-cli
apple-submission-monitor
```

The first launch creates a private local configuration file and shows its path
in the footer. It never copies `asc` credentials.

## Controls

| Input | Action |
|---|---|
| Arrow keys or `h`/`j`/`k`/`l` | Select a card |
| Click or `Enter`/`Space` | Acknowledge a changed card and stop its flashing |
| `o` | Open the selected review in App Store Connect |
| `D` or `Delete` | Remove an acknowledged terminal/retained card |
| `r` | Refresh and discover now |
| `?` | Toggle help |
| `q` or `Ctrl-C` | Quit cleanly |

## Configuration

The generated YAML contains polling, timeout, animation, executable, and
announcement settings. The committed
[`config/config.example.yaml`](config/config.example.yaml) documents every
field.

```yaml
announcements:
  approved: "Great news. {{.AppName}} has been approved."
  rejected: "Attention. {{.AppName}} was rejected. Its status is {{.NewStatus}}."
  status_changed: "{{.AppName}} changed from {{.OldStatus}} to {{.NewStatus}}."
```

Available template fields are `.AppName`, `.OldStatus`, and `.NewStatus`.
Announcements execute `say` with the rendered sentence as the only argument;
the app never specifies a voice.

CLI options:

```text
--config PATH   Use an explicit local YAML file
--no-speech     Disable spoken announcements
--version       Print the binary version
```

See [Operations](docs/OPERATIONS.md) for long-session behavior, persistence,
failure recovery, and tuning.

## Architecture

```
asc CLI → monitor engine → persisted snapshot → Bubble Tea dashboard
```

- Full discovery runs less frequently than active-card polling, uses a bounded
  worker pool, and enriches cards with submitted-build and review-detail checks.
- Status changes compare a stable fingerprint, preventing metadata-only flashes.
- Prepared dual-hero FIGlet masks are cached; animation recolors the app name
  and submission age without rebuilding fonts or strobing an acknowledged
  card's background.
- State writes use private permissions and atomic replacement.
- Terminal results remain until acknowledgment and explicit removal.

## Development

```bash
make setup
make dev
make validate
```

`make validate` runs formatting/static analysis, race-enabled tests, coverage,
architecture limits, the public-safety scan, module verification, reachable
vulnerability analysis, and the version/changelog gate. `VIBE.yaml` is the
machine-readable repository policy. Docker and remote deployment are
intentionally out of scope. External CI is not coupled; `make validate` is the
authoritative local completion gate.

## Privacy

This is a public repository. Example data is synthetic. The program never
stores credentials, and generated configuration, state, logs, app names, IDs,
and other live `asc` output must never be committed. The safety gate also fails
if a Go source file is accidentally hidden by `.gitignore`.
