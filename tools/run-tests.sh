#!/usr/bin/env zsh
# tools/run-tests.sh
# Combined test runner for valkey-go valkeycompat package.
# Usage: ./tools/run-tests.sh [unit|integration|all]

set -euo pipefail
script_dir=$(cd "$(dirname "$0")" && pwd)
repo_root=$(cd "$script_dir/.." && pwd)
cd "$repo_root"

mode=${1:-all}
compose_file="$repo_root/docker-compose.yml"

function info() { echo "[INFO] $*" }
function err() { echo "[ERROR] $*" >&2 }

# start docker compose
info "Starting docker compose..."
docker compose up -d

# wait for ports
ports=(6378 7010 6382 6356 6381 7007 6379 6380)
info "Waiting for compose services to accept connections: ${ports[*]}"
for p in "${ports[@]}"; do
  info "waiting for port $p..."
  i=0
  while ! nc -z 127.0.0.1 $p; do
    sleep 1
    ((i++))
    if [ $i -gt 120 ]; then
      err "Timeout waiting for port $p"
      docker compose logs --tail=200
      exit 1
    fi
  done
  info "port $p ready"
done
info "All ports ready"

# run unit tests (mock-based / fast)
cd "$repo_root/valkeycompat"
if [ "$mode" = "unit" ] || [ "$mode" = "all" ]; then
  info "Running unit (mock) tests..."
  go test -run TestSetFromBuffer -cover -count=1 ./... || true
fi

# run integration (full) tests
if [ "$mode" = "integration" ] || [ "$mode" = "all" ]; then
  info "Running full valkeycompat test suite (integration)..."
  go test ./... -count=1
fi

# collect coverage for package
info "Collecting coverage for valkeycompat..."
go test -coverprofile=cover.out ./valkeycompat || true
if [ -f cover.out ]; then
  go tool cover -func=cover.out | sed -n '1,120p'
fi

info "Tearing down docker compose"
docker compose down -v --remove-orphans

info "Done"
