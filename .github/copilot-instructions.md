# Copilot Instructions

This repository contains `go-a2cp`, a Go library for parsing and manipulating Apache2 `.conf` files.

## Goals

- Keep changes small and focused.
- Preserve public API stability unless the task explicitly requires a breaking change.
- Add or update tests for behavior changes.
- Update documentation when public API or workflow changes.

## Codebase Map

- `ast.go`: AST types (`Document`, `Directive`, `Block`, `Position`).
- `parser.go`: parsing logic.
- `render.go`: rendering logic.
- `manipulate.go`: mutation helpers and builders.
- `io.go`: file helpers.
- `examples/`: runnable example programs.
- `testdata/`: parser and round-trip fixtures.

## Practical Rules

- Prefer minimal, focused edits.
- Avoid destructive git operations.
- Keep ASCII unless the existing file uses Unicode.
- Use existing naming and formatting conventions.
- Do not add inline comments unless they improve clarity.

## Validation

- Run `go test ./...` after code changes.
- Run `go run ./examples/...` when example programs are affected.

## Useful Paths

- Repository guide: `AGENTS.md`
- Versioning policy: `VERSIONING.md`
- Changelog: `CHANGELOG.md`
- Contributing guide: `CONTRIBUTING.md`
