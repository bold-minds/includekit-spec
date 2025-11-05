// Package tests provides cross-language conformance tests for
// IncludeKit Universal Format implementations.
package tests_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bold-minds/includekit-spec/go/tests"
	"github.com/bold-minds/includekit-spec/go/types"
)

type Vector struct {
	Name              string          `json:"name"`
	Shape             types.Statement `json:"shape"`
	ExpectedCanonical string          `json:"expectedCanonical"`
	ExpectedShapeID   string          `json:"expectedShapeId"`
}

func TestConformanceQueryShapes(t *testing.T) {
	// Load shared test vectors from tools directory
	vectorsPath := filepath.Join("..", "..", "..", "tools", "tests", "vectors", "query-shapes.json")
	data, err := os.ReadFile(vectorsPath)
	if err != nil {
		t.Fatalf("Failed to read vectors: %v", err)
	}

	var vectors []Vector
	if err := json.Unmarshal(data, &vectors); err != nil {
		t.Fatalf("Failed to parse vectors: %v", err)
	}

	for _, v := range vectors {
		t.Run(v.Name, func(t *testing.T) {
			// Validate
			if err := tests.ValidateQueryShape(&v.Shape); err != nil {
				t.Errorf("Validation failed: %v", err)
			}

			// Canonicalize
			canonical, err := tests.CanonicalizeQueryShape(&v.Shape)
			if err != nil {
				t.Errorf("Canonicalization failed: %v", err)
			}

			// CRITICAL: Compare against expected canonical JSON
			if canonical != v.ExpectedCanonical {
				t.Errorf("Canonical JSON mismatch for %s:\n  got:  %s\n  want: %s",
					v.Name, canonical, v.ExpectedCanonical)
			}

			// Compute shapeId
			shapeID := tests.ComputeShapeID(canonical)

			// CRITICAL: Compare against expected shapeId
			if shapeID != v.ExpectedShapeID {
				t.Errorf("ShapeID mismatch for %s:\n  got:  %s\n  want: %s",
					v.Name, shapeID, v.ExpectedShapeID)
			}

			// Basic format checks
			if !strings.HasPrefix(shapeID, "s_") {
				t.Errorf("ShapeID should start with s_, got: %s", shapeID)
			}
			if len(shapeID) != 66 {
				t.Errorf("ShapeID should be 66 chars, got: %d", len(shapeID))
			}

			// Verify determinism
			canonical2, _ := tests.CanonicalizeQueryShape(&v.Shape)
			if canonical != canonical2 {
				t.Error("Canonicalization should be deterministic")
			}

			shapeID2 := tests.ComputeShapeID(canonical2)
			if shapeID != shapeID2 {
				t.Error("ShapeID should be deterministic")
			}
		})
	}
}

func TestValidationRejectsInvalidShapes(t *testing.T) {
	invalidShape := &types.Statement{
		Query: &types.Query{
			Model: "", // empty model
		},
	}

	err := tests.ValidateQueryShape(invalidShape)
	if err == nil {
		t.Error("Should reject empty model")
	}
}
