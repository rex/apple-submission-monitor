#!/usr/bin/env python3
"""Enforce the repository's Go source-file size and public-surface limits."""

from __future__ import annotations

import re
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parent.parent
SOURCE_ROOTS = (ROOT / "cmd", ROOT / "internal")
SOFT_LIMIT = 250
HARD_LIMIT = 400
PUBLIC_FUNCTION_LIMIT = 8
PUBLIC_FUNCTION = re.compile(r"^func (?:\([^)]*\) )?([A-Z][A-Za-z0-9_]*)\(")


def source_files() -> list[Path]:
    """Return application Go files in deterministic path order."""
    return sorted(
        path
        for root in SOURCE_ROOTS
        if root.exists()
        for path in root.rglob("*.go")
        if not path.name.endswith(".generated.go")
    )


def main() -> int:
    """Report architecture warnings and fail on hard policy violations."""
    failed = False
    warnings = 0
    files = source_files()

    for path in files:
        lines = path.read_text(encoding="utf-8").splitlines()
        relative = path.relative_to(ROOT)
        if len(lines) > HARD_LIMIT:
            print(f"architecture: {relative} has {len(lines)} lines (hard {HARD_LIMIT})")
            failed = True
        elif len(lines) > SOFT_LIMIT:
            print(f"architecture: warning: {relative} has {len(lines)} lines (soft {SOFT_LIMIT})")
            warnings += 1

        public_functions = sum(bool(PUBLIC_FUNCTION.match(line)) for line in lines)
        if public_functions > PUBLIC_FUNCTION_LIMIT:
            print(
                f"architecture: {relative} exposes {public_functions} functions "
                f"(maximum {PUBLIC_FUNCTION_LIMIT})"
            )
            failed = True

    if failed:
        return 1
    print(f"architecture: {len(files)} Go files checked, {warnings} warnings")
    return 0


if __name__ == "__main__":
    sys.exit(main())
