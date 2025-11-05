package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bold-minds/includekit-spec/codegen/internal/generators"
	"github.com/bold-minds/includekit-spec/codegen/internal/parser"
)

func main() {
	languages := flag.String("lang", "all", "Languages to generate (all,ts,go,java,dotnet,python,php)")
	schemaPath := flag.String("schema", "schema/v0-1-0.json", "Path to JSON Schema")
	outputDir := flag.String("output", "pkgs", "Output directory")
	verbose := flag.Bool("v", false, "Verbose output")

	flag.Parse()

	fmt.Println("📦 Generating code from schema...")

	// Parse schema
	s, err := parser.Parse(*schemaPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to parse schema: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Parsed schema: %s (v%s)\n", s.Title, s.Version)
	}

	// Determine which languages to generate
	langs := parseLangs(*languages)

	// Generate for each language
	for _, lang := range langs {
		gen := generators.Get(lang)
		if gen == nil {
			fmt.Fprintf(os.Stderr, "❌ Unknown language: %s\n", lang)
			os.Exit(1)
		}

		fmt.Printf("Generating %s...\n", gen.Language())

		if err := gen.Generate(s, *outputDir); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to generate %s: %v\n", lang, err)
			os.Exit(1)
		}

		fmt.Printf("✓ Generated %s\n", gen.Language())
	}

	fmt.Println("✓ Code generation complete!")
}

func parseLangs(input string) []string {
	if input == "all" {
		return []string{"typescript", "go"}
	}
	return strings.Split(input, ",")
}
