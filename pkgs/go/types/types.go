// Package types provides type definitions for IncludeKit Universal Format v0.1
// This is a PRODUCTION package - types only, no runtime utilities.
//
// IMPORTANT: These types are HAND-WRITTEN (not auto-generated) to preserve idiomatic
// Go patterns like pointer-to-slice for optionals.
// When schema/v0-1-0.json changes, these types must be updated manually.
package types

// Statement is the normalized, language-agnostic description of a read
type Statement struct {
	Query      *Query      `json:"query,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	GroupBy    *[]string   `json:"group_by,omitempty"`
	Having     *Filter     `json:"having,omitempty"`
	Includes   []Include   `json:"includes,omitempty"`
	ORMVersion *string     `json:"orm_version,omitempty"` // diagnostic only
	SDKVersion *string     `json:"sdk_version,omitempty"`
}

type Query struct {
	Model    string     `json:"model"` // target relation name (e.g., "posts", "author")
	Fields   *[]string  `json:"fields,omitempty"`
	Where    *Filter    `json:"where,omitempty"`
	OrderBy  *[]OrderBy `json:"order_by,omitempty"`
	Limit    *int       `json:"limit,omitempty"`
	Offset   *int       `json:"offset,omitempty"`
	Distinct *[]string  `json:"distinct,omitempty"`
}

// Include defines nested relation loading and optional relation-based filtering.
// When Kind is nil, this loads the relation data.
// When Kind is set, this filters the parent records based on the relation.
// When Fields is nil/empty and Kind is set, this filters without loading data.
type Include struct {
	Query    *Query    `json:"query,omitempty"`
	Kind     *string   `json:"kind,omitempty"`     // "some" | "every" | "none" - filters parent by relation
	Includes []Include `json:"includes,omitempty"` // nested includes
}

// Filter composes predicates with boolean logic
type Filter struct {
	And        *[]Filter    `json:"and,omitempty"`
	Or         *[]Filter    `json:"or,omitempty"`
	Not        *Filter      `json:"not,omitempty"`
	Conditions *[]Condition `json:"conditions,omitempty"`
}

// Condition is a leaf-level predicate
type Condition struct {
	Field     string   `json:"field"`
	FieldPath []string `json:"field_path,omitempty"`
	Op        string   `json:"op"`
	Value     any      `json:"value,omitempty"`
}

// OrderBy defines field ordering
type OrderBy struct {
	Field         string `json:"field"`
	Descending    *bool  `json:"descending,omitempty"`     // true = DESCENDING, false = ASCENDING
	NullsFirst    *bool  `json:"nulls_first,omitempty"`    // true = NULLS FIRST, false = NULLS LAST
	CaseSensitive *bool  `json:"case_sensitive,omitempty"` // true = case-sensitive, false = case-insensitive
}

// Pagination defines cursor-based pagination parameters.
// Uses opaque cursors (base64-encoded JSON) for SDK abstraction.
// Forward pagination: use First + After
// Backward pagination: use Last + Before
type Pagination struct {
	First  *int    `json:"first,omitempty"`  // Forward limit
	Last   *int    `json:"last,omitempty"`   // Backward limit
	After  *string `json:"after,omitempty"`  // Opaque cursor to start after (forward)
	Before *string `json:"before,omitempty"` // Opaque cursor to start before (backward)
}

// Mutation describes writes that could affect reads
type Mutation struct {
	TxID    *string  `json:"tx_id,omitempty"`
	Changes []Change `json:"changes"`
}

// Change represents a single mutation operation (insert/update/delete)
type Change struct {
	Model  string  `json:"model"`
	Action string  `json:"action"` // "insert" | "update" | "delete"
	Sets   []KV    `json:"sets,omitempty"`
	Where  *Filter `json:"where,omitempty"`
}

// Dependencies tracks what a read depends on (engine output)
type Dependencies struct {
	ShapeID  string              `json:"shape_id"`
	Records  map[string][]string `json:"records"`
	Filters  []Filter            `json:"filters"`
	Includes []Include           `json:"includes"` // includes with Kind set
	LastRow  *PaginationBoundary `json:"last_row,omitempty"`
	GroupBy  *GroupByKV          `json:"group_by,omitempty"`
}

// PaginationBoundary tracks the last included row for paginated queries
type PaginationBoundary struct {
	OrderBy []OrderBy `json:"order_by"`
	// Field values of the last included row
	Row map[string]any `json:"row"`
	// Cursor identifies the stable pagination cursor
	Cursor *KV `json:"cursor,omitempty"`
}

// GroupByKV tracks group-by dimensions
type GroupByKV struct {
	Keys   []string         `json:"keys"`
	Values []map[string]any `json:"values"`
}

type KV struct {
	Field string `json:"field"`
	Value any    `json:"value"`
}
