// Package main synchronizes version across all files from VERSION file (SSOT)
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// Read version from VERSION file (single source of truth)
	versionBytes, err := os.ReadFile("VERSION")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading VERSION file: %v\n", err)
		os.Exit(1)
	}
	version := strings.TrimSpace(string(versionBytes))

	// Validate semver format
	if !regexp.MustCompile(`^\d+\.\d+\.\d+$`).MatchString(version) {
		fmt.Fprintf(os.Stderr, "Invalid version format: %s (expected: X.Y.Z)\n", version)
		os.Exit(1)
	}

	versionDashed := strings.ReplaceAll(version, ".", "-")
	versionMajorMinor := version[:strings.LastIndex(version, ".")]

	fmt.Printf("📦 Syncing version %s across all files...\n", version)

	// 1. Update schema file (rename and update contents)
	if err := syncSchema(version, versionDashed, versionMajorMinor); err != nil {
		fmt.Fprintf(os.Stderr, "Error syncing schema: %v\n", err)
		os.Exit(1)
	}

	// 2. Update package.json files
	if err := syncPackageJSON("pkgs/ts/types/package.json", version); err != nil {
		fmt.Fprintf(os.Stderr, "Error syncing types package.json: %v\n", err)
		os.Exit(1)
	}

	if err := syncPackageJSON("pkgs/ts/tests/package.json", version); err != nil {
		fmt.Fprintf(os.Stderr, "Error syncing testkit package.json: %v\n", err)
		os.Exit(1)
	}

	// 3. Update codegen default
	if err := updateFile("codegen/main.go",
		regexp.MustCompile(`schema/v\d+-\d+-\d+\.json`),
		fmt.Sprintf("schema/v%s.json", versionDashed)); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating codegen: %v\n", err)
		os.Exit(1)
	}

	// 4. Update CI workflow
	if err := updateFile(".github/workflows/ci.yml",
		regexp.MustCompile(`schema/v\d+-\d+-\d+\.json`),
		fmt.Sprintf("schema/v%s.json", versionDashed)); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating CI workflow: %v\n", err)
		os.Exit(1)
	}

	// 5. Update release workflow
	if err := updateFile(".github/workflows/release.yml",
		regexp.MustCompile(`schema/v\d+-\d+-\d+\.json`),
		fmt.Sprintf("schema/v%s.json", versionDashed)); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating release workflow: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Version sync complete!")
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  1. Run: go run tools/version/sync.go\n")
	fmt.Printf("  2. Regenerate code: cd codegen && go run .\n")
	fmt.Printf("  3. Run tests: ./scripts/test.sh\n")
	fmt.Printf("  4. Commit: git add -A && git commit -m 'chore: bump version to v%s'\n", version)
}

func syncSchema(version, versionDashed, versionMajorMinor string) error {
	oldPattern := "schema/v*-*-*.json"
	matches, err := filepath.Glob(oldPattern)
	if err != nil {
		return err
	}

	newPath := fmt.Sprintf("schema/v%s.json", versionDashed)

	// Read existing schema
	var schemaPath string
	if len(matches) > 0 {
		schemaPath = matches[0]
	} else {
		return fmt.Errorf("no schema file found matching %s", oldPattern)
	}

	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return err
	}

	var schema map[string]interface{}
	if err := json.Unmarshal(data, &schema); err != nil {
		return err
	}

	// Update schema metadata
	schema["$id"] = fmt.Sprintf("https://github.com/bold-minds/ik-spec/schema/v%s.json", versionDashed)
	schema["title"] = fmt.Sprintf("IncludeKit Universal Format v%s", versionMajorMinor)

	// Write updated schema
	updatedData, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(newPath, updatedData, 0644); err != nil {
		return err
	}

	// Remove old file if different
	if schemaPath != newPath {
		if err := os.Remove(schemaPath); err != nil {
			return err
		}
		fmt.Printf("  ✓ Renamed %s → %s\n", schemaPath, newPath)
	}

	fmt.Printf("  ✓ Updated schema to v%s\n", version)
	return nil
}

func syncPackageJSON(path, version string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return err
	}

	pkg["version"] = version

	updatedData, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return err
	}

	// Add newline at end
	updatedData = append(updatedData, '\n')

	if err := os.WriteFile(path, updatedData, 0644); err != nil {
		return err
	}

	fmt.Printf("  ✓ Updated %s to v%s\n", path, version)
	return nil
}

func updateFile(path string, pattern *regexp.Regexp, replacement string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	updated := pattern.ReplaceAllString(string(data), replacement)

	if err := os.WriteFile(path, []byte(updated), 0644); err != nil {
		return err
	}

	fmt.Printf("  ✓ Updated %s\n", path)
	return nil
}
