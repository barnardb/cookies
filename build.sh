#!/usr/bin/env bash
set -euo pipefail

[ $# -eq 1 ] && grep -qE '^\d+\.\d+\.\d+$' <<<"$1" || {
  echo "usage: $0 <semantic-version>" >&2
  exit 2
}

version="$1"

go build -ldflags="-X 'main.version=${version}'"
