# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2024-11-04

### Added
- Initial pre-release of IncludeKit Universal Format
- JSON Schema for QueryShape, MutationEvent, Dependencies
- TypeScript spec package (`@includekit/spec`)
- TypeScript testkit package (`@includekit/spec-testkit`)
- Go spec package (`github.com/bold-minds/includekit-spec/go`)
- Go testkit package (`github.com/bold-minds/includekit-spec/go/tests`)
- Cross-language conformance test vectors (15 comprehensive tests)
- Automated version management via `VERSION` file and sync tool
- Go-based code generator (`codegen/`)
- CI/CD workflows for testing and release
- GitHub Actions for version sync verification

### Schema Features
- QueryShape with filters, relations, pagination, ordering, grouping
- FilterSpec with boolean logic (AND, OR, NOT)
- Comparison operators: eq, ne, gt, gte, lt, lte, in, notIn, contains, startsWith, endsWith
- MutationEvent for write tracking
- Deterministic JSON canonicalization (JCS RFC 8785)
- SHA-256 shape ID computation

[0.1.0]: https://github.com/bold-minds/includekit-spec/releases/tag/v0.1.0
