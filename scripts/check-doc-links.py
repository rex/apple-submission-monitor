#!/usr/bin/env python3
"""Validate repository-local Markdown links without network access."""

from __future__ import annotations

import re
import sys
from pathlib import Path
from urllib.parse import unquote


ROOT = Path(__file__).resolve().parent.parent
LINK = re.compile(r"\[[^\]]+\]\(([^)]+)\)")
REMOTE_PREFIXES = ("http://", "https://", "mailto:", "#")


def markdown_files() -> list[Path]:
    """Return committed and candidate Markdown files in stable order."""
    excluded = {".git", ".serena"}
    return sorted(
        path
        for path in ROOT.rglob("*.md")
        if not any(part in excluded for part in path.relative_to(ROOT).parts)
    )


def main() -> int:
    """Report missing local link targets."""
    failures: list[tuple[Path, str]] = []
    checked = 0
    for source in markdown_files():
        content = source.read_text(encoding="utf-8")
        for match in LINK.finditer(content):
            value = match.group(1).strip().strip("<>")
            if not value or value.startswith(REMOTE_PREFIXES):
                continue
            target_value = unquote(value.split("#", 1)[0])
            if not target_value:
                continue
            checked += 1
            target = (source.parent / target_value).resolve()
            if not target.exists():
                failures.append((source.relative_to(ROOT), value))

    for source, value in failures:
        print(f"docs: missing link from {source}: {value}")
    if failures:
        return 1
    print(f"docs: {checked} local links checked")
    return 0


if __name__ == "__main__":
    sys.exit(main())
