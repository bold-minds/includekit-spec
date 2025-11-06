# IncludeKit Spec

[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![CI](https://github.com/bold-minds/includekit-spec/actions/workflows/ci.yml/badge.svg)](https://github.com/bold-minds/includekit-spec/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)](https://golang.org/doc/go1.22)
[![Node Version](https://img.shields.io/badge/Node-20%2B-339933?logo=node.js)](https://nodejs.org/)
[![Latest Release](https://img.shields.io/github/v/release/bold-minds/includekit-spec?logo=github&color=blueviolet)](https://github.com/bold-minds/includekit-spec/releases)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/bold-minds/includekit-spec/blob/main/CONTRIBUTING.md)

The **single source of truth** for IncludeKit's Universal Format.  
This repo defines the canonical, language-agnostic **specification** used by:
- SDK adapters (e.g., Prisma) to describe **read queries** and **write events**
- The core Go→WASM engine to compute **dependency summaries** and drive **invalidation**

## What this provides

- **Normative spec:** [`schema/README.md`](schema/README.md)
- **Machine schema:** [`schema/v0-1-0.json`](schema/v0-1-0.json)
- **TypeScript spec:** `@includekit/spec` (production, types-only)
- **TypeScript testkit:** `@includekit/spec-testkit` (validators, JCS, shapeId)
- **Go spec:** `github.com/bold-minds/includekit-spec/go` (production, types-only)
- **Go testkit:** `github.com/bold-minds/includekit-spec/go/tests` (validators, JCS, shapeId)
- **Conformance tests:** Cross-language test vectors (TS ↔ Go)

## Why this exists

- **Determinism:** Same input **QueryShape** → same canonical JSON → same `shapeId` (all languages)
- **Safety:** Precise write descriptions (MutationEvents) prevent stale caches
- **Portability:** SDKs and engine speak identical wire format

---

## Quick Start

### Prerequisites

- **Node.js** ≥ 20
- **Go** ≥ 1.22

### Setup

```bash
# Clone and run automated setup
git clone https://github.com/bold-minds/includekit-spec.git
cd includekit-spec
./scripts/test.sh  # Builds codegen, regenerates types, runs all tests
```

### Repository Structure

```
includekit-spec/
├─ VERSION                    # Single source of truth for version
├─ schema/
│  └─ v0-1-0.json            # JSON Schema (source of truth)
├─ codegen/                   # Go code generator
├─ pkgs/
│  ├─ ts/types/              # TypeScript types (production)
│  ├─ ts/tests/              # TypeScript testkit (dev/test only)
│  └─ go/                    # Go types and tests
├─ tools/
│  ├─ version/sync.go        # Version synchronization tool
│  └─ tests/                 # Test vector generation
└─ scripts/
   └─ test.sh                # Main workflow: build + test + verify
```

**Production packages** (types only, no runtime):
- `@includekit/spec` (TypeScript)
- `github.com/bold-minds/includekit-spec/go` (Go)

**Testkit packages** (validators, JCS, shapeId - dev/test only):
- `@includekit/spec-testkit` (TypeScript)
- `github.com/bold-minds/includekit-spec/go/tests` (Go)

---

## Usage

### TypeScript

```bash
npm install @includekit/spec
```

```typescript
import type { Statement, Query, Filter } from '@includekit/spec';

const statement: Statement = {
  query: {
    model: 'Post',
    where: {
      conditions: [{ field: 'published', op: 'eq', value: true }]
    },
    order_by: [{ field: 'createdAt', descending: true }],
    limit: 10,
  }
};
```

### Go

```bash
go get github.com/bold-minds/includekit-spec/go
```

```go
import "github.com/bold-minds/includekit-spec/go/types"

statement := types.Statement{
  Query: &types.Query{
    Model: "Post",
    Where: &types.Filter{
      Conditions: &[]types.Condition{
        {Field: "published", Op: "eq", Value: true},
      },
    },
    OrderBy: &[]types.OrderBy{
      {Field: "createdAt", Descending: boolPtr(true)},
    },
    Limit: intPtr(10),
  },
}
```

---

## Development

### Making Changes

```bash
# 1. Edit schema
vim schema/v0-1-0.json

# 2. Run tests (auto-regenerates code)
./scripts/test.sh

# 3. Update test vectors if needed
go run tools/tests/generate-vectors.go
```

### Version Bumps

**Single source of truth:** `VERSION` file

```bash
# 1. Update version
echo "0.2.0" > VERSION

# 2. Sync all files (schema, package.json, workflows, etc.)
go run tools/version/sync.go

# 3. Regenerate code with new version
cd codegen && go run .

# 4. Test and commit
./scripts/test.sh
git add -A && git commit -m "chore: bump version to v0.2.0"
git tag v0.2.0 && git push --tags
```

The sync tool automatically updates:
- Schema filename and metadata
- TypeScript package.json files  
- Go codegen default path
- CI/CD workflows
- Generated code headers

### Testing

```bash
./scripts/test.sh    # Builds codegen, regenerates types, runs all tests
```

This runs:
1. Codegen build
2. Type generation (TS + Go validation)
3. Testkit build
4. Conformance tests (TS ↔ Go)
5. No-runtime constraint verification

---

## Versioning

Follows [Semantic Versioning](https://semver.org/):
- **MAJOR:** Breaking schema changes
- **MINOR:** Additive features (new operators, optional fields)
- **PATCH:** Documentation, bug fixes

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on:
- Development setup
- Testing requirements
- Pull request process
- Commit conventions

## 📄 License

This project is licensed under the Apache-2.0 License - see the [LICENSE](LICENSE) file for details.

## 🙏 Community

- **Questions**: Open a [GitHub Discussion](https://github.com/bold-minds/includekit-spec/discussions)
- **Bugs**: Report via [GitHub Issues](https://github.com/bold-minds/includekit-spec/issues)
- **Security**: See [SECURITY.md](SECURITY.md) for responsible disclosure

## 📚 Related Resources

- **IncludeKit Docs**: [includekit.dev](https://includekit.dev) (coming soon)
- **Examples**: Check out the test vectors in `tools/tests/vectors/`
- **Code of Conduct**: [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)