# Agent Guide

This repository contains `go-a2cp`, a Go library for parsing and manipulating Apache2 `.conf` files.

## What This Project Does

- Parse Apache2 config into an AST.
- Manipulate directives and nested blocks in memory.
- Render config back to text.
- Build configs from scratch without parsing an existing file.

## Key Files

- `ast.go`: AST types (`Document`, `Directive`, `Block`, `Position`).
- `parser.go`: parsing logic.
- `render.go`: string rendering.
- `manipulate.go`: mutation helpers and builders.
- `io.go`: file/stream helpers.
- `examples/`: runnable examples.
- `testdata/`: fixture `.conf` files for tests.

## Working Rules

- Prefer minimal, focused changes.
- Do not revert user changes you did not make.
- Use `apply_patch` for edits.
- Keep public API changes deliberate and documented.
- Add or update tests for behavior changes.
- Preserve existing style and ASCII unless the file already uses otherwise.

## Common Commands

- `go test ./...` to run the test suite.
- `go test -v` for verbose test output.
- `go run ./examples/parse-string`
- `go run ./examples/parse-file`
- `go run ./examples/manipulate-save`
- `go run ./examples/from-scratch`
- `go run ./examples/include-resolution`
- `go run ./examples/include-optional-skip`

## Test Fixture Conventions

- Put parser cases under `testdata/parser/`.
- Put round-trip cases under `testdata/roundtrip/`.
- Put larger manual examples under `testdata/examples/`.
- Load fixtures with `ParseFile("testdata/.../*.conf")`.

## Examples Conventions

- Add standalone runnable programs under `examples/<name>/main.go`.
- Keep example output files out of version control if they are generated.
- Prefer examples that show one clear use case.

## Public API Guidance

If you add or change exported types or methods:

- Update `README.md`.
- Update `CHANGELOG.md` if the change is user-visible.
- Update `VERSIONING.md` if the release process or policy changes.

## Release and Versioning

- Follow `VERSIONING.md` for SemVer rules.
- Release tags use the form `vMAJOR.MINOR.PATCH`.
- Tag pushes trigger GitHub Actions release automation.

## Safety

- Avoid destructive git commands unless explicitly requested.
- Do not delete or overwrite generated user files unless asked.
- If the task is unclear, inspect the repo first before editing.
