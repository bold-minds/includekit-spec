// Package mock provides a mock IncludeKit engine for testing without WASM
package mock

import (
	"github.com/bold-minds/includekit-spec/go/types"
)

// AppSchema is engine-specific (not in universal format spec)
type AppSchema struct {
	Version int     `json:"version"`
	Models  []Model `json:"models"`
}

// Model represents a model in the schema
type Model struct {
	Name      string     `json:"name"`
	ID        IDConfig   `json:"id"`
	Relations []Relation `json:"relations,omitempty"`
}

// IDConfig represents ID field configuration
type IDConfig struct {
	Kind string `json:"kind"`
}

// Relation represents a model relation
type Relation struct {
	Name   string `json:"name"`
	Target string `json:"target"`
	Kind   string `json:"kind"`
}

// AddQueryRequest wraps a shape with optional result hint
type AddQueryRequest struct {
	Shape      types.Statement          `json:"shape"`
	ResultHint map[string][]interface{} `json:"result_hint,omitempty"`
}

// AddQueryResponse contains shape ID and dependencies
type AddQueryResponse struct {
	ShapeID      string             `json:"shape_id"`
	Dependencies types.Dependencies `json:"dependencies"`
}

// ShapeIDResponse contains the computed shape ID
type ShapeIDResponse struct {
	ShapeID string `json:"shape_id"`
}

// InvalidateResponse contains shape IDs to evict
type InvalidateResponse struct {
	Evict []string `json:"evict"`
}

// ExplainRequest contains mutation and shape ID for explanation
type ExplainRequest struct {
	Mutation types.Mutation `json:"mutation"`
	ShapeID  string         `json:"shape_id"`
}

// ExplainResponse explains why a shape would be invalidated
type ExplainResponse struct {
	Invalidate bool     `json:"invalidate"`
	Reasons    []string `json:"reasons"`
}

// VersionInfo contains engine version information
type VersionInfo struct {
	Core     string `json:"core"`
	Contract string `json:"contract"`
	ABI      string `json:"abi"`
}

// Engine interface matching WASM exports
type Engine interface {
	SetSchema(schema AppSchema) error
	ComputeShapeID(statement types.Statement) (ShapeIDResponse, error)
	AddQuery(request AddQueryRequest) (AddQueryResponse, error)
	Invalidate(mutation types.Mutation) (InvalidateResponse, error)
	ExplainInvalidation(request ExplainRequest) (ExplainResponse, error)
	Reset()
	GetVersion() VersionInfo
}
