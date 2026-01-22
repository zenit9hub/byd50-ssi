#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT_DIR}"

EXCLUDE_REGEX='pkg/did/core/driver|pkg/did/core/driver/scdid|pkg/did/c-shared|pkg/did/service|pkg/did/core/rc'

PKGS=$(go list ./pkg/did/... | grep -vE "${EXCLUDE_REGEX}")
COVERPROFILE="coverage.out"

if [[ -z "${PKGS}" ]]; then
  echo "[coverage-did] no packages selected"
  exit 1
fi

COVERPKG=$(echo "${PKGS}" | paste -sd, -)

go test -v -coverprofile="${COVERPROFILE}" -coverpkg="${COVERPKG}" ${PKGS}

TOTAL=$(go tool cover -func="${COVERPROFILE}" | awk '/^total:/ {gsub(/%/,"",$3); print $3}')
if [[ -z "${TOTAL}" ]]; then
  echo "[coverage-did] failed to read total coverage"
  exit 1
fi

echo "[coverage-did] total coverage: ${TOTAL}%"
awk -v total="${TOTAL}" 'BEGIN { if (total+0 < 95) exit 1 }'
