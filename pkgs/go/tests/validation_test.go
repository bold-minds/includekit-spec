package tests_test

import (
	"testing"

	"github.com/bold-minds/ik-spec/go/tests"
	"github.com/bold-minds/ik-spec/go/types"
)

func TestValidateQueryShape_Comprehensive(t *testing.T) {
	tcs := []struct {
		name    string
		shape   *types.QueryShape
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil shape",
			shape:   nil,
			wantErr: true,
			errMsg:  "cannot be nil",
		},
		{
			name: "empty model",
			shape: &types.QueryShape{
				Model: "",
			},
			wantErr: true,
			errMsg:  "model must be",
		},
		{
			name: "negative take",
			shape: &types.QueryShape{
				Model: "Post",
				Take:  intPtr(-1),
			},
			wantErr: true,
			errMsg:  "take must be non-negative",
		},
		{
			name: "negative skip",
			shape: &types.QueryShape{
				Model: "Post",
				Skip:  intPtr(-5),
			},
			wantErr: true,
			errMsg:  "skip must be non-negative",
		},
		{
			name: "invalid order direction",
			shape: &types.QueryShape{
				Model: "Post",
				OrderBy: &[]types.OrderBySpec{
					{Field: "id", Direction: "invalid"},
				},
			},
			wantErr: true,
			errMsg:  "direction must be",
		},
		{
			name: "empty distinct field",
			shape: &types.QueryShape{
				Model:    "Post",
				Distinct: &[]string{""},
			},
			wantErr: true,
			errMsg:  "distinct field must be non-empty",
		},
		{
			name: "empty groupBy field",
			shape: &types.QueryShape{
				Model:   "Post",
				GroupBy: &[]string{"", "category"},
			},
			wantErr: true,
			errMsg:  "groupBy field must be non-empty",
		},
		{
			name: "valid simple shape",
			shape: &types.QueryShape{
				Model: "Post",
				Take:  intPtr(10),
				Skip:  intPtr(0),
			},
			wantErr: false,
		},
		{
			name: "valid with orderBy",
			shape: &types.QueryShape{
				Model: "Post",
				OrderBy: &[]types.OrderBySpec{
					{Field: "createdAt", Direction: "desc"},
					{Field: "id", Direction: "asc"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid with distinct",
			shape: &types.QueryShape{
				Model:    "Post",
				Distinct: &[]string{"authorId", "status"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tcs {
		t.Run(tt.name, func(t *testing.T) {
			err := tests.ValidateQueryShape(tt.shape)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateQueryShape() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateQueryShape() error = %v, want error containing %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestValidateDependencies_Comprehensive(t *testing.T) {
	tcs := []struct {
		name    string
		deps    *types.Dependencies
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil dependencies",
			deps:    nil,
			wantErr: true,
			errMsg:  "cannot be nil",
		},
		{
			name: "invalid shapeId format",
			deps: &types.Dependencies{
				ShapeID:        "invalid",
				Records:        map[string][]string{},
				FilterBounds:   []types.FilterSpec{},
				RelationBounds: []types.RelationFilterBound{},
			},
			wantErr: true,
			errMsg:  "shapeId must match pattern",
		},
		{
			name: "shapeId too short",
			deps: &types.Dependencies{
				ShapeID:        "s_abc",
				Records:        map[string][]string{},
				FilterBounds:   []types.FilterSpec{},
				RelationBounds: []types.RelationFilterBound{},
			},
			wantErr: true,
			errMsg:  "shapeId must match pattern",
		},
		{
			name: "missing records",
			deps: &types.Dependencies{
				ShapeID:        "s_" + string(make([]byte, 64)),
				FilterBounds:   []types.FilterSpec{},
				RelationBounds: []types.RelationFilterBound{},
			},
			wantErr: true,
			errMsg:  "records must be",
		},
		{
			name: "valid dependencies",
			deps: &types.Dependencies{
				ShapeID:        "s_0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				Records:        map[string][]string{"Post": {"1", "2"}},
				FilterBounds:   []types.FilterSpec{},
				RelationBounds: []types.RelationFilterBound{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tcs {
		t.Run(tt.name, func(t *testing.T) {
			err := tests.ValidateDependencies(tt.deps)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDependencies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateDependencies() error = %v, want error containing %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestCanonicalizeQueryShape_Determinism(t *testing.T) {
	shape := &types.QueryShape{
		Model: "Post",
		Where: &types.FilterSpec{
			Atoms: &[]types.FilterAtom{
				{Field: "status", Op: "eq", Value: "published"},
				{Field: "views", Op: "gt", Value: 100},
			},
		},
		OrderBy: &[]types.OrderBySpec{
			{Field: "createdAt", Direction: "desc"},
		},
		Take: intPtr(10),
	}

	// Canonicalize multiple times
	canonical1, err1 := tests.CanonicalizeQueryShape(shape)
	canonical2, err2 := tests.CanonicalizeQueryShape(shape)
	canonical3, err3 := tests.CanonicalizeQueryShape(shape)

	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("Canonicalization failed: %v, %v, %v", err1, err2, err3)
	}

	if canonical1 != canonical2 || canonical2 != canonical3 {
		t.Error("Canonicalization is not deterministic")
	}

	// Verify shapeId is also deterministic
	shapeId1 := tests.ComputeShapeID(canonical1)
	shapeId2 := tests.ComputeShapeID(canonical2)

	if shapeId1 != shapeId2 {
		t.Error("ShapeID computation is not deterministic")
	}

	// Verify format
	if len(shapeId1) != 66 {
		t.Errorf("ShapeID length = %d, want 66", len(shapeId1))
	}
	if shapeId1[:2] != "s_" {
		t.Errorf("ShapeID prefix = %s, want s_", shapeId1[:2])
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
