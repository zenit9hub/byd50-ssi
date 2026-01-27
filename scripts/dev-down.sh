#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOG_DIR="${ROOT_DIR}/.devlogs"
DEFAULT_PORTS="50051 50052 50053 50054 50055 8080"
PORTS="${DEV_PORTS:-$DEFAULT_PORTS}"
SUDO="${DEV_DOWN_SUDO:-}"
PATTERN="${DEV_DOWN_PATTERN:-apps/(did-|demo-|did_service_endpoint|geth_client)|byd50-ssi}"

run_cmd() {
  if [[ -n "${SUDO}" ]]; then
    sudo "$@"
  else
    "$@"
  fi
}

kill_pid() {
  local pid="$1"
  if [[ -z "${pid}" ]]; then
    return
  fi
  if kill -0 "${pid}" 2>/dev/null; then
    echo "[dev-down] stopping pid ${pid}"
    run_cmd kill "${pid}" || true
    sleep 1
    if kill -0 "${pid}" 2>/dev/null; then
      echo "[dev-down] force killing pid ${pid}"
      run_cmd kill -9 "${pid}" || true
    fi
  fi
}

if [[ ! -d "${LOG_DIR}" ]]; then
  echo "[dev-down] no .devlogs directory found"
  exit 0
fi

for pid_file in "${LOG_DIR}"/*.pid; do
  if [[ ! -f "${pid_file}" ]]; then
    continue
  fi
  pid="$(cat "${pid_file}")"
  echo "[dev-down] stopping pid ${pid} (${pid_file})"
  kill_pid "${pid}"
  rm -f "${pid_file}"
done

list_pids_by_port() {
  local port="$1"
  run_cmd lsof -tiTCP:LISTEN -i :"${port}" 2>/dev/null || true
}

if command -v lsof >/dev/null 2>&1; then
  for port in ${PORTS}; do
    pids="$(list_pids_by_port "${port}")"
    if [[ -n "${pids}" ]]; then
      echo "[dev-down] stopping port ${port} (pids: ${pids})"
      for pid in ${pids}; do
        kill_pid "${pid}"
      done
    fi
  done

  # Fallback: kill known app processes by pattern
  if command -v pgrep >/dev/null 2>&1; then
    pids="$(run_cmd pgrep -f "${PATTERN}" 2>/dev/null || true)"
    if [[ -n "${pids}" ]]; then
      echo "[dev-down] stopping matching processes (pattern: ${PATTERN})"
      for pid in ${pids}; do
        kill_pid "${pid}"
      done
    fi
  else
    echo "[dev-down] pgrep not found; skipping pattern-based cleanup"
  fi

  # Recheck ports to ensure they are free
  for port in ${PORTS}; do
    pids="$(list_pids_by_port "${port}")"
    if [[ -n "${pids}" ]]; then
      echo "[dev-down] warning: port ${port} still in use (pids: ${pids})"
      echo "[dev-down] hint: try 'DEV_DOWN_SUDO=1 make dev-down' if owned by another user"
    fi
  done
else
  echo "[dev-down] lsof not found; skipping port-based cleanup"
fi

echo "[dev-down] done"
