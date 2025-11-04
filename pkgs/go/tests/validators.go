// Package tests provides runtime validators, JCS canonicalization, and shapeId computation
// for the IncludeKit Universal Format.
//
// This is a TESTKIT package - for testing and development only.
// DO NOT import this package in production code. Use the production
// package (github.com/bold-minds/ik-spec/go/types) instead.
package tests

import (
	"fmt"

	"github.com/bold-minds/ik-spec/go/types"
)

// Constants for validation
const (
	ShapeIDPrefix    = "s_"
	ShapeIDLength    = 66 // s_ + 64 hex chars (sha256)
	ShapeIDHexLength = 64
)

// ValidationError represents a validation failure
type ValidationError struct {
	Message string
	Path    string
}

func (e *ValidationError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("%s at %s", e.Message, e.Path)
	}
	return e.Message
}

// ValidateQueryShape validates a QueryShape according to the
// IncludeKit Universal Format specification.
//
// It checks that:
//   - Model is a non-empty string
//   - All filter specifications are valid
//   - OrderBy specifications are valid
//   - Take and Skip are non-negative if present
//   - Distinct fields are non-empty if present
//   - Nested includes are valid
//
// Returns a ValidationError if any constraint is violated.
func ValidateQueryShape(shape *types.QueryShape) error {
	if shape == nil {
		return &ValidationError{Message: "QueryShape cannot be nil", Path: "queryShape"}
	}
	if shape.Model == "" {
		return &ValidationError{Message: "model must be a non-empty string", Path: "queryShape.model"}
	}

	// Validate where clause
	if shape.Where != nil {
		if err := validateFilterSpec(shape.Where, "queryShape.where"); err != nil {
			return err
		}
	}

	// Validate orderBy
	if shape.OrderBy != nil {
		for i, ob := range *shape.OrderBy {
			if err := validateOrderBy(&ob, fmt.Sprintf("queryShape.orderBy[%d]", i)); err != nil {
				return err
			}
		}
	}

	// Validate take (must be non-negative)
	if shape.Take != nil && *shape.Take < 0 {
		return &ValidationError{Message: "take must be non-negative", Path: "queryShape.take"}
	}

	// Validate skip (must be non-negative)
	if shape.Skip != nil && *shape.Skip < 0 {
		return &ValidationError{Message: "skip must be non-negative", Path: "queryShape.skip"}
	}

	// Validate distinct fields
	if shape.Distinct != nil {
		for i, field := range *shape.Distinct {
			if field == "" {
				return &ValidationError{
					Message: "distinct field must be non-empty",
					Path:    fmt.Sprintf("queryShape.distinct[%d]", i),
				}
			}
		}
	}

	// Validate groupBy fields
	if shape.GroupBy != nil {
		for i, field := range *shape.GroupBy {
			if field == "" {
				return &ValidationError{
					Message: "groupBy field must be non-empty",
					Path:    fmt.Sprintf("queryShape.groupBy[%d]", i),
				}
			}
		}
	}

	// Validate having clause
	if shape.Having != nil {
		if err := validateFilterSpec(shape.Having, "queryShape.having"); err != nil {
			return err
		}
	}

	// Validate cursor
	if shape.Cursor != nil {
		if err := validateCursor(shape.Cursor, "queryShape.cursor"); err != nil {
			return err
		}
	}

	// Validate nested includes
	if shape.Include != nil {
		for relation, include := range shape.Include {
			if err := validateInclude(&include, fmt.Sprintf("queryShape.include[%s]", relation)); err != nil {
				return err
			}
		}
	}

	return nil
}

// ValidateMutationEvent validates a MutationEvent
func ValidateMutationEvent(event *types.MutationEvent) error {
	if event == nil {
		return &ValidationError{Message: "MutationEvent cannot be nil", Path: "mutationEvent"}
	}
	if event.Changes == nil {
		return &ValidationError{Message: "changes must be an array", Path: "mutationEvent.changes"}
	}

	for i, change := range event.Changes {
		validOps := map[string]bool{"create": true, "update": true, "delete": true, "link": true, "unlink": true}
		if !validOps[change.Op] {
			return &ValidationError{
				Message: fmt.Sprintf("invalid operation: %s", change.Op),
				Path:    fmt.Sprintf("mutationEvent.changes[%d].op", i),
			}
		}
	}

	return nil
}

// ValidateDependencies validates a Dependencies structure.
//
// It checks that the shapeId follows the correct format (s_ + 64 hex chars)
// and that all required fields are present and valid.
func ValidateDependencies(deps *types.Dependencies) error {
	if deps == nil {
		return &ValidationError{Message: "Dependencies cannot be nil", Path: "dependencies"}
	}
	if deps.ShapeID == "" || len(deps.ShapeID) != ShapeIDLength || deps.ShapeID[:len(ShapeIDPrefix)] != ShapeIDPrefix {
		return &ValidationError{
			Message: fmt.Sprintf("shapeId must match pattern ^%s[0-9a-f]{%d}$", ShapeIDPrefix, ShapeIDHexLength),
			Path:    "dependencies.shapeId",
		}
	}
	if deps.Records == nil {
		return &ValidationError{Message: "records must be an object", Path: "dependencies.records"}
	}
	if deps.FilterBounds == nil {
		return &ValidationError{Message: "filterBounds must be an array", Path: "dependencies.filterBounds"}
	}
	if deps.RelationBounds == nil {
		return &ValidationError{Message: "relationBounds must be an array", Path: "dependencies.relationBounds"}
	}

	return nil
}

func validateFilterSpec(spec *types.FilterSpec, path string) error {
	if spec == nil {
		return nil
	}

	if spec.And != nil {
		for i, s := range *spec.And {
			if err := validateFilterSpec(&s, fmt.Sprintf("%s.and[%d]", path, i)); err != nil {
				return err
			}
		}
	}
	if spec.Or != nil {
		for i, s := range *spec.Or {
			if err := validateFilterSpec(&s, fmt.Sprintf("%s.or[%d]", path, i)); err != nil {
				return err
			}
		}
	}
	if spec.Not != nil {
		if err := validateFilterSpec(spec.Not, fmt.Sprintf("%s.not", path)); err != nil {
			return err
		}
	}
	if spec.Atoms != nil {
		for i, a := range *spec.Atoms {
			if err := validateFilterAtom(&a, fmt.Sprintf("%s.atoms[%d]", path, i)); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateFilterAtom(atom *types.FilterAtom, path string) error {
	if atom.Field == "" {
		return &ValidationError{Message: "field must be a non-empty string", Path: fmt.Sprintf("%s.field", path)}
	}
	if atom.Op == "" {
		return &ValidationError{Message: "op must be a non-empty string", Path: fmt.Sprintf("%s.op", path)}
	}

	validOps := map[string]bool{
		"eq": true, "ne": true, "in": true, "notIn": true, "isNull": true,
		"gt": true, "gte": true, "lt": true, "lte": true, "between": true,
		"contains": true, "startsWith": true, "endsWith": true,
		"like": true, "ilike": true, "regex": true,
		"has": true, "hasSome": true, "hasEvery": true, "jsonContains": true,
		"lenEq": true, "lenGt": true, "lenLt": true, "exists": true,
	}

	isCustomOp := len(atom.Op) >= 7 && atom.Op[:7] == "custom:"
	if !validOps[atom.Op] && !isCustomOp {
		return &ValidationError{Message: fmt.Sprintf("invalid operator: %s", atom.Op), Path: fmt.Sprintf("%s.op", path)}
	}

	return nil
}

func validateOrderBy(ob *types.OrderBySpec, path string) error {
	if ob.Field == "" {
		return &ValidationError{Message: "field must be a non-empty string", Path: fmt.Sprintf("%s.field", path)}
	}
	if ob.Direction != "asc" && ob.Direction != "desc" {
		return &ValidationError{Message: "direction must be 'asc' or 'desc'", Path: fmt.Sprintf("%s.direction", path)}
	}
	if ob.Nulls != nil && *ob.Nulls != "first" && *ob.Nulls != "last" {
		return &ValidationError{Message: "nulls must be 'first' or 'last'", Path: fmt.Sprintf("%s.nulls", path)}
	}
	if ob.Collation != nil && *ob.Collation != "ci" && *ob.Collation != "cs" {
		return &ValidationError{Message: "collation must be 'ci' or 'cs'", Path: fmt.Sprintf("%s.collation", path)}
	}
	return nil
}

func validateCursor(cursor *types.CursorSpec, path string) error {
	if len(cursor.Fields) == 0 {
		return &ValidationError{Message: "cursor fields must not be empty", Path: fmt.Sprintf("%s.fields", path)}
	}
	if len(cursor.Values) != len(cursor.Fields) {
		return &ValidationError{
			Message: fmt.Sprintf("cursor values length (%d) must match fields length (%d)", len(cursor.Values), len(cursor.Fields)),
			Path:    path,
		}
	}
	return nil
}

func validateInclude(include *types.IncludeSpec, path string) error {
	if include.Where != nil {
		if err := validateFilterSpec(include.Where, fmt.Sprintf("%s.where", path)); err != nil {
			return err
		}
	}
	if include.OrderBy != nil {
		for i, ob := range *include.OrderBy {
			if err := validateOrderBy(&ob, fmt.Sprintf("%s.orderBy[%d]", path, i)); err != nil {
				return err
			}
		}
	}
	if include.Take != nil && *include.Take < 0 {
		return &ValidationError{Message: "take must be non-negative", Path: fmt.Sprintf("%s.take", path)}
	}
	if include.Skip != nil && *include.Skip < 0 {
		return &ValidationError{Message: "skip must be non-negative", Path: fmt.Sprintf("%s.skip", path)}
	}
	// Recursively validate nested includes
	if include.Include != nil {
		for relation, nested := range include.Include {
			if err := validateInclude(&nested, fmt.Sprintf("%s.include[%s]", path, relation)); err != nil {
				return err
			}
		}
	}
	return nil
}
