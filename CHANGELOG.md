# Changelog

All notable changes to this project will be documented in this file.

## Unreleased

### Added
- New AST node `Comment` for preserved Apache2 comments.
- Runnable example `examples/comment-roundtrip` for comment-preserving round-trip rendering.
- Runnable example `examples/backslash-comments` for continuation lines with trailing `\\` and comment preservation.
- Helper APIs `(*Document).AddComment(text string, opts ...CommentOption) error` and `(*Block).AddComment(text string, opts ...CommentOption) error`.
- Inline comment option `WithInlineComment()`.
- Raw text option `WithRawCommentText()` for verbatim comment spacing.
- Trailing comments on closing tags are now preserved inline via `Block.EndComment`.
- Parser fixture `testdata/parser/include-inline-comment/` for include-resolution inline comment behavior.

### Changed
- Parser now preserves full-line and inline comments as `Comment` statements.
- Parser now preserves comments that appear on physical lines consumed by line continuation (`\\`).
- Parser now keeps quote state across continued physical lines so `#` after a closed quote is parsed as a comment.
- Parser no longer rescans the full accumulated continuation buffer to track quote state (incremental state tracking).
- Continuation line joining now normalizes whitespace to avoid doubled spaces in rendered output.
- Include resolution now consumes inline comments attached to `Include`/`IncludeOptional` directives so they stay near the include expansion.
- Renderer now writes `Comment` statements back and keeps inline comments on the same line.
- Renderer now keeps closing-tag comments on the same line as `</...>`.
- Renderer now escapes `\r` and `\n` inside comment text to prevent multi-line comment injection in generated config.
- `AddComment` now normalizes non-empty text to render with `# ` by default.
- `AddInlineComment` remains available as a compatibility alias for `AddComment(text, WithInlineComment())`.

## v0.2.0 - 2026-04-17

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

## v0.1.0 - 2026-04-16

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
