// Package tests provides runtime validators, JCS canonicalization, and shapeId computation
// for the IncludeKit Universal Format.
//
// This is a TESTKIT package - for testing and development only.
// DO NOT import this package in production code. Use the production
// package (github.com/bold-minds/includekit-spec/go/types) instead.
package tests

import (
	"fmt"

	"github.com/bold-minds/includekit-spec/go/types"
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

// ValidateQueryShape validates a Statement structure.
//
// It checks that:
//   - Query is present with non-empty model
//   - All filters, orderBy specs, pagination are valid
//   - Limit and offset are non-negative
//   - Distinct and groupBy fields are non-empty strings
//   - Nested includes are valid
//
// Returns a ValidationError if any constraint is violated.
func ValidateQueryShape(stmt *types.Statement) error {
	if stmt == nil {
		return &ValidationError{Message: "Statement cannot be nil", Path: "statement"}
	}

	// Validate query
	if stmt.Query != nil {
		if err := validateQuery(stmt.Query, "statement.query"); err != nil {
			return err
		}
	}

	// Validate groupBy fields
	if stmt.GroupBy != nil {
		for i, field := range *stmt.GroupBy {
			if field == "" {
				return &ValidationError{
					Message: "groupBy field must be non-empty",
					Path:    fmt.Sprintf("statement.groupBy[%d]", i),
				}
			}
		}
	}

	// Validate having clause
	if stmt.Having != nil {
		if err := validateFilterSpec(stmt.Having, "statement.having"); err != nil {
			return err
		}
	}

	// Validate pagination
	if stmt.Pagination != nil {
		if err := validatePagination(stmt.Pagination, "statement.pagination"); err != nil {
			return err
		}
	}

	// Validate includes
	if stmt.Includes != nil {
		for i, include := range stmt.Includes {
			if err := validateInclude(&include, fmt.Sprintf("statement.includes[%d]", i)); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateQuery(q *types.Query, path string) error {
	if q.Model == "" {
		return &ValidationError{Message: "model must be a non-empty string", Path: fmt.Sprintf("%s.model", path)}
	}

	// Validate where clause
	if q.Where != nil {
		if err := validateFilterSpec(q.Where, fmt.Sprintf("%s.where", path)); err != nil {
			return err
		}
	}

	// Validate orderBy
	if q.OrderBy != nil {
		for i, ob := range *q.OrderBy {
			if err := validateOrderBy(&ob, fmt.Sprintf("%s.orderBy[%d]", path, i)); err != nil {
				return err
			}
		}
	}

	// Validate limit (must be non-negative)
	if q.Limit != nil && *q.Limit < 0 {
		return &ValidationError{Message: "limit must be non-negative", Path: fmt.Sprintf("%s.limit", path)}
	}

	// Validate offset (must be non-negative)
	if q.Offset != nil && *q.Offset < 0 {
		return &ValidationError{Message: "offset must be non-negative", Path: fmt.Sprintf("%s.offset", path)}
	}

	// Validate distinct fields
	if q.Distinct != nil {
		for i, field := range *q.Distinct {
			if field == "" {
				return &ValidationError{
					Message: "distinct field must be non-empty",
					Path:    fmt.Sprintf("%s.distinct[%d]", path, i),
				}
			}
		}
	}

	return nil
}

// ValidateMutationEvent validates a Mutation
func ValidateMutationEvent(event *types.Mutation) error {
	if event == nil {
		return &ValidationError{Message: "Mutation cannot be nil", Path: "mutation"}
	}
	if event.Changes == nil {
		return &ValidationError{Message: "changes must be an array", Path: "mutation.changes"}
	}

	for i, change := range event.Changes {
		if err := validateDataChange(&change, fmt.Sprintf("mutation.changes[%d]", i)); err != nil {
			return err
		}
	}

	return nil
}

func validateDataChange(change *types.Change, path string) error {
	// Validate model
	if change.Model == "" {
		return &ValidationError{Message: "model must be non-empty", Path: fmt.Sprintf("%s.model", path)}
	}

	// Validate action
	validActions := map[string]bool{"insert": true, "update": true, "delete": true}
	if !validActions[change.Action] {
		return &ValidationError{
			Message: fmt.Sprintf("action must be 'insert', 'update', or 'delete', got: %s", change.Action),
			Path:    fmt.Sprintf("%s.action", path),
		}
	}

	// Validate based on action type
	switch change.Action {
	case "insert":
		// Insert requires Set, no Where
		if len(change.Sets) == 0 {
			return &ValidationError{
				Message: "insert requires non-empty set",
				Path:    fmt.Sprintf("%s.set", path),
			}
		}
		if change.Where != nil {
			return &ValidationError{
				Message: "insert cannot have where clause",
				Path:    fmt.Sprintf("%s.where", path),
			}
		}

	case "update":
		// Update requires both Set and Where
		if len(change.Sets) == 0 {
			return &ValidationError{
				Message: "update requires non-empty set",
				Path:    fmt.Sprintf("%s.set", path),
			}
		}
		if change.Where == nil {
			return &ValidationError{
				Message: "update requires where clause",
				Path:    fmt.Sprintf("%s.where", path),
			}
		}

	case "delete":
		// Delete requires Where, no Set
		if len(change.Sets) > 0 {
			return &ValidationError{
				Message: "delete cannot have set clause",
				Path:    fmt.Sprintf("%s.set", path),
			}
		}
		if change.Where == nil {
			return &ValidationError{
				Message: "delete requires where clause",
				Path:    fmt.Sprintf("%s.where", path),
			}
		}
	}

	// Validate Set clauses
	for j, setClause := range change.Sets {
		if setClause.Field == "" {
			return &ValidationError{
				Message: "set clause field must be non-empty",
				Path:    fmt.Sprintf("%s.set[%d].field", path, j),
			}
		}
	}

	// Validate Where clause if present
	if change.Where != nil {
		if err := validateFilterSpec(change.Where, fmt.Sprintf("%s.where", path)); err != nil {
			return err
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
	if deps.Filters == nil {
		return &ValidationError{Message: "filterBounds must be an array", Path: "dependencies.filterBounds"}
	}
	if deps.Includes == nil {
		return &ValidationError{Message: "relationBounds must be an array", Path: "dependencies.relationBounds"}
	}

	return nil
}

func validateFilterSpec(spec *types.Filter, path string) error {
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
	if spec.Conditions != nil {
		for i, a := range *spec.Conditions {
			if err := validateFilterAtom(&a, fmt.Sprintf("%s.atoms[%d]", path, i)); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateFilterAtom(atom *types.Condition, path string) error {
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

func validateOrderBy(ob *types.OrderBy, path string) error {
	if ob.Field == "" {
		return &ValidationError{Message: "field must be a non-empty string", Path: fmt.Sprintf("%s.field", path)}
	}
	// Descending, NullsFirst and CaseSensitive are bools - no validation needed
	return nil
}

func validatePagination(p *types.Pagination, path string) error {
	// Can't mix forward and backward pagination
	hasForward := p.First != nil || p.After != nil
	hasBackward := p.Last != nil || p.Before != nil

	if hasForward && hasBackward {
		return &ValidationError{
			Message: "cannot mix forward pagination (first/after) with backward pagination (last/before)",
			Path:    path,
		}
	}

	// Validate First (must be positive)
	if p.First != nil && *p.First <= 0 {
		return &ValidationError{
			Message: "first must be a positive integer",
			Path:    fmt.Sprintf("%s.first", path),
		}
	}

	// Validate Last (must be positive)
	if p.Last != nil && *p.Last <= 0 {
		return &ValidationError{
			Message: "last must be a positive integer",
			Path:    fmt.Sprintf("%s.last", path),
		}
	}

	// After/Before are opaque strings, no validation needed
	// (SDKs encode them as base64 JSON)

	return nil
}

func validateInclude(include *types.Include, path string) error {
	// Validate query if present
	if include.Query != nil {
		if err := validateQuery(include.Query, fmt.Sprintf("%s.query", path)); err != nil {
			return err
		}
	}

	// Validate kind if present
	if include.Kind != nil {
		validKinds := map[string]bool{"some": true, "every": true, "none": true}
		if !validKinds[*include.Kind] {
			return &ValidationError{
				Message: "kind must be 'some', 'every', or 'none'",
				Path:    fmt.Sprintf("%s.kind", path),
			}
		}
	}

	// Recursively validate nested includes
	if include.Includes != nil {
		for i, nested := range include.Includes {
			if err := validateInclude(&nested, fmt.Sprintf("%s.includes[%d]", path, i)); err != nil {
				return err
			}
		}
	}
	return nil
}
