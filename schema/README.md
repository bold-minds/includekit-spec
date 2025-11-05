# IncludeKit Universal Format v0.1 — Complete Reference

This document provides comprehensive documentation for every type and field in the IncludeKit Universal Format, explaining **when**, **why**, and **how** each is used.

---

## Table of Contents

1. [Statement](#statement) - Top-level read operation
2. [Query](#query) - Core query parameters
3. [Include](#include) - Nested relation loading
4. [Filter](#filter) - Boolean predicate logic
5. [Condition](#condition) - Leaf-level predicates
6. [OrderBy](#orderby) - Sort specifications
7. [Pagination](#pagination) - Cursor-based pagination
8. [Mutation](#mutation) - Write operations
9. [Change](#change) - Individual write changes
10. [KV](#kv) - Key-value pairs
11. [Dependencies](#dependencies) - Cache tracking (engine output)
12. [PaginationBoundary](#paginationboundary) - Page edge tracking
13. [GroupByKV](#groupbykv) - Aggregation keys

---

## Statement

The top-level type representing a complete read operation. This is what adapters send to the engine.

```go
type Statement struct {
    Query      *Query      `json:"query,omitempty"`
    Pagination *Pagination `json:"pagination,omitempty"`
    GroupBy    *[]string   `json:"group_by,omitempty"`
    Having     *Filter     `json:"having,omitempty"`
    Includes   []Include   `json:"includes,omitempty"`
    ORMVersion *string     `json:"orm_version,omitempty"` // diagnostic only
    SDKVersion *string     `json:"sdk_version,omitempty"`
}
```

### Fields

#### `Query` (*Query)
- **When**: Always present when fetching data
- **Why**: Specifies what model to query and basic filtering/sorting
- **Example**: Query for published posts
  ```json
  {
    "query": {
      "model": "posts",
      "where": { "conditions": [{"field": "status", "op": "eq", "value": "published"}] }
    }
  }
  ```

#### `Pagination` (*Pagination)
- **When**: Implementing cursor-based pagination
- **Why**: Enables stable pagination across result sets
- **Example**: Get first 20 posts after a cursor
  ```json
  {
    "query": { "model": "posts" },
    "pagination": {
      "first": 20,
      "after": "eyJpZCI6InBvc3RfMTIzIn0="
    }
  }
  ```

#### `GroupBy` (*[]string)
- **When**: Performing aggregations (COUNT, SUM, etc.)
- **Why**: Groups rows by specified fields for aggregate functions
- **Example**: Count posts per author
  ```json
  {
    "query": {
      "model": "posts",
      "fields": ["author_id", "COUNT(*) as count"]
    },
    "group_by": ["author_id"]
  }
  ```

#### `Having` (*Filter)
- **When**: Filtering aggregated results
- **Why**: Applies predicates after grouping (like SQL HAVING)
- **Example**: Only authors with more than 5 posts
  ```json
  {
    "query": { "model": "posts", "fields": ["author_id", "COUNT(*) as count"] },
    "group_by": ["author_id"],
    "having": { "conditions": [{"field": "count", "op": "gt", "value": 5}] }
  }
  ```

#### `Includes` ([]Include)
- **When**: Loading related data (joins)
- **Why**: Enables fetching parent/child relationships in one query
- **Example**: Get posts with their authors and comments
  ```json
  {
    "query": { "model": "posts" },
    "includes": [
      { "query": { "model": "author" } },
      { "query": { "model": "comments", "limit": 5 } }
    ]
  }
  ```

#### `ORMVersion`, `SDKVersion` (*string)
- **When**: Always set by adapters for diagnostics
- **Why**: Helps debug adapter issues, excluded from cache keys
- **Example**: `"orm_version": "prisma@5.0.0", "sdk_version": "@includekit/prisma@1.0.0"`

---

## Query

Core query parameters - the "SELECT * FROM model WHERE..." part.

```go
type Query struct {
    Model    string     `json:"model"`
    Fields   *[]string  `json:"fields,omitempty"`
    Where    *Filter    `json:"where,omitempty"`
    OrderBy  *[]OrderBy `json:"order_by,omitempty"`
    Limit    *int       `json:"limit,omitempty"`
    Offset   *int       `json:"offset,omitempty"`
    Distinct *[]string  `json:"distinct,omitempty"`
}
```

### Fields

#### `Model` (string, required)
- **When**: Every query
- **Why**: Identifies which table/collection to query
- **Example**: `"model": "posts"` queries the posts table

#### `Fields` (*[]string)
- **When**: Selecting specific columns (projection)
- **Why**: Reduces data transfer, enables aggregates
- **Example**: Select only id and title
  ```json
  { "model": "posts", "fields": ["id", "title"] }
  ```

#### `Where` (*Filter)
- **When**: Filtering rows
- **Why**: Applies predicates to limit results
- **Example**: Only published posts
  ```json
  {
    "model": "posts",
    "where": { "conditions": [{"field": "published", "op": "eq", "value": true}] }
  }
  ```

#### `OrderBy` (*[]OrderBy)
- **When**: Sorting results
- **Why**: Deterministic ordering for pagination and UX
- **Example**: Sort by creation date descending
  ```json
  {
    "model": "posts",
    "order_by": [{"field": "created_at", "descending": true}]
  }
  ```

#### `Limit` (*int)
- **When**: Restricting result count
- **Why**: Pagination, performance
- **Example**: Get top 10 posts
  ```json
  { "model": "posts", "limit": 10 }
  ```

#### `Offset` (*int)
- **When**: Offset-based pagination
- **Why**: Skip N rows (less stable than cursors)
- **Example**: Skip first 20, get next 10
  ```json
  { "model": "posts", "limit": 10, "offset": 20 }
  ```

#### `Distinct` (*[]string)
- **When**: Removing duplicate rows
- **Why**: Get unique values for specific fields
- **Example**: Unique author IDs
  ```json
  { "model": "posts", "fields": ["author_id"], "distinct": ["author_id"] }
  ```

---

## Include

Defines nested relation loading (joins) and relation-based filtering.

```go
type Include struct {
    Query    *Query    `json:"query,omitempty"`
    Kind     *string   `json:"kind,omitempty"`     // "some" | "every" | "none"
    Includes []Include `json:"includes,omitempty"` // nested includes
}
```

### Fields

#### `Query` (*Query)
- **When**: Loading relation data
- **Why**: Specifies what/how to load from the related model
- **Example**: Load first 5 comments per post
  ```json
  {
    "query": { "model": "comments", "limit": 5, "order_by": [{"field": "created_at"}] }
  }
  ```

#### `Kind` (*string)
- **When**: Filtering parent by relation existence
- **Why**: "Show me users who have SOME/EVERY/NO posts matching X"
- **Values**: `"some"`, `"every"`, `"none"`
- **Example**: Users with at least one published post
  ```json
  {
    "kind": "some",
    "query": {
      "model": "posts",
      "where": { "conditions": [{"field": "published", "op": "eq", "value": true}] }
    }
  }
  ```

**Semantics:**
- `"some"`: At least one related row matches
- `"every"`: All related rows match (empty = true)
- `"none"`: No related rows match (empty = true)

#### `Includes` ([]Include)
- **When**: Loading nested relations (multi-level joins)
- **Why**: Get posts → comments → replies in one query
- **Example**: Posts with comments and their authors
  ```json
  {
    "query": { "model": "posts" },
    "includes": [
      {
        "query": { "model": "comments" },
        "includes": [
          { "query": { "model": "author" } }
        ]
      }
    ]
  }
  ```

---

## Filter

Boolean logic for composing predicates (AND/OR/NOT).

```go
type Filter struct {
    And        *[]Filter    `json:"and,omitempty"`
    Or         *[]Filter    `json:"or,omitempty"`
    Not        *Filter      `json:"not,omitempty"`
    Conditions *[]Condition `json:"conditions,omitempty"`
}
```

### Fields

#### `And` (*[]Filter)
- **When**: All conditions must be true
- **Why**: Combine multiple predicates (intersection)
- **Example**: Published AND featured posts
  ```json
  {
    "and": [
      { "conditions": [{"field": "published", "op": "eq", "value": true}] },
      { "conditions": [{"field": "featured", "op": "eq", "value": true}] }
    ]
  }
  ```

#### `Or` (*[]Filter)
- **When**: Any condition can be true
- **Why**: Alternative predicates (union)
- **Example**: Posts with high views OR featured
  ```json
  {
    "or": [
      { "conditions": [{"field": "views", "op": "gte", "value": 1000}] },
      { "conditions": [{"field": "featured", "op": "eq", "value": true}] }
    ]
  }
  ```

#### `Not` (*Filter)
- **When**: Negating a condition
- **Why**: Exclusion logic
- **Example**: Posts not deleted
  ```json
  {
    "not": { "conditions": [{"field": "deleted", "op": "eq", "value": true}] }
  }
  ```

#### `Conditions` (*[]Condition)
- **When**: Leaf-level predicates
- **Why**: Actual field comparisons
- **Example**: Status equals "published"
  ```json
  {
    "conditions": [{"field": "status", "op": "eq", "value": "published"}]
  }
  ```

---

## Condition

Leaf-level predicate (field op value).

```go
type Condition struct {
    Field     string   `json:"field"`
    FieldPath []string `json:"field_path,omitempty"`
    Op        string   `json:"op"`
    Value     any      `json:"value,omitempty"`
}
```

### Fields

#### `Field` (string, required)
- **When**: Every condition
- **Why**: Identifies which column to filter
- **Example**: `"field": "status"`

#### `FieldPath` ([]string)
- **When**: Filtering on nested JSON fields
- **Why**: Navigate into JSON columns (e.g., PostgreSQL JSONB)
- **Example**: Filter on metadata.user.country
  ```json
  {
    "field": "metadata",
    "field_path": ["user", "country"],
    "op": "eq",
    "value": "US"
  }
  ```
  Generates: `metadata->'user'->>'country' = 'US'`

#### `Op` (string, required)
- **When**: Every condition
- **Why**: Defines the comparison operation
- **Values**:
  - Equality: `eq`, `ne`, `in`, `notIn`, `isNull`
  - Numeric: `gt`, `gte`, `lt`, `lte`, `between`
  - Text: `contains`, `startsWith`, `endsWith`, `like`, `ilike`, `regex`
  - Arrays/JSON: `has`, `hasSome`, `hasEvery`, `jsonContains`
  - Length: `lenEq`, `lenGt`, `lenLt`
  - Relation: `exists`
  - Extension: `custom:*`
- **Example**: `"op": "gte"` for greater-than-or-equal

#### `Value` (any)
- **When**: Most operators (except `isNull`, `exists`)
- **Why**: The comparison value
- **Types**: string, number, boolean, null, array, object
- **Examples**:
  ```json
  {"field": "age", "op": "gte", "value": 18}
  {"field": "status", "op": "in", "value": ["published", "featured"]}
  {"field": "email", "op": "contains", "value": "@example.com"}
  ```

---

## OrderBy

Sort specification for result ordering.

```go
type OrderBy struct {
    Field         string `json:"field"`
    Descending    *bool  `json:"descending,omitempty"`
    NullsFirst    *bool  `json:"nulls_first,omitempty"`
    CaseSensitive *bool  `json:"case_sensitive,omitempty"`
}
```

### Fields

#### `Field` (string, required)
- **When**: Every sort spec
- **Why**: Column to sort by
- **Example**: `"field": "created_at"`

#### `Descending` (*bool)
- **When**: Controlling sort direction
- **Why**: DESC vs ASC ordering
- **Default**: false (ascending)
- **Example**: Sort newest first
  ```json
  {"field": "created_at", "descending": true}
  ```

#### `NullsFirst` (*bool)
- **When**: Controlling NULL placement
- **Why**: Some databases default NULLS LAST, some FIRST
- **Example**: Show posts without dates first
  ```json
  {"field": "published_at", "descending": true, "nulls_first": true}
  ```

#### `CaseSensitive` (*bool)
- **When**: Sorting text fields
- **Why**: Control case-sensitive vs case-insensitive collation
- **Example**: Case-insensitive title sort
  ```json
  {"field": "title", "case_sensitive": false}
  ```

---

## Pagination

Cursor-based pagination (Relay-style).

```go
type Pagination struct {
    First  *int    `json:"first,omitempty"`
    Last   *int    `json:"last,omitempty"`
    After  *string `json:"after,omitempty"`
    Before *string `json:"before,omitempty"`
}
```

### Fields

#### `First` (*int)
- **When**: Forward pagination
- **Why**: Limit for "next N items"
- **Example**: Get first 20 posts
  ```json
  {"first": 20}
  ```

#### `Last` (*int)
- **When**: Backward pagination
- **Why**: Limit for "previous N items"
- **Example**: Get last 20 posts
  ```json
  {"last": 20}
  ```

#### `After` (*string)
- **When**: Forward pagination from a specific point
- **Why**: Resume from where previous page ended
- **Format**: Opaque cursor (base64-encoded JSON)
- **Example**: Next page after cursor
  ```json
  {"first": 20, "after": "eyJpZCI6InBvc3RfMTIzIiwiY3JlYXRlZEF0IjoiMjAyNC0wMS0xNSJ9"}
  ```

#### `Before` (*string)
- **When**: Backward pagination before a specific point
- **Why**: Go back from current position
- **Example**: Previous page before cursor
  ```json
  {"last": 20, "before": "eyJpZCI6InBvc3RfMTIzIn0="}
  ```

**Rules:**
- Cannot mix forward (first/after) with backward (last/before)
- Cursors are opaque - SDKs create/decode them, clients treat as strings

---

## Mutation

Describes write operations that may invalidate cached reads.

```go
type Mutation struct {
    TxID    *string  `json:"tx_id,omitempty"`
    Changes []Change `json:"changes"`
}
```

### Fields

#### `TxID` (*string)
- **When**: Available transaction ID
- **Why**: Correlate related changes, debugging
- **Example**: `"tx_id": "tx_abc123xyz"`

#### `Changes` ([]Change, required)
- **When**: Every write
- **Why**: Enumerates all modifications
- **Example**: Insert one post, update another
  ```json
  {
    "changes": [
      {
        "model": "posts",
        "action": "insert",
        "sets": [{"field": "id", "value": "post_1"}, {"field": "title", "value": "Hello"}]
      },
      {
        "model": "posts",
        "action": "update",
        "sets": [{"field": "status", "value": "published"}],
        "where": {"conditions": [{"field": "id", "op": "eq", "value": "post_2"}]}
      }
    ]
  }
  ```

---

## Change

Individual mutation operation (insert/update/delete).

```go
type Change struct {
    Model  string  `json:"model"`
    Action string  `json:"action"` // "insert" | "update" | "delete"
    Sets   []KV    `json:"sets,omitempty"`
    Where  *Filter `json:"where,omitempty"`
}
```

### Fields

#### `Model` (string, required)
- **When**: Every change
- **Why**: Identifies which table is modified
- **Example**: `"model": "posts"`

#### `Action` (string, required)
- **When**: Every change
- **Why**: Type of operation
- **Values**: `"insert"`, `"update"`, `"delete"`
- **Example**: `"action": "insert"`

#### `Sets` ([]KV)
- **When**: insert, update
- **Why**: Field assignments (SET clause in SQL)
- **Required for**: insert, update
- **Example**: Set title and status
  ```json
  {
    "action": "insert",
    "model": "posts",
    "sets": [
      {"field": "title", "value": "New Post"},
      {"field": "status", "value": "draft"}
    ]
  }
  ```

#### `Where` (*Filter)
- **When**: update, delete
- **Why**: Identifies which rows to modify
- **Required for**: update, delete
- **Example**: Update specific post
  ```json
  {
    "action": "update",
    "model": "posts",
    "sets": [{"field": "status", "value": "published"}],
    "where": {"conditions": [{"field": "id", "op": "eq", "value": "post_1"}]}
  }
  ```

**Validation Rules:**
- `insert`: Requires `sets`, no `where`
- `update`: Requires `sets` and `where`
- `delete`: Requires `where`, no `sets`

---

## KV

Simple key-value pair (field assignment or cursor field).

```go
type KV struct {
    Field string `json:"field"`
    Value any    `json:"value"`
}
```

### Fields

#### `Field` (string, required)
- **When**: Always
- **Why**: Column name
- **Example**: `"field": "status"`

#### `Value` (any, required)
- **When**: Always
- **Why**: The assigned value
- **Example**: `"value": "published"`

**Used in:**
1. **Change.Sets**: Field assignments in INSERT/UPDATE
2. **PaginationBoundary.Cursor**: Stable cursor field

---

## Dependencies

Engine output tracking what a cached result depends on (for invalidation).

```go
type Dependencies struct {
    ShapeID  string              `json:"shape_id"`
    Records  map[string][]string `json:"records"`
    Filters  []Filter            `json:"filters"`
    Includes []Include           `json:"includes"`
    LastRow  *PaginationBoundary `json:"last_row,omitempty"`
    GroupBy  *GroupByKV          `json:"group_by,omitempty"`
}
```

### Fields

#### `ShapeID` (string, required)
- **When**: Every query result
- **Why**: Cache key (deterministic hash of query)
- **Format**: `"s_" + hex(SHA-256(JCS(Statement)))`
- **Example**: `"shape_id": "s_a1b2c3d4..."`

#### `Records` (map[string][]string, required)
- **When**: Every query result
- **Why**: Track which specific rows were returned
- **Format**: `{"model": ["id1", "id2", ...]}`
- **Example**:
  ```json
  {
    "records": {
      "posts": ["post_1", "post_2", "post_3"],
      "users": ["user_10", "user_20"]
    }
  }
  ```
- **Invalidation**: If any listed ID is modified/deleted, invalidate

#### `Filters` ([]Filter, required)
- **When**: Query has WHERE/HAVING clauses
- **Why**: Detect when changes cross filter boundaries
- **Example**: If cached results had `views >= 100`, a post changing from 99 → 101 crosses the boundary
- **Invalidation**: If old→new value crosses any filter predicate, invalidate

#### `Includes` ([]Include, required)
- **When**: Query has relation filters (kind: some/every/none)
- **Why**: Track relation-based filtering
- **Example**: Users with "some" published posts - if a post's status changes, invalidate
- **Invalidation**: If link/unlink or child changes affect relation semantics, invalidate

#### `LastRow` (*PaginationBoundary)
- **When**: Query uses pagination with ORDER BY
- **Why**: Detect when new rows would fall into the page
- **Example**: Cached top 10 posts by views - if a post gets more views than #10, invalidate
- **Invalidation**: If a write creates/updates a row that sorts into the window, invalidate

#### `GroupBy` (*GroupByKV)
- **When**: Query uses GROUP BY
- **Why**: Track which groups existed
- **Example**: Cached counts per author - if a new author appears, invalidate
- **Invalidation**: If group membership changes, invalidate

---

## PaginationBoundary

Tracks the edge of a paginated result set for invalidation.

```go
type PaginationBoundary struct {
    OrderBy []OrderBy      `json:"order_by"`
    Row     map[string]any `json:"row"`
    Cursor  *KV            `json:"cursor,omitempty"`
}
```

### Fields

#### `OrderBy` ([]OrderBy, required)
- **When**: Always (matches query's order)
- **Why**: Defines sort direction for boundary comparison
- **Example**:
  ```json
  {
    "order_by": [
      {"field": "views", "descending": true},
      {"field": "id"}
    ]
  }
  ```

#### `Row` (map[string]any, required)
- **When**: Always
- **Why**: Field values of the last included row
- **Example**: Last post in top 10 by views
  ```json
  {
    "row": {
      "views": 250,
      "id": "post_789",
      "created_at": "2024-01-15T10:30:00Z"
    }
  }
  ```
- **Usage**: Compare new/updated rows against these values to detect boundary crossing

#### `Cursor` (*KV)
- **When**: Stable tiebreaker field available (usually ID)
- **Why**: Ensures deterministic pagination even with duplicate sort values
- **Example**:
  ```json
  {"cursor": {"field": "id", "value": "post_789"}}
  ```

---

## GroupByKV

Tracks aggregation group keys for invalidation.

```go
type GroupByKV struct {
    Keys   []string         `json:"keys"`
    Values []map[string]any `json:"values"`
}
```

### Fields

#### `Keys` ([]string, required)
- **When**: Always (GROUP BY fields)
- **Why**: Identifies grouping dimensions
- **Example**: `"keys": ["author_id", "status"]`

#### `Values` ([]map[string]any, required)
- **When**: Always
- **Why**: Actual group key combinations that existed
- **Example**: Groups that appeared in result
  ```json
  {
    "keys": ["author_id", "status"],
    "values": [
      {"author_id": "user_1", "status": "published"},
      {"author_id": "user_2", "status": "draft"},
      {"author_id": "user_2", "status": "published"}
    ]
  }
  ```
- **Invalidation**: If a write creates a new group combination, invalidate

---

## Canonicalization & Shape Hashing

The engine produces deterministic `ShapeID` values for cache keys:

1. **Remove diagnostics**: Strip `orm_version`, `sdk_version`
2. **Apply JCS** (JSON Canonicalization Scheme):
   - Sort object keys lexicographically
   - Preserve array order
   - Canonical number encoding
3. **Hash**: SHA-256 of canonical JSON string
4. **Prefix**: `"s_" + hex(hash)`

**Example:**
```
Statement → Remove diagnostics → JCS → SHA-256 → "s_a1b2c3d4..."
```

---

## Invalidation Rules Summary

Given a `Mutation` and cached `Dependencies`:

1. **Record Match**: Any changed/deleted ID in `Records` → invalidate
2. **Filter Crossing**: Old→new value crosses any `Filters` predicate → invalidate
3. **Relation Change**: Link/unlink or child changes affecting `Includes` semantics → invalidate
4. **Boundary Shift**: Write creates/updates row sorting into `LastRow` window → invalidate
5. **Group Change**: Write creates new group in `GroupBy` → invalidate
6. **Unknown Operator**: Any `custom:*` operator in bounds → invalidate conservatively

---

## JSON Schema

Machine-readable schema: `schema/v0-1-0.json`

All types, fields, and validation rules are formally defined in the JSON Schema file.

---

## Versioning

- **Breaking changes**: Bump MAJOR (remove fields, change semantics)
- **Additive changes**: Bump MINOR (new fields, new operators)
- **Fixes**: Bump PATCH (documentation, implementation bugs)

Current version: **v0.1.0** (pre-release)
