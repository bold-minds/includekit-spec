package mock_test

import (
	"testing"

	"github.com/bold-minds/includekit-spec/go/tests/mock"
	"github.com/bold-minds/includekit-spec/go/types"
)

func TestSetSchema(t *testing.T) {
	engine := mock.NewMockEngine(mock.MockEngineConfig{TrackCalls: true})
	schema := mock.AppSchema{
		Version: 1,
		Models: []mock.Model{
			{Name: "users", ID: mock.IDConfig{Kind: "string"}},
			{Name: "posts", ID: mock.IDConfig{Kind: "string"}},
		},
	}

	err := engine.SetSchema(schema)
	if err != nil {
		t.Fatalf("SetSchema failed: %v", err)
	}

	calls := engine.GetCalls()
	if len(calls.SetSchema) != 1 {
		t.Errorf("Expected 1 SetSchema call, got %d", len(calls.SetSchema))
	}
}

func TestComputeShapeID(t *testing.T) {
	engine := mock.NewMockEngine(mock.MockEngineConfig{})
	
	stmt := types.Statement{
		Query: &types.Query{
			Model: "users",
			Where: &types.Filter{
				Conditions: &[]types.Condition{
					{Field: "id", Op: "eq", Value: "1"},
				},
			},
		},
	}

	result1, err := engine.ComputeShapeID(stmt)
	if err != nil {
		t.Fatalf("ComputeShapeID failed: %v", err)
	}

	result2, err := engine.ComputeShapeID(stmt)
	if err != nil {
		t.Fatalf("ComputeShapeID failed: %v", err)
	}

	if result1.ShapeID != result2.ShapeID {
		t.Error("ShapeID should be deterministic")
	}

	if result1.ShapeID[:2] != "s_" {
		t.Error("ShapeID should start with s_")
	}

	if len(result1.ShapeID) != 66 {
		t.Errorf("ShapeID should be 66 chars, got %d", len(result1.ShapeID))
	}
}

func TestComputeShapeIDCustomGenerator(t *testing.T) {
	engine := mock.NewMockEngine(mock.MockEngineConfig{
		ShapeIDGenerator: func(stmt types.Statement) string {
			return "s_custom_test"
		},
	})

	result, err := engine.ComputeShapeID(types.Statement{
		Query: &types.Query{Model: "users"},
	})
	if err != nil {
		t.Fatalf("ComputeShapeID failed: %v", err)
	}

	if result.ShapeID != "s_custom_test" {
		t.Errorf("Expected s_custom_test, got %s", result.ShapeID)
	}
}

func TestAddQuery(t *testing.T) {
	engine := mock.NewMockEngine(mock.MockEngineConfig{})

	stmt := types.Statement{
		Query: &types.Query{
			Model: "users",
			Where: &types.Filter{
				Conditions: &[]types.Condition{
					{Field: "id", Op: "eq", Value: "1"},
				},
			},
		},
	}

	result, err := engine.AddQuery(mock.AddQueryRequest{Shape: stmt})
	if err != nil {
		t.Fatalf("AddQuery failed: %v", err)
	}

	if result.ShapeID == "" {
		t.Error("ShapeID should not be empty")
	}

	if result.Dependencies.ShapeID != result.ShapeID {
		t.Error("Dependencies.ShapeID should match result.ShapeID")
	}

	if len(result.Dependencies.Records) != 0 {
		t.Error("Records should be empty without result_hint")
	}

	if len(result.Dependencies.Filters) != 1 {
		t.Errorf("Expected 1 filter, got %d", len(result.Dependencies.Filters))
	}
}

func TestAddQueryExtractsRecords(t *testing.T) {
	engine := mock.NewMockEngine(mock.MockEngineConfig{})

	stmt := types.Statement{
		Query: &types.Query{Model: "users"},
	}

	result, err := engine.AddQuery(mock.AddQueryRequest{
		Shape: stmt,
		ResultHint: map[string][]interface{}{
			"users": {
				map[string]interface{}{"id": "1", "name": "Alice"},
				map[string]interface{}{"id": "2", "name": "Bob"},
			},
		},
	})
	if err != nil {
		t.Fatalf("AddQuery failed: %v", err)
	}

	userIDs, exists := result.Dependencies.Records["users"]
	if !exists {
		t.Fatal("Expected users in records")
	}

	if len(userIDs) != 2 {
		t.Errorf("Expected 2 user IDs, got %d", len(userIDs))
	}

	if userIDs[0] != "1" || userIDs[1] != "2" {
		t.Errorf("Expected IDs [1, 2], got %v", userIDs)
	}
}

func TestInvalidateEvictsAffectedShapes(t *testing.T) {
	engine := mock.NewMockEngine(mock.MockEngineConfig{})

	stmt := types.Statement{
		Query: &types.Query{Model: "users"},
	}

	addResult, err := engine.AddQuery(mock.AddQueryRequest{
		Shape: stmt,
		ResultHint: map[string][]interface{}{
			"users": {
				map[string]interface{}{"id": "1", "name": "Alice"},
			},
		},
	})
	if err != nil {
		t.Fatalf("AddQuery failed: %v", err)
	}

	mutation := types.Mutation{
		Changes: []types.Change{
			{
				Model:  "users",
				Action: "update",
				Sets:   []types.KV{{Field: "name", Value: "Alice Updated"}},
				Where: &types.Filter{
					Conditions: &[]types.Condition{
						{Field: "id", Op: "eq", Value: "1"},
					},
				},
			},
		},
	}

	invalidateResult, err := engine.Invalidate(mutation)
	if err != nil {
		t.Fatalf("Invalidate failed: %v", err)
	}

	found := false
	for _, shapeID := range invalidateResult.Evict {
		if shapeID == addResult.ShapeID {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected shape to be evicted")
	}
}

func TestInvalidateDoesNotEvictUnrelatedShapes(t *testing.T) {
	engine := mock.NewMockEngine(mock.MockEngineConfig{})

	stmt := types.Statement{
		Query: &types.Query{Model: "posts"},
	}

	addResult, err := engine.AddQuery(mock.AddQueryRequest{Shape: stmt})
	if err != nil {
		t.Fatalf("AddQuery failed: %v", err)
	}

	mutation := types.Mutation{
		Changes: []types.Change{
			{
				Model:  "users",
				Action: "update",
				Sets:   []types.KV{{Field: "name", Value: "Test"}},
				Where:  &types.Filter{},
			},
		},
	}

	invalidateResult, err := engine.Invalidate(mutation)
	if err != nil {
		t.Fatalf("Invalidate failed: %v", err)
	}

	for _, shapeID := range invalidateResult.Evict {
		if shapeID == addResult.ShapeID {
			t.Error("Shape should not be evicted")
		}
	}
}

func TestInvalidateCustomEvictList(t *testing.T) {
	engine := mock.NewMockEngine(mock.MockEngineConfig{})
	engine.SetEvictList([]string{"s_custom_1", "s_custom_2"})

	mutation := types.Mutation{
		Changes: []types.Change{
			{Model: "users", Action: "update", Sets: []types.KV{}},
		},
	}

	result, err := engine.Invalidate(mutation)
	if err != nil {
		t.Fatalf("Invalidate failed: %v", err)
	}

	if len(result.Evict) != 2 {
		t.Errorf("Expected 2 evicted shapes, got %d", len(result.Evict))
	}

	if result.Evict[0] != "s_custom_1" || result.Evict[1] != "s_custom_2" {
		t.Errorf("Unexpected evict list: %v", result.Evict)
	}
}

func TestExplainInvalidation(t *testing.T) {
	engine := mock.NewMockEngine(mock.MockEngineConfig{})

	stmt := types.Statement{
		Query: &types.Query{Model: "users"},
	}

	addResult, err := engine.AddQuery(mock.AddQueryRequest{
		Shape: stmt,
		ResultHint: map[string][]interface{}{
			"users": {map[string]interface{}{"id": "1"}},
		},
	})
	if err != nil {
		t.Fatalf("AddQuery failed: %v", err)
	}

	mutation := types.Mutation{
		Changes: []types.Change{
			{Model: "users", Action: "update", Sets: []types.KV{}},
		},
	}

	result, err := engine.ExplainInvalidation(mock.ExplainRequest{
		Mutation: mutation,
		ShapeID:  addResult.ShapeID,
	})
	if err != nil {
		t.Fatalf("ExplainInvalidation failed: %v", err)
	}

	if !result.Invalidate {
		t.Error("Expected invalidate to be true")
	}

	found := false
	for _, reason := range result.Reasons {
		if reason == "record_membership" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected record_membership reason, got %v", result.Reasons)
	}
}

func TestExplainInvalidationUnknownShape(t *testing.T) {
	engine := mock.NewMockEngine(mock.MockEngineConfig{})

	mutation := types.Mutation{
		Changes: []types.Change{
			{Model: "users", Action: "update", Sets: []types.KV{}},
		},
	}

	result, err := engine.ExplainInvalidation(mock.ExplainRequest{
		Mutation: mutation,
		ShapeID:  "s_unknown",
	})
	if err != nil {
		t.Fatalf("ExplainInvalidation failed: %v", err)
	}

	if result.Invalidate {
		t.Error("Expected invalidate to be false for unknown shape")
	}

	if len(result.Reasons) != 0 {
		t.Errorf("Expected empty reasons, got %v", result.Reasons)
	}
}

func TestReset(t *testing.T) {
	engine := mock.NewMockEngine(mock.MockEngineConfig{TrackCalls: true})

	engine.AddQuery(mock.AddQueryRequest{
		Shape: types.Statement{Query: &types.Query{Model: "users"}},
	})

	calls := engine.GetCalls()
	if len(calls.AddQuery) != 1 {
		t.Error("Expected 1 AddQuery call before reset")
	}

	engine.Reset()

	calls = engine.GetCalls()
	if len(calls.AddQuery) != 0 {
		t.Error("Expected 0 AddQuery calls after reset")
	}
}

func TestGetVersion(t *testing.T) {
	engine := mock.NewMockEngine(mock.MockEngineConfig{})

	version := engine.GetVersion()

	if version.Core != "mock-0.1.0" {
		t.Errorf("Expected core mock-0.1.0, got %s", version.Core)
	}

	if version.Contract != "0.1.0" {
		t.Errorf("Expected contract 0.1.0, got %s", version.Contract)
	}

	if version.ABI != "1" {
		t.Errorf("Expected ABI 1, got %s", version.ABI)
	}
}

func TestTrackCalls(t *testing.T) {
	engine := mock.NewMockEngine(mock.MockEngineConfig{TrackCalls: true})

	engine.SetSchema(mock.AppSchema{Version: 1, Models: []mock.Model{}})
	engine.ComputeShapeID(types.Statement{Query: &types.Query{Model: "users"}})
	engine.AddQuery(mock.AddQueryRequest{
		Shape: types.Statement{Query: &types.Query{Model: "posts"}},
	})
	engine.GetVersion()

	calls := engine.GetCalls()

	if len(calls.SetSchema) != 1 {
		t.Errorf("Expected 1 SetSchema call, got %d", len(calls.SetSchema))
	}

	if len(calls.ComputeShapeID) != 2 {
		t.Errorf("Expected 2 ComputeShapeID calls, got %d", len(calls.ComputeShapeID))
	}

	if len(calls.AddQuery) != 1 {
		t.Errorf("Expected 1 AddQuery call, got %d", len(calls.AddQuery))
	}

	if len(calls.GetVersion) != 1 {
		t.Errorf("Expected 1 GetVersion call, got %d", len(calls.GetVersion))
	}
}
