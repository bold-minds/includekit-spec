package tests

import (
	"crypto/sha256"
	"fmt"

	"github.com/bold-minds/ik-spec/go/types"
)

// ComputeShapeID computes shapeId from canonical JSON
func ComputeShapeID(canonicalJSON string) string {
	hash := sha256.Sum256([]byte(canonicalJSON))
	return fmt.Sprintf("s_%x", hash)
}

// ComputeQueryShapeID is a convenience wrapper
func ComputeQueryShapeID(shape *types.QueryShape) (string, error) {
	canonical, err := CanonicalizeQueryShape(shape)
	if err != nil {
		return "", err
	}
	return ComputeShapeID(canonical), nil
}
