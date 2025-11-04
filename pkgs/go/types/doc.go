// Package types provides production type definitions for the IncludeKit Universal Format.
//
// This is a PRODUCTION package containing only type definitions with no runtime utilities.
// For validation, canonicalization, and shape ID computation, use the testkit package:
// github.com/bold-minds/ik-spec/go/tests
//
// # Overview
//
// The IncludeKit Universal Format defines a cross-ORM, language-agnostic way to express
// database queries, mutations, and dependencies. It enables:
//
//   - Deterministic canonical JSON (JCS RFC 8785) representation
//   - SHA-256 based shape hashing for cache keys
//   - Conservative invalidation when query results might change
//
// # Core Types
//
// QueryShape represents a normalized read operation:
//
//	shape := &types.QueryShape{
//	    Model: "Post",
//	    Where: &types.FilterSpec{...},
//	    OrderBy: &[]types.OrderBySpec{...},
//	    Take: intPtr(10),
//	}
//
// MutationEvent describes write operations that may invalidate cached reads:
//
//	event := &types.MutationEvent{
//	    Changes: []types.WriteChange{
//	        {Op: "create", Model: "Post", ID: "1", After: {...}},
//	    },
//	}
//
// Dependencies tracks what a read depends on (engine output):
//
//	deps := &types.Dependencies{
//	    ShapeID: "s_abc...",
//	    Records: map[string][]string{"Post": {"1", "2"}},
//	    FilterBounds: []types.FilterSpec{...},
//	}
//
// # Implementation Boundary
//
// Production code should ONLY import this package for type definitions.
// Runtime utilities (validators, JCS, hashing) belong in separate testkit packages:
//
//   - TypeScript: @includekit/types-testkit
//   - Go: github.com/bold-minds/ik-spec/go/tests
//
// This separation ensures production bundles remain lightweight and type-focused.
//
// # Schema Definition
//
// Types are generated from the JSON Schema at:
// https://github.com/bold-minds/ik-spec/schema/v0-1-0.json
//
// For the full specification, see:
// https://github.com/bold-minds/ik-spec/schema/README.md
package types
