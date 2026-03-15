# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go monorepo of utility libraries for the `go-openapi` / `go-swagger` ecosystem. The root package
is a **backward-compat facade** — all its exported functions are deprecated and delegate to sub-packages.

See [docs/MAINTAINERS.md](../docs/MAINTAINERS.md) for CI/CD, release process, and repo structure details.

## Workspace & Modules

Go workspace (`go.work`, Go 1.24). Contains 15+ modules:

| Module | Purpose |
|--------|---------|
| `.` (root) | Deprecated shims forwarding to sub-packages |
| `cmdutils` | CLI utility helpers |
| `conv` | Type conversions: string-to-value, value-to-pointer, pointer-to-value (generics) |
| `fileutils` | File and path utilities |
| `jsonname` | Infer JSON field names from Go struct tags |
| `jsonutils` | JSON read/write, `ConcatJSON`, `JSONMapSlice`; pluggable adapter system |
| `jsonutils/adapters/easyjson` | easyjson adapter (separate module to isolate dependency) |
| `loading` | Load specs from filesystem or HTTP |
| `mangling` | Name mangling: `ToGoName`, `ToFileName`, `ToVarName`, etc.; configurable initialisms |
| `netutils` | Host/port parsing |
| `stringutils` | Slice search, query parameter splitting |
| `typeutils` | Zero-value and nil-safe interface checks |
| `yamlutils` | YAML-to-JSON conversion, ordered YAML documents |

Inter-module dependencies use `replace` directives pointing to local paths.

## Testing

```sh
# Run all tests across every workspace module
go test work ./...

# Run tests for a single module
go test ./conv/...
```

Note: plain `go test ./...` only tests the root module. The `work` pattern expands to all
modules listed in `go.work`.

CI runs tests on `{ubuntu, macos, windows} x {stable, oldstable}` with `-race` via `gotestsum`.

### Fuzz tests

```sh
# List all fuzz targets across the workspace
go test work -list Fuzz ./...

# Run a specific fuzz target (go test -fuzz cannot span multiple packages)
go test -fuzz=Fuzz -run='FuzzToGoName$' -fuzztime=1m30s ./mangling
```

Fuzz corpus lives in `testdata/fuzz/` within each package. CI runs each fuzz target for 1m30s
with a 5m minimize timeout.

### Test framework

`github.com/go-openapi/testify/v2` — a zero-dep fork of `stretchr/testify`.
Because it's a fork, `testifylint` does not work.

Patterns: table-driven tests with `t.Run`, fuzz tests, benchmarks, and integration
tests in `jsonutils/adapters/testintegration/`.

## Linting

```sh
golangci-lint run
```

Config: `.golangci.yml` — posture is `default: all` with explicit disables.
See [docs/STYLE.md](../docs/STYLE.md) for the rationale behind each disabled linter.

Key rules:
- Every `//nolint` directive **must** have an inline comment explaining why.
- Prefer disabling a linter over scattering `//nolint` across the codebase.

## Code Conventions

- All files must have SPDX license headers (Apache-2.0).
- Go version policy: support the 2 latest stable Go minor versions.
- Commits require DCO sign-off (`git commit -s`).
- Root-level `*_iface.go` files are thin deprecated wrappers — do not add new public API there.
- `jsonutils` uses a pluggable adapter registry (`adapters.Registry`); new JSON backends
  should be separate modules implementing interfaces from `jsonutils/adapters/ifaces`.
- `mangling.NameMangler` is configured via functional options and is concurrency-safe.
