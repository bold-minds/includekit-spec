// Package parser provides JSON Schema parsing functionality for the
// IncludeKit Universal Format code generator.
//
// It reads and validates JSON Schema files, extracting metadata and
// definitions needed for language-specific code generation.
package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Schema represents a parsed JSON Schema from the schema/ directory
type Schema struct {
	Path        string
	Schema      string                 `json:"$schema"`
	ID          string                 `json:"$id"`
	Title       string                 `json:"title"`
	Type        string                 `json:"type"`
	Version     string                 // Extracted from title or filename
	Definitions map[string]interface{} `json:"$defs"`
	Properties  map[string]interface{} `json:"properties"`
	Raw         map[string]interface{} // Full raw schema
}

// Parse reads and parses a JSON Schema file from the given path.
// It validates the file exists, parses the JSON structure, and extracts
// version information from the schema title or filename.
//
// Returns an error if the file doesn't exist, isn't valid JSON,
// or if the path contains directory traversal attempts.
func Parse(path string) (*Schema, error) {
	// Validate path safety
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return nil, fmt.Errorf("invalid schema path (directory traversal detected): %s", path)
	}

	// Check file exists with better error message
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("schema file not found: %s", cleanPath)
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	var s Schema
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	// Parse raw for full access
	if err := json.Unmarshal(data, &s.Raw); err != nil {
		return nil, err
	}

	s.Path = cleanPath
	s.Version = extractVersion(s.Title, cleanPath)

	return &s, nil
}

// extractVersion extracts the version string from the schema title or filename.
// It tries the title first (e.g., "IncludeKit Universal Format v1.0"),
// then falls back to the filename pattern (e.g., "v1-0-0.json").
//
// Returns "unknown" if no version can be extracted.
func extractVersion(title, path string) string {
	// Try to extract from title: "IncludeKit Universal Format v1.0" or "v1.2.3"
	if title != "" {
		re := regexp.MustCompile(`v(\d+\.\d+(?:\.\d+)?)`)
		if matches := re.FindStringSubmatch(title); len(matches) > 1 {
			return matches[1]
		}
	}

	// Fallback to filename pattern: "v1-0-0.json" -> "1.0.0"
	basename := filepath.Base(path)
	re := regexp.MustCompile(`v(\d+)-(\d+)-(\d+)`)
	if matches := re.FindStringSubmatch(basename); len(matches) > 3 {
		return fmt.Sprintf("%s.%s.%s", matches[1], matches[2], matches[3])
	}

	return "unknown"
}
