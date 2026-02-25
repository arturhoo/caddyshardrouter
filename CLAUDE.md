# CLAUDE.md

## Project Overview

Caddy Shard Router — a set of Caddy middleware modules that route reverse-proxy
requests to backend shards based on a customer identifier extracted from either a
JWT Bearer token or the JSON request body. Shard mappings are stored in Redis.

## Build & Run

```bash
# Build the Caddy binary
CGO_ENABLED=0 go build -cover -o caddy cmd/main.go

# Start the full stack (Caddy + upstreams + Redis)
docker compose up -d --build

# Run unit tests
go test ./...

# Run integration tests (requires running containers)
docker run --net=host -v "$(pwd):/app" ghcr.io/orange-opensource/hurl:latest --test /app/test.hurl

# Lint
golangci-lint run --timeout=5m
```

## Clean Code Rules

Follow Uncle Bob's Clean Code principles. Every change must leave the codebase
cleaner than you found it (the Boy Scout Rule).

### Naming

- Names must reveal intent. A reader should understand purpose from the name
  alone — no comments needed to explain what a variable holds or a function does.
- Use pronounceable, searchable names. Avoid single-letter variables except for
  short-lived loop indices (`i`, `j`).
- Functions are named as verbs or verb phrases (`ParseJWT`, `ServeHTTP`). Types
  and structs are named as nouns (`JWTShardRouter`, `BodyShardRouter`).
- Don't encode type information into names. No Hungarian notation, no prefixes
  like `I` for interfaces.
- Pick one word per concept and stick with it across the codebase. Don't mix
  `get`, `fetch`, and `retrieve` for the same kind of operation.

### Functions

- Functions do one thing. If you can extract another function from it with a
  name that is not merely a restatement of its implementation, it does too much.
- Keep functions short. A function that scrolls beyond a single screen is a
  candidate for extraction.
- Limit arguments. Zero is ideal, one or two is fine, three requires strong
  justification. When you need more, group related arguments into a struct.
- No side effects. A function named `ParseJWT` must not secretly modify global
  state. If it must do more, rename it to say so.
- Prefer returning errors over panic. Follow Go's `(result, error)` convention.
- Extract `try/catch` bodies into their own functions — in Go this means error
  handling blocks should be concise and not bury the happy path.

### Comments

- Don't comment bad code — rewrite it. A comment is a failure to express intent
  through code.
- The only good comments are: legal headers, explanations of intent when no
  better alternative exists, clarifications of third-party APIs, warnings of
  consequences, and TODO markers for work you intend to finish.
- Never leave commented-out code. Delete it. Version control remembers.
- Don't add doc comments to unexported helpers if the name already communicates
  the purpose.

### Formatting

- Run `gofmt` (or `goimports`) on all Go code. No exceptions.
- `golangci-lint` must pass cleanly. The CI runs it on every pull request.
- Keep files focused. Each file should contain closely related functionality.
  The current layout (one file per middleware, a separate file for JWT utilities,
  a separate file for Redis initialization) is the standard to follow.
- Vertical ordering: callers above callees. High-level functions appear near the
  top of a file, lower-level helpers below.

### Error Handling

- Errors are not exceptional — they are part of the contract. Always check the
  returned `error` and handle it or propagate it.
- Return errors with enough context for the caller to act
  (`fmt.Errorf("failed to query redis")` not `fmt.Errorf("error")`).
- Don't silently swallow errors. If you must ignore one, assign it to `_` and
  leave a brief reason why.

### Tests

- Tests are first-class code. They follow the same clean code standards as
  production code: clear names, no duplication, readable assertions.
- Unit tests live alongside the code they test (`jwtutils_test.go` next to
  `jwtutils.go`). Integration tests live in `test.hurl`.
- Every new function or bug fix gets a test. Minimum 80% code coverage is
  enforced in CI — treat it as a floor, not a ceiling.
- Test names should describe the scenario: `TestParseJWT_InvalidToken` not
  `TestParse2`.
- Keep tests fast. No network calls in unit tests. Integration tests that need
  Redis or upstream services belong in the Hurl suite.

### Dependencies & Structure

- No dead code. If a function, type, or import is unused, delete it.
- Minimize coupling. Modules depend on Caddy interfaces and Redis — keep it that
  way. Don't introduce new external dependencies without a strong reason.
- One level of abstraction per function. Don't mix HTTP response writing with
  business logic in the same block.

### The Boy Scout Rule

Leave every file you touch cleaner than you found it. That does not mean
rewriting unrelated code in the same PR — it means fixing a misleading name, an
unnecessary comment, or a missing error check when you encounter one during your
work.
