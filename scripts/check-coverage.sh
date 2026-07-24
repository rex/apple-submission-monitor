#!/usr/bin/env bash
set -euo pipefail

profile=${1:?coverage profile path is required}
minimum=${2:?minimum coverage percentage is required}

actual=$(go tool cover -func="$profile" | awk '/^total:/ {gsub("%", "", $3); print $3}')
if [[ -z "$actual" ]]; then
  printf 'coverage: unable to read total from %s\n' "$profile" >&2
  exit 1
fi

if ! awk -v actual="$actual" -v minimum="$minimum" 'BEGIN { exit !(actual >= minimum) }'; then
  printf 'coverage: %s%% is below required %s%%\n' "$actual" "$minimum" >&2
  exit 1
fi

printf 'coverage: %s%% (minimum %s%%)\n' "$actual" "$minimum"
