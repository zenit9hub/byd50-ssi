#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${ROOT_DIR}/.env"

if [[ -f "${ENV_FILE}" ]]; then
  # shellcheck disable=SC1090
  source "${ENV_FILE}"
fi

if [[ -z "${ETH_PRIVATE_KEY_HEX:-}" ]]; then
  echo "[integration-test] missing ETH_PRIVATE_KEY_HEX in .env" >&2
  exit 1
fi

cleanup() {
  "${ROOT_DIR}/scripts/dev-down.sh" || true
}
trap cleanup EXIT

echo "[integration-test] starting services"
make -C "${ROOT_DIR}" dev-up

echo "[integration-test] running integration tests"
go test -v -tags=integration ./test/...

echo "[integration-test] done"
