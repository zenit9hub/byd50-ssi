#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT_DIR}"

if ! command -v npx >/dev/null 2>&1; then
  echo "[swagger-lint] npx not found. Install Node.js first." >&2
  exit 1
fi

npx @redocly/cli lint apps/did_service_endpoint/docs/swagger.json
