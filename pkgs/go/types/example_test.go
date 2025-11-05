package types_test

import (
	"encoding/json"
	"fmt"

	"github.com/bold-minds/includekit-spec/go/types"
)

// Example demonstrates basic QueryShape construction
func Example() {
	shape := &types.Query{
		Model: "Post",
		Where: &types.Filter{
			Conditions: &[]types.Condition{
				{Field: "status", Op: "eq", Value: "published"},
				{Field: "views", Op: "gt", Value: 100},
			},
		},
		OrderBy: &[]types.OrderBy{
			{Field: "createdAt", Descending: boolPtr(true)},
		},
		Limit:  intPtr(10),
		Offset: intPtr(0),
	}

	data, _ := json.MarshalIndent(shape, "", "  ")
	fmt.Println(string(data))
}

// ExampleStatement_withIncludes demonstrates nested relation loading
func ExampleStatement_withIncludes() {
	shape := &types.Statement{
		Query: &types.Query{
			Model: "Post",
		},
		Includes: []types.Include{
			{
				Query: &types.Query{
					Model:  "author",
					Fields: &[]string{"id", "name", "email"},
				},
			},
			{
				Query: &types.Query{
					Model: "comments",
					Where: &types.Filter{
						Conditions: &[]types.Condition{
							{Field: "approved", Op: "eq", Value: true},
						},
					},
					OrderBy: &[]types.OrderBy{
						{Field: "createdAt", Descending: boolPtr(true)},
					},
					Limit: intPtr(5),
				},
			},
		},
	}

	fmt.Printf("Model: %s\n", shape.Query.Model)
	fmt.Printf("Includes: %d relations\n", len(shape.Includes))
	// Output:
	// Model: Post
	// Includes: 2 relations
}

// ExampleFilter demonstrates complex filter construction
func ExampleFilter() {
	filter := &types.Filter{
		Or: &[]types.Filter{
			{
				Conditions: &[]types.Condition{
					{Field: "status", Op: "eq", Value: "published"},
				},
			},
			{
				And: &[]types.Filter{
					{
						Conditions: &[]types.Condition{
							{Field: "status", Op: "eq", Value: "draft"},
							{Field: "authorId", Op: "eq", Value: "123"},
						},
					},
				},
			},
		},
	}

	data, _ := json.Marshal(filter)
	fmt.Println("Complex filter created")
	fmt.Printf("JSON length: %d bytes\n", len(data))
	// Output:
	// Complex filter created
	// JSON length: 191 bytes
}

// ExampleMutation demonstrates write event tracking
func ExampleMutation() {
	event := &types.Mutation{
		TxID: strPtr("tx_abc123"),
		Changes: []types.Change{
			{
				Model:  "posts",
				Action: "insert",
				Sets: []types.KV{
					{Field: "id", Value: "post_1"},
					{Field: "title", Value: "New Post"},
					{Field: "status", Value: "draft"},
				},
			},
			{
				Model:  "posts",
				Action: "update",
				Sets: []types.KV{
					{Field: "status", Value: "published"},
				},
				Where: &types.Filter{
					Conditions: &[]types.Condition{
						{Field: "id", Op: "eq", Value: "post_2"},
					},
				},
			},
			{
				Model:  "posts",
				Action: "delete",
				Where: &types.Filter{
					Conditions: &[]types.Condition{
						{Field: "id", Op: "eq", Value: "post_3"},
					},
				},
			},
		},
	}

	data, _ := json.MarshalIndent(event, "", "  ")
	fmt.Println(string(data))
}

// ExamplePagination demonstrates cursor-based pagination
func ExamplePagination() {
	// Forward pagination (first page)
	firstPage := &types.Statement{
		Query: &types.Query{
			Model: "posts",
			OrderBy: &[]types.OrderBy{
				{Field: "createdAt", Descending: boolPtr(true)},
				{Field: "id"},
			},
		},
		Pagination: &types.Pagination{
			First: intPtr(20), // Get first 20 results
		},
	}

	// Forward pagination (next page)
	// SDK encodes cursor as base64 JSON: {"createdAt": "2024-01-15T10:30:00Z", "id": "post_123"}
	nextPage := &types.Statement{
		Query: &types.Query{
			Model: "posts",
			OrderBy: &[]types.OrderBy{
				{Field: "createdAt", Descending: boolPtr(true)},
				{Field: "id"},
			},
		},
		Pagination: &types.Pagination{
			First: intPtr(20),
			After: strPtr("eyJjcmVhdGVkQXQiOiIyMDI0LTAxLTE1VDEwOjMwOjAwWiIsImlkIjoicG9zdF8xMjMifQ=="),
		},
	}

	// Backward pagination (previous page)
	prevPage := &types.Statement{
		Query: &types.Query{
			Model: "posts",
			OrderBy: &[]types.OrderBy{
				{Field: "createdAt", Descending: boolPtr(true)},
				{Field: "id"},
			},
		},
		Pagination: &types.Pagination{
			Last:   intPtr(20),
			Before: strPtr("eyJjcmVhdGVkQXQiOiIyMDI0LTAxLTE1VDEwOjMwOjAwWiIsImlkIjoicG9zdF8xMjMifQ=="),
		},
	}

	fmt.Printf("First page: first=%d\n", *firstPage.Pagination.First)
	fmt.Printf("Next page: first=%d, after cursor present=%v\n", *nextPage.Pagination.First, nextPage.Pagination.After != nil)
	fmt.Printf("Previous page: last=%d, before cursor present=%v\n", *prevPage.Pagination.Last, prevPage.Pagination.Before != nil)
	// Output:
	// First page: first=20
	// Next page: first=20, after cursor present=true
	// Previous page: last=20, before cursor present=true
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

func strPtr(s string) *string {
	return &s
}
