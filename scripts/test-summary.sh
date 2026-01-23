#!/usr/bin/env bash
set -euo pipefail

tmp_log="$(mktemp)"
go test -v ./... | tee "${tmp_log}"
awk '
/^=== RUN /{run++}
/^--- PASS: /{pass++}
/^--- FAIL: /{fail++}
END{printf "tests: run=%d pass=%d fail=%d\n", run, pass, fail}' "${tmp_log}"
rm -f "${tmp_log}"
