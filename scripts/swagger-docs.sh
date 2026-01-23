#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT_DIR}"

if ! command -v npx >/dev/null 2>&1; then
  echo "[swagger-docs] npx not found. Install Node.js first." >&2
  exit 1
fi

mkdir -p api-docs

npx @redocly/cli build-docs apps/did_service_endpoint/docs/swagger.json \
  --title "DID ServiceEndpoint API" \
  -o api-docs/api-docs.html
