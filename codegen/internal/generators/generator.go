package generators

import (
	"fmt"

	"github.com/bold-minds/ik-spec/codegen/internal/parser"
)

// Generator defines the interface for language-specific code generators
type Generator interface {
	Generate(s *parser.Schema, outputDir string) error
	Language() string
	NeedsExternal() bool // Does it need external tools like npm?
}

// Get returns a generator for the specified language
func Get(lang string) Generator {
	switch lang {
	case "typescript", "ts":
		return &TypeScriptGenerator{}
	case "go", "golang":
		return &GoGenerator{}
	case "java":
		return &JavaGenerator{}
	case "dotnet", "csharp", "c#":
		return &DotNetGenerator{}
	case "python", "py":
		return &PythonGenerator{}
	case "php":
		return &PHPGenerator{}
	default:
		return nil
	}
}

// Placeholder generators for future languages

type JavaGenerator struct{}

func (g *JavaGenerator) Generate(s *parser.Schema, outputDir string) error {
	return fmt.Errorf("java generator not yet implemented")
}

func (g *JavaGenerator) Language() string { return "Java" }

func (g *JavaGenerator) NeedsExternal() bool { return false }

type DotNetGenerator struct{}

func (g *DotNetGenerator) Generate(s *parser.Schema, outputDir string) error {
	return fmt.Errorf(".NET generator not yet implemented")
}

func (g *DotNetGenerator) Language() string { return ".NET/C#" }

func (g *DotNetGenerator) NeedsExternal() bool { return false }

type PythonGenerator struct{}

func (g *PythonGenerator) Generate(s *parser.Schema, outputDir string) error {
	return fmt.Errorf("python generator not yet implemented")
}

func (g *PythonGenerator) Language() string { return "Python" }

func (g *PythonGenerator) NeedsExternal() bool { return false }

type PHPGenerator struct{}

func (g *PHPGenerator) Generate(s *parser.Schema, outputDir string) error {
	return fmt.Errorf("php generator not yet implemented")
}

func (g *PHPGenerator) Language() string { return "PHP" }

func (g *PHPGenerator) NeedsExternal() bool { return false }
