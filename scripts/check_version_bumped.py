#!/usr/bin/env python3
"""check_version_bumped.py — gate that blocks commits without a version bump.

Wired into:
  - `make validate` (skeleton Makefile target)
  - `.pre-commit-config.yaml` (humans)
  - `auto-commit.sh` (Claude Code autonomous commits)

Rules enforced:
  1. VERSION file exists and contains a semver string.
  2. VERSION differs from HEAD's VERSION (unless HEAD doesn't exist —
     bootstrap exemption for the first commit).
  3. CHANGELOG.md has a `## [<NEW_VERSION>] — ` header matching the
     current VERSION.

Exit codes:
  0  ok
  1  version not bumped
  2  changelog missing or mismatched
  3  malformed version
"""

from __future__ import annotations

import os
import re
import subprocess
import sys
from pathlib import Path


_TTY = sys.stderr.isatty() and "NO_COLOR" not in os.environ
_R = "\033[31m" if _TTY else ""
_G = "\033[32m" if _TTY else ""
_X = "\033[0m" if _TTY else ""

SEMVER_RE = re.compile(r"^[0-9]+\.[0-9]+\.[0-9]+$")


def fail(msg: str, code: int) -> None:
    print(f"{_R}✗ version-gate:{_X} {msg}", file=sys.stderr)
    sys.exit(code)


def head_version() -> str | None:
    """Return HEAD's VERSION contents, or None if no HEAD or no VERSION at HEAD."""
    try:
        result = subprocess.run(
            ["git", "rev-parse", "--git-dir"],
            capture_output=True, check=False,
        )
        if result.returncode != 0:
            return None  # not in a git repo
    except FileNotFoundError:
        return None  # git not available

    # Bootstrap exemption: no HEAD yet
    head_check = subprocess.run(
        ["git", "rev-parse", "HEAD"],
        capture_output=True, check=False,
    )
    if head_check.returncode != 0:
        return None

    show = subprocess.run(
        ["git", "show", "HEAD:VERSION"],
        capture_output=True, text=True, check=False,
    )
    if show.returncode != 0:
        return None  # VERSION didn't exist at HEAD
    return show.stdout.strip()


def main() -> int:
    version_file = Path("VERSION")
    if not version_file.is_file():
        fail("VERSION file missing at repo root. Run: scripts/bump_version.py patch", 1)

    current = version_file.read_text().strip()
    if not SEMVER_RE.match(current):
        fail(f"VERSION is not semver: '{current}'", 3)

    prev = head_version()
    if prev is not None and prev == current:
        fail(
            f"VERSION unchanged ({current}). Bump before committing.\n"
            "    Run: make bump-patch | bump-minor | bump-major",
            1,
        )

    changelog = Path("CHANGELOG.md")
    if not changelog.is_file():
        fail("CHANGELOG.md missing. bump_version.py creates it — re-run the bump.", 2)

    changelog_text = changelog.read_text()
    # Match `## [X.Y.Z] — ` or `## [X.Y.Z] -` at any line
    pattern = rf"^## \[{re.escape(current)}\][[:space:]]+[—-]"
    # Python's re doesn't support [:space:] — use \s
    pattern = pattern.replace("[[:space:]]+", r"\s+")
    if not re.search(pattern, changelog_text, re.MULTILINE):
        fail(
            f"CHANGELOG.md has no entry for version {current}. "
            f"Expected: '## [{current}] — <date> — Agent: <name>'",
            2,
        )

    print(f"{_G}✓ version-gate:{_X} VERSION={current}, CHANGELOG has matching entry")
    return 0


if __name__ == "__main__":
    sys.exit(main())
