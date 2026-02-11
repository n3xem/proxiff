# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is Proxiff?

A Diffy-like HTTP proxy written in Go. Forwards requests to two servers (newer/current), compares responses, logs differences, and returns the current server's response to clients.

## Build & Development Commands

```bash
make build             # Build proxiff binary to ./proxiff
make build-examples    # Build example servers
make test              # Run tests with race detector and coverage
make hadolint          # Lint Dockerfile
make clean             # Remove built binaries
```

Run a single test:
```bash
go test ./proxy/ -run TestProxyForwardsToNewerAndCurrent -v
```

## CI Checks

CI (`test.yml`) enforces: `go vet`, `gofmt -s` formatting, cyclomatic complexity ≤ 10 (via gocyclo), race detection, and license allowlist (MIT, Apache-2.0, BSD-2/3-Clause, ISC, MPL-2.0).

## Architecture

**Request flow:** Client → `Proxy` → forwards to both newer & current servers → `Comparator` compares responses → logs diffs → returns current's response.

**Key packages:**
- `cmd/proxiff/cmd/` — CLI (Cobra). `start.go` wires up the proxy; `root.go`/`version.go` are boilerplate.
- `proxy/` — Core `Proxy` struct. Handles HTTP forwarding and comparison orchestration.
- `comparator/` — `Comparator` interface, `Response`/`Result` types, and `SimpleComparator` implementation (uses `google/go-cmp` for detailed diffs).

## Running Locally

```bash
make build && make build-examples
# Terminal 1: current server
./sample-server --port 8081
# Terminal 2: newer server (with differences)
./sample-server --port 8082 --newer
# Terminal 3: proxy
./proxiff start --newer http://localhost:8082 --current http://localhost:8081
# Terminal 4: test
curl http://localhost:8080/api/users
```

Or use `./example/run-demo.sh` which handles all of the above.
