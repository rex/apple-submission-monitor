#!/usr/bin/env bash
set -euo pipefail

# Scan tracked files and non-ignored candidates so the gate protects the first
# commit as well as later changes.
patterns=(
  '/Users/[^/[:space:]]+/'
  '/home/[^/[:space:]]+/'
  'AKIA[0-9A-Z]{16}'
  '-----BEGIN (RSA |EC |OPENSSH )?PRIVATE KEY-----'
  '(ASC_PRIVATE_KEY|ASC_KEY_ID|ASC_ISSUER_ID|GITHUB_TOKEN)[[:space:]]*[:=][[:space:]]*["'\'']?[[:alnum:]]'
  '[[:alnum:]._%+-]+@[[:alnum:].-]+\.[[:alpha:]]{2,}'
)

failed=0
while IFS= read -r -d '' file; do
  [[ -f "$file" ]] || continue
  [[ "$file" == "scripts/check-public-safety.sh" ]] && continue
  grep -Iq . "$file" || continue
  for pattern in "${patterns[@]}"; do
    if grep -Eq -- "$pattern" "$file"; then
      printf 'public-safety: prohibited content pattern in %s\n' "$file" >&2
      failed=1
      break
    fi
  done
done < <(git ls-files --cached --others --exclude-standard -z)

if (( failed != 0 )); then
  exit 1
fi

ignored_go=$(git ls-files --others --ignored --exclude-standard -- "*.go")
if [[ -n "$ignored_go" ]]; then
  printf 'public-safety: Go source is hidden by .gitignore:\n%s\n' "$ignored_go" >&2
  exit 1
fi

printf 'public-safety: no secrets, email addresses, or machine paths found\n'
