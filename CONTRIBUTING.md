# Contributing to includekit-spec

Thanks for helping improve IncludeKit’s spec. This repo is the **source of truth** for SDKs and the core engine, so precision and determinism are critical.

## Code of Conduct

We pledge to make participation in our community a harassment-free experience for everyone. Please see [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for full details.

## 🚀 Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

### Pull Requests

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed the schema, regenerate types and update documentation.
4. Ensure the test suite passes.
5. Make sure your code follows the project style.
6. Issue that pull request!

## 🏗️ Development Setup

### Prerequisites

- **Node.js** ≥ 20
- **Go** ≥ 1.22
- **Git**

### Clone and Setup

```bash
git clone https://github.com/bold-minds/includekit-spec.git
cd includekit-spec
./scripts/test.sh
```

This will:
- Build the code generator
- Generate TypeScript and Go types
- Run all tests (TypeScript + Go)
- Verify the no-runtime constraint

## 🧪 Testing

We maintain high test coverage and all contributions should include appropriate tests.

### TypeScript Tests

```bash
cd pkgs/ts/tests
npm install
npm run build
npm test
```

### Go Tests

```bash
cd pkgs/go
go test -v ./...
go test -race ./...  # With race detection
go test -cover ./... # With coverage
```

### Full Test Suite

```bash
./scripts/test.sh
```

## 📋 Commit Guidelines

We follow conventional commits for clear history:

- `feat:` new feature (e.g., `feat: add pagination support`)
- `fix:` bug fix (e.g., `fix: correct schema field name`)
- `docs:` documentation changes (e.g., `docs: update README examples`)
- `test:` adding or updating tests (e.g., `test: add edge cases for filters`)
- `refactor:` code refactoring (e.g., `refactor: simplify canonicalization`)
- `perf:` performance improvements (e.g., `perf: optimize shapeId computation`)
- `chore:` maintenance tasks (e.g., `chore: bump version to v0.2.0`)

## 🐛 Bug Reports

Great bug reports include:

1. **Quick summary** and/or background
2. **Steps to reproduce** - Be specific! Provide sample code if possible
3. **What you expected** to happen
4. **What actually happens**
5. **Notes** (why you think this might be happening, things you tried)

### Bug Report Template

```markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce:
1. ...
2. ...

**Expected behavior**
What you expected to happen.

**Actual behavior**
What actually happened.

**Environment**
- OS: [e.g., macOS 14.0]
- Node.js: [e.g., 20.9.0]
- Go: [e.g., 1.22.0]
- Package version: [e.g., 0.1.0]
```

## 💡 Feature Requests

We welcome feature requests! Please provide:

- **Use case**: Describe the problem you're trying to solve
- **Proposed solution**: How you envision the feature working
- **Alternatives considered**: Other approaches you've thought about
- **Additional context**: Any other relevant information

## ✅ Pull Request Checklist

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
git clone https://github.com/YOUR_USERNAME/includekit-spec.git
cd includekit-spec
./scripts/test.sh  # Installs deps, builds, tests everything
```

- **Do not import testkit packages from production code.** CI will fail if:
  - `@includekit/spec-testkit` is imported outside test files (TypeScript)
  - `github.com/bold-minds/includekit-spec/go/tests` is imported from `pkgs/go/types` (Go)
- All validators, canonicalization, and shapeId utilities live in testkit packages only.
- Production packages (`@includekit/spec`, `includekit-spec/go`) contain types only.
