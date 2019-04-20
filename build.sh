#!/usr/bin/env bash
set -euo pipefail

export GO111MODULE=on
exec go build
