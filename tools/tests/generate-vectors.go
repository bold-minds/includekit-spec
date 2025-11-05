// Package main generates test vectors for the new schema format
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type TestVector struct {
	Name              string      `json:"name"`
	Shape             interface{} `json:"shape"`
	ExpectedCanonical string      `json:"expectedCanonical"`
	ExpectedShapeID   string      `json:"expectedShapeId"`
}

func main() {
	vectors := []TestVector{
		{
			Name: "minimal-query",
			Shape: map[string]interface{}{
				"query": map[string]interface{}{
					"model": "Post",
				},
			},
		},
		{
			Name: "simple-query-with-filter",
			Shape: map[string]interface{}{
				"query": map[string]interface{}{
					"model": "Post",
					"where": map[string]interface{}{
						"conditions": []map[string]interface{}{
							{"field": "published", "op": "eq", "value": true},
						},
					},
				},
			},
		},
		{
			Name: "with-order-and-limit",
			Shape: map[string]interface{}{
				"query": map[string]interface{}{
					"model": "Post",
					"order_by": []map[string]interface{}{
						{"field": "createdAt", "descending": true},
					},
					"limit": 10,
				},
			},
		},
		{
			Name: "with-fields-and-distinct",
			Shape: map[string]interface{}{
				"query": map[string]interface{}{
					"model":    "Post",
					"fields":   []string{"id", "title"},
					"distinct": []string{"authorId"},
				},
			},
		},
		{
			Name: "with-includes",
			Shape: map[string]interface{}{
				"query": map[string]interface{}{
					"model": "User",
				},
				"includes": []map[string]interface{}{
					{
						"query": map[string]interface{}{
							"model":  "posts",
							"fields": []string{"id", "title"},
						},
					},
				},
			},
		},
		{
			Name: "with-nested-includes",
			Shape: map[string]interface{}{
				"query": map[string]interface{}{
					"model": "User",
				},
				"includes": []map[string]interface{}{
					{
						"query": map[string]interface{}{
							"model": "posts",
						},
						"includes": []map[string]interface{}{
							{
								"query": map[string]interface{}{
									"model": "comments",
									"limit": 5,
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "with-relation-filter",
			Shape: map[string]interface{}{
				"query": map[string]interface{}{
					"model": "User",
				},
				"includes": []map[string]interface{}{
					{
						"kind": "some",
						"query": map[string]interface{}{
							"model": "posts",
							"where": map[string]interface{}{
								"conditions": []map[string]interface{}{
									{"field": "published", "op": "eq", "value": true},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "with-pagination",
			Shape: map[string]interface{}{
				"query": map[string]interface{}{
					"model": "Post",
					"order_by": []map[string]interface{}{
						{"field": "createdAt", "descending": true},
						{"field": "id"},
					},
				},
				"pagination": map[string]interface{}{
					"first": 20,
					"after": "eyJpZCI6InBvc3RfMTIzIn0=",
				},
			},
		},
		{
			Name: "complex-filter",
			Shape: map[string]interface{}{
				"query": map[string]interface{}{
					"model": "Post",
					"where": map[string]interface{}{
						"and": []map[string]interface{}{
							{
								"conditions": []map[string]interface{}{
									{"field": "published", "op": "eq", "value": true},
								},
							},
							{
								"or": []map[string]interface{}{
									{
										"conditions": []map[string]interface{}{
											{"field": "featured", "op": "eq", "value": true},
										},
									},
									{
										"conditions": []map[string]interface{}{
											{"field": "views", "op": "gte", "value": 100},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "with-group-by-having",
			Shape: map[string]interface{}{
				"query": map[string]interface{}{
					"model":  "Post",
					"fields": []string{"authorId", "COUNT(*) as count"},
				},
				"group_by": []string{"authorId"},
				"having": map[string]interface{}{
					"conditions": []map[string]interface{}{
						{"field": "count", "op": "gt", "value": 5},
					},
				},
			},
		},
	}

	// Compute canonical JSON and shape IDs
	for i := range vectors {
		canonical, err := canonicalize(vectors[i].Shape)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error canonicalizing %s: %v\n", vectors[i].Name, err)
			os.Exit(1)
		}
		vectors[i].ExpectedCanonical = canonical
		vectors[i].ExpectedShapeID = computeShapeID(canonical)
	}

	// Write to file
	outputPath := filepath.Join("tools", "tests", "vectors", "query-shapes.json")
	data, err := json.MarshalIndent(vectors, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling vectors: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Generated %d test vectors in %s\n", len(vectors), outputPath)
}

// canonicalize produces JCS (RFC 8785) canonical JSON
func canonicalize(v interface{}) (string, error) {
	// Marshal to JSON first
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	// Unmarshal to generic interface
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return "", err
	}

	// Canonicalize
	canonical := canonicalizeValue(obj)

	// Marshal back to canonical JSON
	result, err := json.Marshal(canonical)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func canonicalizeValue(val interface{}) interface{} {
	switch val := val.(type) {
	case map[string]interface{}:
		// Sort keys
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		result := make(map[string]interface{})
		for _, k := range keys {
			result[k] = canonicalizeValue(val[k])
		}
		return result

	case []interface{}:
		result := make([]interface{}, len(val))
		for i, item := range val {
			result[i] = canonicalizeValue(item)
		}
		return result

	default:
		return val
	}
}

func computeShapeID(canonical string) string {
	hash := sha256.Sum256([]byte(canonical))
	return "s_" + hex.EncodeToString(hash[:])
}
