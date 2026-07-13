#!/usr/bin/env bash
set -euo pipefail

config_path="${1:-${CONFIG:-config.yaml}}"
work_dir="$(mktemp -d)"
trap 'rm -rf "$work_dir"' EXIT

go build -o "$work_dir/fnctl" ./cmd/fnctl

first="$("$work_dir/fnctl" migrate --config "$config_path")"
second="$("$work_dir/fnctl" migrate --config "$config_path")"
status="$("$work_dir/fnctl" db-status --config "$config_path")"

grep -Fq 'migrations_total=6' <<<"$first"
grep -Fq 'migrations_total=6 applied=0 skipped=6' <<<"$second"
grep -Fq 'migrations=6 tables=22 empty_table_comments=0 empty_column_comments=0 enabled_policies=6' <<<"$status"
grep -Fq 'read_only=false' <<<"$status"

printf '%s\n' "$first"
printf '%s\n' "$second"
printf '%s\n' "$status"
