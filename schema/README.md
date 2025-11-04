# IncludeKit Universal Format — v1.0 (Normative)

This document defines the **Universal Format** used by all IncludeKit SDKs and the core Go/WASM engine.

## 1. Goals

- Cross-ORM expressiveness (Prisma, TypeORM, Sequelize, GORM, EF Core, Hibernate/JPA, Django ORM, SQLAlchemy, Drizzle)
- Deterministic **canonical JSON (JCS)** and **shape hashing** across languages
- Cache-safety: conservative invalidation when uncertain
- Simplicity: adapters map **what was requested**, not how the ORM executes

> **Implementation boundary (non-normative note):**  
> This document is normative for data shapes and algorithms (JCS + sha256, dependency semantics). In **production**, SDKs do **not** implement these algorithms; the Go→WASM core implements them. Reference validators/JCS exist only in separate *types-testkit* packages for CI/conformance.

## 2. Scalars & IDs

- `Scalar = string | number | boolean | null | { json: any }`
  - Dates: ISO-8601 strings; decimals/bigints: strings
- IDs are **strings**. Composite keys encode as:
  - `"fieldA=<valA>|fieldB=<valB>"` with field names sorted lexicographically

## 3. QueryShape

A normalized, language-agnostic description of a read.

- **Top-level**: `model`, `select`, `where`, `orderBy`, `take`, `skip`, `cursor`, `distinct`, `groupBy`, `having`
- **Relations**: each include may specify `where`, `orderBy`, `take`, `skip`, `distinct`, nested `include`, and `relationFilter` (`some | every | none`)
- Optional hints (diagnostics only): `orm`, `adapterVersion` — **excluded** from canonicalization

## 4. Filters

- Leaf predicates are `FilterAtom { field, op, value?, path? }`
- `op` includes equality, comparisons, text, json/array, `exists`, and `custom:*` (escape hatch)
- `FilterSpec` composes atoms via `and | or | not`

## 5. Relation filters

`RelationFilterBound { relation, kind: some|every|none, where? }`

- Captures relational predicates that affect inclusion even without direct foreign-key equality

## 6. Writes (MutationEvent)

Adapters must enumerate every write that could affect reads:

- `create`, `update`, `delete` (with `before`/`after` values)
- `link`/`unlink` for relationship changes
- Bulk ops must **prefetch IDs** and expand cascades; otherwise the adapter should refuse upstream

## 7. Dependencies (engine output per read)

- `shapeId`: `"s_" + sha256(JCS(QueryShape))`
- `records`: model → list of IDs observed in the result
- `filterBounds`: predicates that, if crossed by a future write, change inclusion
- `relationBounds`: relational `some/every/none` semantics
- `topN`: for `orderBy+take`, the boundary values of the last included row
- `groups`: (optional) group keys for aggregate shapes

## 8. Canonicalization & hashing

- Remove `orm` and `adapterVersion`
- Apply **JSON Canonicalization Scheme (JCS)**:
  1. Object keys sorted lexicographically
  2. Arrays remain ordered
  3. Numbers remain JSON numbers; bigints/decimals stay as strings
- Hash the canonical string with **SHA-256**; prefix with `"s_"`

> **Production note:** The canonicalization and hashing steps are implemented by the core engine (WASM) in production; any reference implementations in *types-testkit* packages are for testing only.

## 9. Invalidation (engine rules, summarized)

Given a `MutationEvent` and cached `Dependencies`:
1. **Record strategy:** any changed/deleted ID in `records` → invalidate
2. **Filter strategy:** any old→new crossing of a `filterBounds`/`having` predicate → invalidate
3. **Relation strategy:** link/unlink or child changes that flip `some/every/none` semantics → invalidate
4. **Top-N strategy:** a write crossing the boundary in `topN` (or modifying the boundary row) → invalidate
5. **Groups/Distinct:** group membership change → invalidate (conservatively if keys absent)
6. **Unknown operator:** any `custom:*` bound → invalidate on any write to that field’s model

> The *engine* implements these rules. The spec merely defines the shapes and events.

## 10. JSON Schema

Authoritative machine form: `schema/v*.json`

## 11. Conformance vectors

Files in `conformance/vectors/` provide:
- Example shapes → canonical JSON and `shapeId`
- Example filters + mutations → expected “invalidate? yes|no” (for engine harnesses)

All language packages must pass conformance tests.

## 12. Versioning

- Breaking schema changes bump **MAJOR**.
- Additive fields/operators bump **MINOR**.
- Implementation/documentation fixes bump **PATCH**.
