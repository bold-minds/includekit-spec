// Package tests provides runtime validation, canonicalization,
// and shape ID computation for the IncludeKit Universal Format.
package tests

import (
	"encoding/json"
	"sort"

	"github.com/bold-minds/ik-spec/go/types"
)

// Canonicalize returns the JCS (RFC 8785) canonical JSON string representation
// of the given object.
//
// It recursively sorts all object keys in lexicographic order and marshals
// the result to JSON. This ensures deterministic output for hashing.
//
// Returns an error if the object cannot be marshaled to JSON.
func Canonicalize(obj interface{}) (string, error) {
	if obj == nil {
		return "null", nil
	}
	normalized := canonicalizeValue(obj)
	data, err := json.Marshal(normalized)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// CanonicalizeQueryShape removes diagnostic fields and canonicalizes
func CanonicalizeQueryShape(shape *types.QueryShape) (string, error) {
	// Make a copy and remove diagnostic fields
	data, err := json.Marshal(shape)
	if err != nil {
		return "", err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return "", err
	}

	delete(m, "orm")
	delete(m, "adapterVersion")

	return Canonicalize(m)
}

func canonicalizeValue(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case map[string]interface{}:
		// Handle nil map
		if val == nil {
			return nil
		}
		// Sort keys lexicographically
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		sorted := make(map[string]interface{}, len(val))
		for _, k := range keys {
			sorted[k] = canonicalizeValue(val[k])
		}
		return sorted

	case []interface{}:
		// Handle nil slice
		if val == nil {
			return nil
		}
		// Recursively canonicalize array elements
		result := make([]interface{}, len(val))
		for i, elem := range val {
			result[i] = canonicalizeValue(elem)
		}
		return result

	default:
		return v
	}
}
