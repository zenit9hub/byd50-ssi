#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOG_DIR="${ROOT_DIR}/.devlogs"

echo "[dev-status] pid files"
if [[ ! -d "${LOG_DIR}" ]]; then
  echo "  no .devlogs directory"
else
  for pid_file in "${LOG_DIR}"/*.pid; do
    if [[ ! -f "${pid_file}" ]]; then
      continue
    fi
    pid="$(cat "${pid_file}")"
    if kill -0 "${pid}" 2>/dev/null; then
      echo "  ${pid_file##*/}: running (${pid})"
    else
      echo "  ${pid_file##*/}: stale (${pid})"
    fi
  done
fi

echo "[dev-status] listening ports"
if command -v lsof >/dev/null 2>&1; then
  for port in 50051 50052 50053 50054 50055; do
    if lsof -iTCP:"${port}" -sTCP:LISTEN >/dev/null 2>&1; then
      echo "  ${port}: in use"
    else
      echo "  ${port}: free"
    fi
  done
else
  echo "  lsof not found; skipping port check"
fi
