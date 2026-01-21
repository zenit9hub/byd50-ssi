#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT_DIR}"

if ! command -v node >/dev/null 2>&1; then
  echo "[swagger-pdf] node not found. Install Node.js first." >&2
  exit 1
fi

node scripts/html-to-pdf.js
