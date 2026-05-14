#!/usr/bin/env bash
set -euo pipefail

PORT="${PORT:-18080}"
BASE_URL="http://localhost:${PORT}"
SERVER_LOG="$(mktemp)"

cleanup() {
  if [[ -n "${SERVER_PID:-}" ]]; then
    kill "${SERVER_PID}" >/dev/null 2>&1 || true
    wait "${SERVER_PID}" >/dev/null 2>&1 || true
  fi
  rm -f "${SERVER_LOG}"
}
trap cleanup EXIT

echo "Running Go tests..."
go test ./...

echo "Starting URL Lookup Service on port ${PORT}..."
PORT="${PORT}" go run ./cmd/urlinfo >"${SERVER_LOG}" 2>&1 &
SERVER_PID="$!"

ready=false
for _ in {1..30}; do
  if curl -fsS "${BASE_URL}/readyz" >/dev/null 2>&1; then
    ready=true
    break
  fi
  sleep 0.2
done

if [[ "${ready}" != "true" ]]; then
  echo "Service did not become ready."
  echo
  cat "${SERVER_LOG}"
  exit 1
fi

passed=0
failed=0

check_verdict() {
  local label="$1"
  local path="$2"
  local expected="$3"
  local body

  body="$(curl -fsS "${BASE_URL}${path}")"

  if [[ "${body}" == *"\"verdict\":\"${expected}\""* ]]; then
    printf "PASS %-10s %s\n" "${expected}" "${label}"
    passed=$((passed + 1))
    return
  fi

  printf "FAIL %-10s %s\n" "${expected}" "${label}"
  echo "  response: ${body}"
  failed=$((failed + 1))
}

check_probe() {
  local label="$1"
  local path="$2"

  if curl -fsS "${BASE_URL}${path}" >/dev/null 2>&1; then
    printf "PASS %-10s %s\n" "probe" "${label}"
    passed=$((passed + 1))
    return
  fi

  printf "FAIL %-10s %s\n" "probe" "${label}"
  failed=$((failed + 1))
}

echo
echo "Running smoke checks..."

check_probe "healthz" "/healthz"
check_probe "readyz"  "/readyz"

check_verdict "malware.test/bad" "/urlinfo/1/malware.test/bad" "malicious"
check_verdict "bad.example:443/download" "/urlinfo/1/bad.example:443/download" "malicious"
check_verdict "google.com/search?q=test" "/urlinfo/1/google.com/search?q=test" "safe"
check_verdict "example.com/docs" "/urlinfo/1/example.com/docs" "safe"

echo
echo "Summary: ${passed} passed, ${failed} failed"

if [[ "${failed}" -gt 0 ]]; then
  exit 1
fi
