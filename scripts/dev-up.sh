#!/usr/bin/env bash
set -euo pipefail

# Simple dev launcher for local PoC services.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOG_DIR="${ROOT_DIR}/.devlogs"
ENV_FILE="${ROOT_DIR}/.env"

mkdir -p "${LOG_DIR}"

if [[ -f "${ENV_FILE}" ]]; then
  # shellcheck disable=SC1090
  source "${ENV_FILE}"
fi

require_env() {
  local name="$1"
  if [[ -z "${!name:-}" ]]; then
    echo "[dev-up] missing required env: ${name}" >&2
    exit 1
  fi
}

# eth driver needs this when using the eth method.
require_env "ETH_PRIVATE_KEY_HEX"

check_port() {
  local port="$1"
  if command -v lsof >/dev/null 2>&1; then
    if lsof -iTCP:"${port}" -sTCP:LISTEN >/dev/null 2>&1; then
      if [[ "${DEV_UP_STRICT_PORTS:-}" == "1" ]]; then
        echo "[dev-up] error: port ${port} already in use" >&2
        exit 1
      fi
      echo "[dev-up] warning: port ${port} already in use" >&2
    fi
  fi
}

for port in 50051 50052 50053 50054 50055; do
  check_port "${port}"
done

stop_if_running() {
  local pid_file="$1"
  if [[ -f "${pid_file}" ]]; then
    local pid
    pid="$(cat "${pid_file}")"
    if kill -0 "${pid}" 2>/dev/null; then
      echo "[dev-up] stopping pid ${pid} (${pid_file})"
      kill "${pid}" || true
      sleep 1
    fi
    rm -f "${pid_file}"
  fi
}

start_service() {
  local name="$1"
  local cmd="$2"
  local pid_file="${LOG_DIR}/${name}.pid"
  local log_file="${LOG_DIR}/${name}.log"

  stop_if_running "${pid_file}"
  echo "[dev-up] starting ${name}"
  (cd "${ROOT_DIR}" && nohup bash -c "${cmd}" >"${log_file}" 2>&1 & echo $! > "${pid_file}")
}

# Start servers
start_service "did-registry" "go run ./apps/did-registry/main.go"
start_service "did-registrar" "go run ./apps/did-registrar/main.go"
start_service "demo-issuer" "go run ./apps/demo-issuer/main.go"
start_service "demo-rp" "go run ./apps/demo-rp/main.go"

# Give servers a moment to bind their ports before the client exercise.
sleep 2

echo "[dev-up] running demo-client to exercise the flow"
go run "${ROOT_DIR}/apps/demo-client/main.go" | tee "${LOG_DIR}/demo-client.log"

echo "[dev-up] all services running; logs under ${LOG_DIR}"
