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
DEFAULT_PORTS="50051 50052 50053 50054 50055 8080"
PORTS="${DEV_PORTS:-$DEFAULT_PORTS}"
SUDO="${DEV_STATUS_SUDO:-}"

run_cmd() {
  if [[ -n "${SUDO}" ]]; then
    sudo "$@"
  else
    "$@"
  fi
}

if command -v lsof >/dev/null 2>&1; then
  for port in ${PORTS}; do
    pids="$(run_cmd lsof -tiTCP:LISTEN -i :"${port}" 2>/dev/null || true)"
    if [[ -n "${pids}" ]]; then
      for pid in ${pids}; do
        cmd="$(ps -p "${pid}" -o comm= 2>/dev/null || true)"
        echo "  ${port}: in use (pid=${pid} ${cmd})"
      done
    else
      echo "  ${port}: free"
    fi
  done
else
  echo "  lsof not found; skipping port check"
fi
