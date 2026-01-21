#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOG_DIR="${ROOT_DIR}/.devlogs"

if [[ ! -d "${LOG_DIR}" ]]; then
  echo "[dev-down] no .devlogs directory found"
  exit 0
fi

for pid_file in "${LOG_DIR}"/*.pid; do
  if [[ ! -f "${pid_file}" ]]; then
    continue
  fi
  pid="$(cat "${pid_file}")"
  if kill -0 "${pid}" 2>/dev/null; then
    echo "[dev-down] stopping pid ${pid} (${pid_file})"
    kill "${pid}" || true
    sleep 1
  fi
  rm -f "${pid_file}"
done

echo "[dev-down] done"
