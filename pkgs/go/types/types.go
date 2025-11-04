// Package types provides type definitions for IncludeKit Universal Format v0.1
// This is a PRODUCTION package - types only, no runtime utilities.
//
// IMPORTANT: These types are HAND-WRITTEN (not auto-generated) to preserve idiomatic
// Go patterns like the sealed Scalar interface and pointer-to-slice for optionals.
// When schema/v0-1-0.json changes, these types must be updated manually.
//
// TypeScript types ARE auto-generated. See GO_GENERATION_ANALYSIS.md for details.
package types

// Scalar represents a supported value type in filters/values
type Scalar interface {
	isScalar()
}

type ScalarString string
type ScalarNumber float64
type ScalarBool bool
type ScalarNull struct{}
type ScalarJSON struct {
	JSON interface{} `json:"json"`
}

func (ScalarString) isScalar() {}
func (ScalarNumber) isScalar() {}
func (ScalarBool) isScalar()   {}
func (ScalarNull) isScalar()   {}
func (ScalarJSON) isScalar()   {}

// OrderBySpec defines field ordering
type OrderBySpec struct {
	Field     string  `json:"field"`
	Direction string  `json:"direction"`           // "asc" | "desc"
	Nulls     *string `json:"nulls,omitempty"`     // "first" | "last"
	Collation *string `json:"collation,omitempty"` // "ci" | "cs"
}

// CursorSpec defines cursor-based pagination
type CursorSpec struct {
	Fields []string `json:"fields"`
	Values []Scalar `json:"values"`
	Before *bool    `json:"before,omitempty"`
}

// FilterAtom is a leaf-level predicate
type FilterAtom struct {
	Field string      `json:"field"`
	Op    string      `json:"op"`
	Value interface{} `json:"value,omitempty"`
	Path  []string    `json:"path,omitempty"`
}

// FilterSpec composes predicates with boolean logic
type FilterSpec struct {
	And   *[]FilterSpec `json:"and,omitempty"`
	Or    *[]FilterSpec `json:"or,omitempty"`
	Not   *FilterSpec   `json:"not,omitempty"`
	Atoms *[]FilterAtom `json:"atoms,omitempty"`
}

// RelationFilterBound captures relational predicates
type RelationFilterBound struct {
	Relation string      `json:"relation"`
	Kind     string      `json:"kind"` // "some" | "every" | "none"
	Where    *FilterSpec `json:"where,omitempty"`
}

// IncludeSpec defines nested relation loading
type IncludeSpec struct {
	Select         *[]string              `json:"select,omitempty"`
	Where          *FilterSpec            `json:"where,omitempty"`
	OrderBy        *[]OrderBySpec         `json:"orderBy,omitempty"`
	Take           *int                   `json:"take,omitempty"`
	Skip           *int                   `json:"skip,omitempty"`
	Distinct       *[]string              `json:"distinct,omitempty"`
	Include        map[string]IncludeSpec `json:"include,omitempty"`
	RelationFilter *[]RelationFilterBound `json:"relationFilter,omitempty"`
	SeparateLoad   *bool                  `json:"separateLoad,omitempty"`
}

// QueryShape is the normalized, language-agnostic description of a read
type QueryShape struct {
	Model          string                 `json:"model"`
	Select         *[]string              `json:"select,omitempty"`
	Where          *FilterSpec            `json:"where,omitempty"`
	OrderBy        *[]OrderBySpec         `json:"orderBy,omitempty"`
	Take           *int                   `json:"take,omitempty"`
	Skip           *int                   `json:"skip,omitempty"`
	Cursor         *CursorSpec            `json:"cursor,omitempty"`
	Distinct       *[]string              `json:"distinct,omitempty"`
	GroupBy        *[]string              `json:"groupBy,omitempty"`
	Having         *FilterSpec            `json:"having,omitempty"`
	Include        map[string]IncludeSpec `json:"include,omitempty"`
	ORM            *string                `json:"orm,omitempty"` // diagnostic only
	AdapterVersion *string                `json:"adapterVersion,omitempty"`
}

// LinkChange represents link/unlink operations
type LinkChange struct {
	Kind        string `json:"kind"` // "link" | "unlink"
	ParentModel string `json:"parentModel"`
	ParentID    string `json:"parentId"`
	Relation    string `json:"relation"`
	ChildModel  string `json:"childModel"`
	ChildID     string `json:"childId"`
}

// WriteChange represents create/update/delete operations
type WriteChange struct {
	Op     string                 `json:"op"` // "create" | "update" | "delete" | "link" | "unlink"
	Model  string                 `json:"model,omitempty"`
	ID     string                 `json:"id,omitempty"`
	Before map[string]interface{} `json:"before,omitempty"`
	After  map[string]interface{} `json:"after,omitempty"`
	// For link/unlink operations
	ParentModel *string `json:"parentModel,omitempty"`
	ParentID    *string `json:"parentId,omitempty"`
	Relation    *string `json:"relation,omitempty"`
	ChildModel  *string `json:"childModel,omitempty"`
	ChildID     *string `json:"childId,omitempty"`
}

// MutationEvent describes writes that could affect reads
type MutationEvent struct {
	TxID    *string       `json:"txId,omitempty"`
	Changes []WriteChange `json:"changes"`
}

// SortThreshold defines boundaries for topN tracking
type SortThreshold struct {
	OrderBy    []OrderBySpec          `json:"orderBy"`
	Boundary   map[string]interface{} `json:"boundary"`
	TieBreaker *struct {
		Field string      `json:"field"`
		Value interface{} `json:"value"`
	} `json:"tieBreaker,omitempty"`
}

// Dependencies tracks what a read depends on (engine output)
type Dependencies struct {
	ShapeID        string                `json:"shapeId"`
	Records        map[string][]string   `json:"records"`
	FilterBounds   []FilterSpec          `json:"filterBounds"`
	RelationBounds []RelationFilterBound `json:"relationBounds"`
	TopN           *SortThreshold        `json:"topN,omitempty"`
	Groups         *struct {
		Keys   []string                 `json:"keys"`
		Values []map[string]interface{} `json:"values"`
	} `json:"groups,omitempty"`
}
