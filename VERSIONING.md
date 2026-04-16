# Versioning Policy

This project follows Semantic Versioning (`MAJOR.MINOR.PATCH`).

## Rules

- `PATCH`: bug fixes and internal improvements with no public API break.
- `MINOR`: backward-compatible new features and API additions.
- `MAJOR`: incompatible public API changes.

Before `v1.0.0`, minor API changes may still occur in minor releases.

## Release Tags

Use annotated tags in this format:

- `v0.1.0`
- `v0.2.3`
- `v1.0.0`

Only tags matching `v*.*.*` trigger the release workflow.

## Release Process

1. Update `CHANGELOG.md` under a new version heading.
2. Commit all release-ready changes.
3. Create and push an annotated tag:

   git tag -a v0.1.0 -m "v0.1.0"
   git push origin v0.1.0

4. GitHub Actions publishes the release from the pushed tag.

## Branch Policy

- Open pull requests into `main`.
- CI must pass before merge.
