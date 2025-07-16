#!/usr/bin/env bash
# task_status.sh — quick helper to change the status tag in a task file
# Usage: ./scripts/task_status.sh <status> <id>
# Example: ./scripts/task_status.sh done 07

set -euo pipefail

if [[ $# -ne 2 ]]; then
  echo "Usage: $0 <status> <task-id>" >&2
  echo "Accepted status values: ready, in_progress, blocked, done" >&2
  exit 1
fi

new_status="$1"
task_id="$2"

case "$new_status" in
  ready|in_progress|blocked|done) ;;
  *) echo "Invalid status: $new_status" >&2; exit 1;;
esac

# locate task file (pattern <id>-*.md)
file_path=$(ls "$(dirname "$0")/../tasks/"${task_id}-*.md 2>/dev/null || true)
if [[ -z "$file_path" ]]; then
  echo "Task file for id $task_id not found" >&2
  exit 1
fi

# Update the first '**Status:** [...]' line in place
# shellcheck disable=SC2016
sed -i -E '0,/\*\*Status:\*\* \[[^]]+\]/s//**Status:** ['"$new_status"']/' "$file_path"

echo "Updated $file_path → [$new_status]" 