#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT_DIR}"

if ! command -v swag >/dev/null 2>&1; then
  echo "[swagger] swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest" >&2
  exit 1
fi

swag init -g apps/did_service_endpoint/main.go -o apps/did_service_endpoint/docs
