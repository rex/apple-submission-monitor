#!/usr/bin/env bash
# SessionStart hook — injects git state + TASK_STATE.md + PROGRESS.md into the
# fresh session context so the agent has its bearings before the first prompt.
#
# Fires on: startup | resume | clear
# Emits:    JSON with hookSpecificOutput.additionalContext

set -euo pipefail

cd "${CLAUDE_PROJECT_DIR:-.}"

ctx="## Git state\n"
ctx+="Branch: $(git branch --show-current 2>/dev/null || echo 'detached')\n"
ctx+="Uncommitted: $(git status --porcelain 2>/dev/null | wc -l | tr -d ' ') files\n\n"

if command -v git >/dev/null 2>&1; then
  ctx+="## Recent commits\n$(git log --oneline -10 2>/dev/null || echo 'no history')\n\n"
fi

if [ -f TASK_STATE.md ]; then
  ctx+="## Current task\n$(head -80 TASK_STATE.md)\n\n"
fi

if [ -f PROGRESS.md ]; then
  ctx+="## Progress\n$(cat PROGRESS.md)\n\n"
fi

if [ -f .claude/SESSION_NOTES.md ]; then
  ctx+="## Previous session notes\n$(cat .claude/SESSION_NOTES.md)\n"
fi

jq -n --arg c "$ctx" \
  '{hookSpecificOutput:{hookEventName:"SessionStart", additionalContext:$c}}'
