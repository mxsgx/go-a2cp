# Changelog

All notable changes to this project will be documented in this file.

## v0.1.0

### Added
- Initial Apache2 .conf parser.
- AST types for directives, blocks, documents, and positions.
- AST manipulation helpers.
- Builder-style API for creating configs from scratch.
- Rendering and file save support.
- Runnable examples under `examples/`.
- Test fixture layout under `testdata/`.
- GitHub Actions CI workflow.
- GitHub Actions release workflow.
- Versioning policy documentation in `VERSIONING.md`.
- Contributor, security, code-of-conduct, and AI guidance docs.

## Unreleased

### Added
- Functional parse option `WithIncludeResolution(basePath string)`.
- Recursive resolution for Apache2 `Include` and `IncludeOptional` directives.
- Glob-based include matching using `filepath.Glob`.
- Circular include detection with path tracking.
- Parser fixtures for include resolution and circular include cases.
- Runnable examples for include resolution and `IncludeOptional` skip behavior.

### Changed
- `ParseReader` now accepts parse options: `ParseReader(r io.Reader, opts ...ParseOption)`.
- `ParseFile` now accepts parse options: `ParseFile(path string, opts ...ParseOption)`.
- `ParseFile` include resolution now consistently honors `basePath` when provided, including nested includes.
