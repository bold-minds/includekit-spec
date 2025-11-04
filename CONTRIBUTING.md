# Contributing to includekit-spec

Thanks for helping improve IncludeKit’s spec. This repo is the **source of truth** for SDKs and the core engine, so precision and determinism are critical.

## Code of Conduct

We follow the [Contributor Covenant v2.1](https://www.contributor-covenant.org/version/2/1/code_of_conduct/).

## Design principles

1. **Single schema**: `/schema/v*.json` generates all language artifacts.
2. **Determinism**: JSON Canonicalization (JCS) + SHA-256 → `shapeId`.
3. **Safety over cleverness**: Unknown operators must be expressible as `custom:*`; downstream engines will invalidate conservatively.
4. **SemVer discipline**: The schema’s version drives the packages’ major versions.

## Development workflow

1. **Discuss changes** in an issue (motivation, ORM mappings, examples).
2. **Update schema** (`schema/v0-1-0.json`) and documentation.
3. **Run the test suite** (auto-regenerates code):

   ```bash
   ./scripts/test.sh
   ```

   This automatically:
   - Builds the Go codegen tool
   - Regenerates TypeScript types
   - Builds testkits
   - Runs all tests (TS + Go)
   - Verifies conformance vectors
   - Checks no-runtime constraints

4. **Update test vectors** if canonical output or shapeId expectations change:

   ```bash
   go run tools/tests/generate-vectors.go
   ```

   Vectors are in `tools/tests/vectors/query-shapes.json`.

### Open a PR with:

- [ ] Schema diff
- [ ] README.md changes
- [ ] New/updated vectors
- [ ] Passing CI (unit + conformance)

### Style & tools

- [ ] Conventional Commits (feat:, fix:, docs:, refactor:…)
- [ ] TS: strict mode, ESLint, Prettier
- [ ] Go: gofmt, golangci-lint
- [ ] Tests: ≥85% coverage per package, conformance must pass

### Adding an operator
- [ ] Propose in an issue (with cross-ORM references).
- [ ] Add to schema enum/pattern in `schema/v0-1-0.json` (^`custom:.*$` allowed).
- [ ] Document in `schema/README.md`.
- [ ] Add test cases to `tools/tests/generate-vectors.go`.
- [ ] Run `go run tools/tests/generate-vectors.go` to update vectors.
- [ ] Run `./scripts/test.sh` to verify.

### Releasing

Version management uses the `VERSION` file as single source of truth:

1. **Update version:**
   ```bash
   echo "0.2.0" > VERSION
   ```

2. **Sync all files:**
   ```bash
   go run tools/version/sync.go
   ```
   This updates schema filename, package.json files, workflows, etc.

3. **Regenerate code:**
   ```bash
   cd codegen && go run .
   ```

4. **Test:**
   ```bash
   ./scripts/test.sh
   ```

5. **Update CHANGELOG.md** with release notes.

6. **Commit and tag:**
   ```bash
   git add -A
   git commit -m "chore: bump version to v0.2.0"
   git tag v0.2.0
   git push origin main --tags
   ```

7. **CI will automatically:**
   - Run full test suite
   - Publish to npm
   - Create GitHub release

### Security
- [ ] No secrets in tree.
- [ ] Canonicalizers must not rely on undefined map iteration order.
- [ ] This repo does not process live customer data.

### Notes

- To get started, fork, clone, and test:

```bash
git clone https://github.com/YOUR_USERNAME/ik-spec.git
cd ik-spec
./scripts/test.sh  # Installs deps, builds, tests everything
```

- **Do not import testkit packages from production code.** CI will fail if:
  - `@includekit/types-testkit` is imported outside test files (TypeScript)
  - `github.com/bold-minds/ik-spec/go/tests` is imported from `pkgs/go/types` (Go)
- All validators, canonicalization, and shapeId utilities live in testkit packages only.
- Production packages (`@includekit/types`, `ik-spec/go`) contain types only.
