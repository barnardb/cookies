#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

[ $# -eq 1 ] && grep -qE '^\d+\.\d+\.\d+$' <<<"$1" || {
  echo "usage: $0 <semantic-version>" >&2
  exit 2
}

set -x

version="$1"
tag="v${version}"

export GOOS=darwin
export GOARCH=amd64
./build.sh "${version}"
archive="cookies_${tag}_${GOOS}_${GOARCH}.gz"
gzip --stdout cookies >"${archive}"
gh release create --draft --generate-notes "${tag}" "${archive}"
