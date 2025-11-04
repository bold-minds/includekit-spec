package parser

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (string, func())
		wantErr     bool
		wantVersion string
	}{
		{
			name: "valid schema file",
			setup: func() (string, func()) {
				tmpfile, _ := os.CreateTemp("", "schema-*.json")
				content := `{
					"$schema": "http://json-schema.org/draft-07/schema#",
					"title": "Test Schema v1.2.3",
					"type": "object",
					"$defs": {}
				}`
				tmpfile.Write([]byte(content))
				tmpfile.Close()
				return tmpfile.Name(), func() { os.Remove(tmpfile.Name()) }
			},
			wantErr:     false,
			wantVersion: "1.2.3",
		},
		{
			name: "file not found",
			setup: func() (string, func()) {
				return "/nonexistent/schema.json", func() {}
			},
			wantErr: true,
		},
		{
			name: "directory traversal attempt",
			setup: func() (string, func()) {
				return "../../../etc/passwd", func() {}
			},
			wantErr: true,
		},
		{
			name: "invalid JSON",
			setup: func() (string, func()) {
				tmpfile, _ := os.CreateTemp("", "schema-*.json")
				tmpfile.Write([]byte("invalid json"))
				tmpfile.Close()
				return tmpfile.Name(), func() { os.Remove(tmpfile.Name()) }
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, cleanup := tt.setup()
			defer cleanup()

			schema, err := Parse(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && schema.Version != tt.wantVersion {
				t.Errorf("Parse() version = %v, want %v", schema.Version, tt.wantVersion)
			}
		})
	}
}

func TestExtractVersion(t *testing.T) {
	tests := []struct {
		name  string
		title string
		path  string
		want  string
	}{
		{
			name:  "from title with v prefix",
			title: "IncludeKit Universal Format v1.0",
			path:  "schema.json",
			want:  "1.0",
		},
		{
			name:  "from title with full version",
			title: "Test Schema v2.1.3",
			path:  "schema.json",
			want:  "2.1.3",
		},
		{
			name:  "from filename",
			title: "",
			path:  "/path/to/v1-0-0.json",
			want:  "1.0.0",
		},
		{
			name:  "from filename with dashes",
			title: "",
			path:  "v2-5-1.json",
			want:  "2.5.1",
		},
		{
			name:  "no version found",
			title: "Schema without version",
			path:  "schema.json",
			want:  "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractVersion(tt.title, tt.path)
			if got != tt.want {
				t.Errorf("extractVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
