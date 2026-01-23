#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT_DIR}"

export PATH="$(go env GOPATH)/bin:${PATH}"

if ! command -v swag >/dev/null 2>&1; then
  echo "[swagger] swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest" >&2
  exit 1
fi

swag init -g main.go -d apps/did_service_endpoint -o apps/did_service_endpoint/docs
