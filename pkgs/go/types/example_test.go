package types_test

import (
	"encoding/json"
	"fmt"

	"github.com/bold-minds/ik-spec/go/types"
)

// Example demonstrates basic QueryShape construction
func Example() {
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
		Skip: intPtr(0),
	}

	data, _ := json.MarshalIndent(shape, "", "  ")
	fmt.Println(string(data))
}

// ExampleQueryShape_withIncludes demonstrates nested relation loading
func ExampleQueryShape_withIncludes() {
	shape := &types.QueryShape{
		Model: "Post",
		Include: map[string]types.IncludeSpec{
			"author": {
				Select: &[]string{"id", "name", "email"},
			},
			"comments": {
				Where: &types.FilterSpec{
					Atoms: &[]types.FilterAtom{
						{Field: "approved", Op: "eq", Value: true},
					},
				},
				OrderBy: &[]types.OrderBySpec{
					{Field: "createdAt", Direction: "desc"},
				},
				Take: intPtr(5),
			},
		},
	}

	fmt.Printf("Model: %s\n", shape.Model)
	fmt.Printf("Includes: %d relations\n", len(shape.Include))
	// Output:
	// Model: Post
	// Includes: 2 relations
}

// ExampleFilterSpec demonstrates complex filter construction
func ExampleFilterSpec() {
	filter := &types.FilterSpec{
		Or: &[]types.FilterSpec{
			{
				Atoms: &[]types.FilterAtom{
					{Field: "status", Op: "eq", Value: "published"},
				},
			},
			{
				And: &[]types.FilterSpec{
					{
						Atoms: &[]types.FilterAtom{
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
	// JSON length: 181 bytes
}

// ExampleMutationEvent demonstrates write event tracking
func ExampleMutationEvent() {
	event := &types.MutationEvent{
		TxID: strPtr("tx_abc123"),
		Changes: []types.WriteChange{
			{
				Op:    "create",
				Model: "Post",
				ID:    "post_1",
				After: map[string]interface{}{
					"title":  "New Post",
					"status": "draft",
				},
			},
			{
				Op:    "update",
				Model: "Post",
				ID:    "post_2",
				Before: map[string]interface{}{
					"status": "draft",
				},
				After: map[string]interface{}{
					"status": "published",
				},
			},
		},
	}

	fmt.Printf("Transaction: %s\n", *event.TxID)
	fmt.Printf("Changes: %d operations\n", len(event.Changes))
	// Output:
	// Transaction: tx_abc123
	// Changes: 2 operations
}

func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}
