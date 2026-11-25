# Integration tests

Integration tests live in `test/integration` and use [Testcontainers](https://golang.testcontainers.org/)
to start a MySQL container and the application container built from `deploy/Dockerfile`
(the deployment image). `test/setup` holds the container setup.

```bash
go test -count=1 ./test/integration/          # run the tests
make test-integration-coverage                # run them and report coverage
```

Docker must be running.

## Coverage

### Why `go test -cover` is not enough

Acceptance tests send HTTP requests to the application running in its own
container, so the handlers, services, middleware and `cmd/service` never run inside
the `go test` process. `go test -coverpkg=...` only counts packages linked into the
test binary, so those packages were silently left out of the total, and the number
it printed only reflected the repository tests that call the database directly.

### Options considered

| | A. Run the server in-process | B. Instrument the container binary (`-cover` + `GOCOVERDIR`) — **chosen** |
|---|---|---|
| How | Extract the wiring in `cmd/service/main.go` into an importable package, start it from `TestMain`, measure with `go test -coverpkg` | Build the binary in `deploy/Dockerfile` with `go build -cover`, set `GOCOVERDIR`, copy the coverage files out of the container after a graceful stop, convert with `go tool covdata` |
| What is exercised | Go code only; the image, `CMD`, `/wait`, goose migrations and env loading are bypassed | The deployment image, started exactly as in production |
| Production code changes | Refactor of `main` into a reusable package | None in Go code; one optional build argument in `deploy/Dockerfile` |
| Test isolation | Tests and app share one process (global logger, gin mode, env) | Separate processes, like production |
| CI complexity | Low: one `-coverprofile` | Medium: collect files and convert them |
| Remote Docker (GitLab DinD) | Not affected | Files are copied through the Docker API, so no bind mount is needed |

B was chosen because it measures the artifact we ship and keeps the integration tests
running against the deployment image, without restructuring `main` for testability.

### How it works

1. `INTEGRATION_COVERDIR=<host dir>` turns coverage on. Without it the image is built
   exactly like a release image.
2. `test/setup` passes `GO_BUILD_FLAGS="-cover -coverpkg=./internal/...,./cmd/..."` to
   `deploy/Dockerfile` and sets `GOCOVERDIR=/coverage` in the container.
3. On shutdown the container is stopped with `SIGTERM` (30s grace period). The binary
   runs as PID 1, handles the signal, shuts down gracefully and writes its counters on
   exit. A non-zero exit code fails the run, because coverage would be incomplete.
4. `/coverage` is copied from the stopped container into `INTEGRATION_COVERDIR`.
5. `go tool covdata` turns the raw files into a report:

```bash
go tool covdata percent -i=coverage/integration                              # per package
go tool covdata textfmt -i=coverage/integration -o coverage/integration.out  # profile
go tool cover -html=coverage/integration.out -o coverage/integration.html    # HTML
```

### Production behaviour

`GO_BUILD_FLAGS` defaults to empty, so release images build the same binary as before.
Coverage instrumentation only exists in the image the tests build.

### CI

- **GitHub Actions** (`Unit Tests` job): prints per-package and total coverage to the job
  summary and uploads `coverage/` (raw data, profile, HTML) as the `integration-coverage` artifact.
- **GitLab CI** (`test` job): prints the same summary, exposes the total through the
  `coverage:` regex (shown on merge requests) and keeps `coverage/` as a job artifact.
