#!/usr/bin/env bash
set -euo pipefail

# Tags and pushes a release. The Release workflow (.github/workflows/release.yml) then builds
# binaries and creates a draft GitHub release to review and publish.

cd "$(dirname "$0")"

[ $# -eq 1 ] && grep -qE '^[0-9]+\.[0-9]+\.[0-9]+$' <<<"$1" || {
  echo "usage: $0 <semantic-version>" >&2
  exit 2
}

version="$1"
tag="v${version}"

fail() {
  echo "error: $*" >&2
  exit 1
}

[ -z "$(git status --porcelain)" ] || fail "the working tree isn't clean"
[ "$(git symbolic-ref --short HEAD)" = main ] || fail "not on main"
git fetch --quiet --tags origin main
[ "$(git rev-parse HEAD)" = "$(git rev-parse origin/main)" ] || fail "main isn't in sync with origin/main"
! git rev-parse --quiet --verify "refs/tags/${tag}" >/dev/null || fail "tag ${tag} already exists"

go test ./...

set -x
git tag --annotate "${tag}" --message "cookies ${version}"
git push origin "${tag}"
{ set +x; } 2>/dev/null

echo "Pushed ${tag}. Watch the Release workflow, then review and publish the draft release:"
echo "  https://github.com/barnardb/cookies/actions/workflows/release.yml"
echo "  https://github.com/barnardb/cookies/releases"
echo "After publishing, run ./update.sh in barnardb/homebrew-cookies."
