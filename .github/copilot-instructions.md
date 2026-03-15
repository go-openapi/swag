# Copilot Instructions

## Project Overview

Go mono-repo of utility libraries for the `go-openapi` / `go-swagger` ecosystem.
The root package is a backward-compatible facade — all its exported functions are
deprecated and delegate to sub-packages. New code belongs in the appropriate
sub-package, not the root.

The repo exists to give the go-swagger code generator and runtime a single,
well-tested set of helpers for name mangling, type conversions, JSON/YAML
handling, file loading, and other common tasks, without pulling in heavy
external dependencies.

## Workspace & Modules

Go workspace (`go.work`, Go 1.24) with 15+ modules including:
`conv`, `cmdutils`, `fileutils`, `jsonname`, `jsonutils`, `loading`,
`mangling`, `netutils`, `stringutils`, `typeutils`, `yamlutils`, and others.

## Conventions

Coding conventions are found beneath `.github/copilot`

### Summary

- All `.go` files must have SPDX license headers (Apache-2.0).
- Commits require DCO sign-off (`git commit -s`).
- Linting: `golangci-lint run` — config in `.golangci.yml` (posture: `default: all` with explicit disables).
- Every `//nolint` directive **must** have an inline comment explaining why.
- Tests: `go test work ./...` (mono-repo). CI runs on `{ubuntu, macos, windows} x {stable, oldstable}` with `-race`.
- Test framework: `github.com/go-openapi/testify/v2` (not `stretchr/testify`; `testifylint` does not work).

See `.github/copilot/` (symlinked to `.claude/rules/`) for detailed rules on Go conventions, linting, testing, and contributions.
