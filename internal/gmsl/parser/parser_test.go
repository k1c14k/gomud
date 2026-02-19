package parser

import (
	"goMud/internal/gmsl/lexer"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParserRegression(t *testing.T) {
	matches, err := filepath.Glob("../../../mudlib/**/*.gms")
	if err != nil {
		t.Fatalf("Failed to glob mudlib files: %v", err)
	}
	rootMatches, err := filepath.Glob("../../../mudlib/*.gms")
	if err != nil {
		t.Fatalf("Failed to glob root mudlib files: %v", err)
	}
	matches = append(matches, rootMatches...)

	if len(matches) == 0 {
		t.Fatal("Found no .gms files for testing")
	}

	for _, match := range matches {
		name := filepath.Base(match)
		t.Run(name, func(t *testing.T) {
			b, err := os.ReadFile(match)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}

			l := lexer.NewLexer(string(b))
			p := NewParser(l)
			ast := p.Parse()

			if ast == nil {
				t.Fatal("Expected AST, got nil")
			}

			out := ast.PrettyPrint(0)

			// Check if expected file exists
			expectedPath := filepath.Join("testdata", name+".expected")

			// Auto-generate if it doesn't exist (baseline creation)
			if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
				os.MkdirAll("testdata", 0755)
				err := os.WriteFile(expectedPath, []byte(out), 0644)
				if err != nil {
					t.Fatalf("Failed to write expected file: %v", err)
				}
				t.Logf("Generated baseline for %s", name)
			} else {
				// Compare with expected
				expectedBytes, err := os.ReadFile(expectedPath)
				if err != nil {
					t.Fatalf("Failed to read expected file: %v", err)
				}
				expectedStr := string(expectedBytes)
				if strings.TrimSpace(out) != strings.TrimSpace(expectedStr) {
					t.Errorf("AST output mismatch.\nExpected:\n%s\nGot:\n%s", expectedStr, out)
				}
			}
		})
	}
}
