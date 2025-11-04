package generators

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bold-minds/ik-spec/codegen/internal/parser"
)

type GoGenerator struct{}

func (g *GoGenerator) Generate(s *parser.Schema, outputDir string) error {
	// NOTE: Go types are currently HAND-WRITTEN to preserve idiomatic patterns:
	// - Sealed Scalar interface (vs plain interface{})
	// - Pointer-to-slice for optional arrays (nil vs empty distinction)
	// - Custom package documentation
	//
	// See GO_GENERATION_ANALYSIS.md for details on why we don't auto-generate.
	//
	// To enable auto-generation, uncomment the following block and update testkit code:
	/*
	if err := g.generateTypes(s, outputDir); err != nil {
		return fmt.Errorf("failed to generate types: %w", err)
	}
	*/

	// For now, just verify the packages exist
	typesDir := filepath.Join(outputDir, "go", "types")
	if _, err := os.Stat(typesDir); os.IsNotExist(err) {
		return fmt.Errorf("go types directory does not exist: %s", typesDir)
	}

	testsDir := filepath.Join(outputDir, "go", "tests")
	if _, err := os.Stat(testsDir); os.IsNotExist(err) {
		return fmt.Errorf("go tests directory does not exist: %s", testsDir)
	}

	return nil
}

func (g *GoGenerator) generateTypes(s *parser.Schema, outputDir string) error {
	// Validate and clean output directory
	outputDir = filepath.Clean(outputDir)
	if strings.Contains(outputDir, "..") {
		return fmt.Errorf("invalid output directory (directory traversal detected): %s", outputDir)
	}

	typesDir := filepath.Join(outputDir, "go", "types")
	if err := os.MkdirAll(typesDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outputFile := filepath.Join(typesDir, "types.go")

	// Check if go-jsonschema is available
	goJsonSchemaPath, err := exec.LookPath("go-jsonschema")
	if err != nil {
		// Try in GOPATH/bin
		cmd := exec.Command("go", "env", "GOPATH")
		gopathBytes, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("go-jsonschema not found and could not determine GOPATH: %w. Install with: go install github.com/atombender/go-jsonschema@latest", err)
		}
		gopath := strings.TrimSpace(string(gopathBytes))
		if gopath == "" {
			return fmt.Errorf("go-jsonschema not found in PATH and GOPATH is empty. Install with: go install github.com/atombender/go-jsonschema@latest")
		}
		goJsonSchemaPath = filepath.Join(gopath, "bin", "go-jsonschema")
		if _, err := os.Stat(goJsonSchemaPath); err != nil {
			return fmt.Errorf("go-jsonschema not found at %s. Install with: go install github.com/atombender/go-jsonschema@latest", goJsonSchemaPath)
		}
	}

	// Validate schema path before using
	schemaPath := filepath.Clean(s.Path)
	if !filepath.IsAbs(schemaPath) {
		var err error
		schemaPath, err = filepath.Abs(schemaPath)
		if err != nil {
			return fmt.Errorf("failed to resolve schema path: %w", err)
		}
	}
	if strings.Contains(schemaPath, "..") {
		return fmt.Errorf("invalid schema path (directory traversal detected): %s", schemaPath)
	}

	// Call go-jsonschema
	cmd := exec.Command(goJsonSchemaPath,
		"-p", "types",
		"--only-models",
		"--capitalization", "ID",
		"--tags", "json",
		"-o", outputFile,
		schemaPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go-jsonschema failed: %w\nOutput: %s", err, output)
	}

	return nil
}

func (g *GoGenerator) Language() string {
	return "Go"
}

func (g *GoGenerator) NeedsExternal() bool {
	return true // Needs go-jsonschema
}
