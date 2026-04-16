# Contributing

Thank you for considering contributing to go-a2cp.

## Development Setup

1. Install Go 1.26 or newer.
2. Clone the repository.
3. Run tests:

   go test ./...

## Project Scope

This project focuses on Apache2 .conf parsing and manipulation with a small, stable API.

When proposing changes, prefer:
- Backward-compatible API additions.
- Clear tests for new parser or renderer behavior.
- Small, focused pull requests.

## Pull Request Checklist

- Add or update tests for behavior changes.
- Update README.md when public API changes.
- Keep formatting consistent with existing code.
- Ensure all tests pass locally.

## Versioning and Releases

This project uses Semantic Versioning. See `VERSIONING.md` for full policy.

- Release tags must follow `vMAJOR.MINOR.PATCH`.
- Pushing a matching tag triggers the release workflow.

## Reporting Bugs

Please include:
- Minimal .conf input that reproduces the issue.
- Expected behavior vs actual behavior.
- Go version and OS.

## Discussion

Open an issue first for large changes so design can be discussed before implementation.
