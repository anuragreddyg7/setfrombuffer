# Test runner scripts

This folder contains helper scripts to run the valkey-go test suites locally using Docker Compose.

## run-tests.sh
Usage:

```zsh
# run both unit and integration tests (default)
./tools/run-tests.sh

# run only unit (mock) tests
./tools/run-tests.sh unit

# run only integration tests (requires Docker/Compose)
./tools/run-tests.sh integration
```

What the script does:
- Starts `docker compose up -d` using the repository `docker-compose.yml`.
- Waits for a set of expected ports to accept connections.
- Runs unit tests (`TestSetFromBuffer`) and/or the entire `valkeycompat` test suite.
- Generates coverage profile `cover.out` for `valkeycompat`.
- Tears down the containers with `docker compose down -v --remove-orphans`.

Notes:
- Requires Docker and `nc` (netcat) installed and running.
- The integration test suite expects services on several ports; the included `docker-compose.yml` contains images and ports used by the test harness.
- If you run into any failing tests, collect the failing `go test` output and `docker compose logs --tail=200`, and paste them here so I can help triage.
