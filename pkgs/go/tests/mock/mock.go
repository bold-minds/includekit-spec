package mock

import (
	"fmt"
	"sync"

	"github.com/bold-minds/includekit-spec/go/tests"
	"github.com/bold-minds/includekit-spec/go/types"
)

// MockEngineConfig configures the mock engine behavior
type MockEngineConfig struct {
	ShapeIDGenerator func(types.Statement) string
	EvictBehavior    string // "conservative" | "custom"
	CustomEvictList  []string
	TrackCalls       bool
}

// MockEngineCalls tracks all method calls when TrackCalls is enabled
type MockEngineCalls struct {
	SetSchema           []AppSchema
	ComputeShapeID      []types.Statement
	AddQuery            []AddQueryRequest
	Invalidate          []types.Mutation
	ExplainInvalidation []ExplainRequest
	Reset               []struct{}
	GetVersion          []struct{}
}

// MockEngine implements the Engine interface for testing
type MockEngine struct {
	mu     sync.RWMutex
	schema *AppSchema
	shapes map[string]types.Dependencies
	calls  MockEngineCalls
	config MockEngineConfig
}

// NewMockEngine creates a new mock engine
func NewMockEngine(config MockEngineConfig) *MockEngine {
	return &MockEngine{
		shapes: make(map[string]types.Dependencies),
		config: config,
		calls:  MockEngineCalls{},
	}
}

// SetSchema stores the application schema
func (m *MockEngine) SetSchema(schema AppSchema) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.config.TrackCalls {
		m.calls.SetSchema = append(m.calls.SetSchema, schema)
	}

	m.schema = &schema
	return nil
}

// ComputeShapeID computes the shape ID for a statement
func (m *MockEngine) ComputeShapeID(stmt types.Statement) (ShapeIDResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.config.TrackCalls {
		m.calls.ComputeShapeID = append(m.calls.ComputeShapeID, stmt)
	}

	shapeID, err := m.computeShapeIDInternal(stmt)
	if err != nil {
		return ShapeIDResponse{}, err
	}

	return ShapeIDResponse{ShapeID: shapeID}, nil
}

// computeShapeIDInternal computes shape ID without locking (internal use)
func (m *MockEngine) computeShapeIDInternal(stmt types.Statement) (string, error) {
	var shapeID string
	if m.config.ShapeIDGenerator != nil {
		shapeID = m.config.ShapeIDGenerator(stmt)
	} else {
		// Use real shapeId computation
		var err error
		shapeID, err = tests.ComputeQueryShapeID(&stmt)
		if err != nil {
			return "", err
		}
	}

	return shapeID, nil
}

// AddQuery adds a query and returns its dependencies
func (m *MockEngine) AddQuery(req AddQueryRequest) (AddQueryResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.config.TrackCalls {
		m.calls.AddQuery = append(m.calls.AddQuery, req)
		// Also track the implicit ComputeShapeID call
		m.calls.ComputeShapeID = append(m.calls.ComputeShapeID, req.Shape)
	}

	// Compute shape ID without locking (we already have the lock)
	shapeID, err := m.computeShapeIDInternal(req.Shape)
	if err != nil {
		return AddQueryResponse{}, err
	}

	deps := types.Dependencies{
		ShapeID:  shapeID,
		Records:  m.extractRecords(req),
		Filters:  m.extractFilters(req.Shape),
		Includes: req.Shape.Includes,
	}

	m.shapes[shapeID] = deps

	return AddQueryResponse{
		ShapeID:      shapeID,
		Dependencies: deps,
	}, nil
}

// Invalidate determines which shapes should be evicted
func (m *MockEngine) Invalidate(mutation types.Mutation) (InvalidateResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.config.TrackCalls {
		m.calls.Invalidate = append(m.calls.Invalidate, mutation)
	}

	// Custom evict list
	if m.config.EvictBehavior == "custom" && len(m.config.CustomEvictList) > 0 {
		return InvalidateResponse{Evict: m.config.CustomEvictList}, nil
	}

	evict := []string{}

	for shapeID, deps := range m.shapes {
		for _, change := range mutation.Changes {
			if m.shouldInvalidate(change, deps) {
				evict = append(evict, shapeID)
				break
			}
		}
	}

	return InvalidateResponse{Evict: evict}, nil
}

// ExplainInvalidation explains why a shape would be invalidated
func (m *MockEngine) ExplainInvalidation(req ExplainRequest) (ExplainResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.config.TrackCalls {
		m.calls.ExplainInvalidation = append(m.calls.ExplainInvalidation, req)
	}

	deps, ok := m.shapes[req.ShapeID]
	if !ok {
		return ExplainResponse{Invalidate: false, Reasons: []string{}}, nil
	}

	reasons := []string{}

	for _, change := range req.Mutation.Changes {
		// Check record membership
		if ids, exists := deps.Records[change.Model]; exists && len(ids) > 0 {
			reasons = append(reasons, "record_membership")
		}

		// Check filter dependencies
		if len(deps.Filters) > 0 {
			for _, filter := range deps.Filters {
				if m.filterReferencesModel(filter, change.Model) {
					reasons = append(reasons, "filter_dependency")
					break
				}
			}
		}

		// Check relation dependencies
		if len(deps.Includes) > 0 {
			for _, include := range deps.Includes {
				if include.Query != nil && include.Query.Model == change.Model {
					reasons = append(reasons, "relation_dependency")
					break
				}
			}
		}
	}

	// Deduplicate reasons
	uniqueReasons := m.deduplicateStrings(reasons)

	return ExplainResponse{
		Invalidate: len(uniqueReasons) > 0,
		Reasons:    uniqueReasons,
	}, nil
}

// Reset clears all engine state
func (m *MockEngine) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.config.TrackCalls {
		m.calls.Reset = append(m.calls.Reset, struct{}{})
	}

	m.schema = nil
	m.shapes = make(map[string]types.Dependencies)

	if m.config.TrackCalls {
		m.calls = MockEngineCalls{}
	}
}

// GetVersion returns version information
func (m *MockEngine) GetVersion() VersionInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.config.TrackCalls {
		m.calls.GetVersion = append(m.calls.GetVersion, struct{}{})
	}

	return VersionInfo{
		Core:     "mock-0.1.0",
		Contract: "0.1.0",
		ABI:      "1",
	}
}

// GetCalls returns all tracked method calls
func (m *MockEngine) GetCalls() MockEngineCalls {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.calls
}

// SetEvictList sets a custom evict list for testing
func (m *MockEngine) SetEvictList(shapeIDs []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config.CustomEvictList = shapeIDs
	m.config.EvictBehavior = "custom"
}

// GetDependencies returns stored dependencies for a shape ID
func (m *MockEngine) GetDependencies(shapeID string) (types.Dependencies, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	deps, ok := m.shapes[shapeID]
	return deps, ok
}

// Helper methods

func (m *MockEngine) extractRecords(req AddQueryRequest) map[string][]string {
	if req.ResultHint == nil || req.Shape.Query == nil {
		return map[string][]string{}
	}

	records := make(map[string][]string)
	model := req.Shape.Query.Model

	if rows, exists := req.ResultHint[model]; exists {
		ids := []string{}
		for _, row := range rows {
			if rowMap, ok := row.(map[string]interface{}); ok {
				if id, ok := rowMap["id"]; ok {
					ids = append(ids, fmt.Sprintf("%v", id))
				}
			}
		}
		if len(ids) > 0 {
			records[model] = ids
		}
	}

	return records
}

func (m *MockEngine) extractFilters(stmt types.Statement) []types.Filter {
	filters := []types.Filter{}

	if stmt.Query != nil && stmt.Query.Where != nil {
		filters = append(filters, *stmt.Query.Where)
	}

	if stmt.Having != nil {
		filters = append(filters, *stmt.Having)
	}

	return filters
}

// filterReferencesModel checks if a filter has any conditions
// Note: This is a simplified implementation for mock/testing purposes.
// A production implementation would parse the filter and check if any
// condition.field references the specified model's fields.
func (m *MockEngine) filterReferencesModel(filter types.Filter, _ string) bool {
	// Simplified: just check if any condition exists
	return filter.Conditions != nil && len(*filter.Conditions) > 0
}

func (m *MockEngine) shouldInvalidate(change types.Change, deps types.Dependencies) bool {
	behavior := m.config.EvictBehavior
	if behavior == "" {
		behavior = "conservative"
	}

	if behavior == "conservative" {
		// Conservative: evict if model is tracked
		_, exists := deps.Records[change.Model]
		return exists
	}

	return false
}

func (m *MockEngine) deduplicateStrings(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range input {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}
