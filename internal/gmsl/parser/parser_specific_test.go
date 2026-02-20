package parser

import (
	"goMud/internal/gmsl/lexer"
	"strings"
	"testing"
)

func TestParserSpecificNodes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string // Strings that should be in the String() output
	}{
		{
			name: "single import declaration",
			input: `import "math"
func test() {}`,
			contains: []string{`(import math)`},
		},
		{
			name: "import declaration list",
			input: `import (
				"math"
				"fmt"
			)
			func test() {}`,
			contains: []string{`(import math fmt)`},
		},
		{
			name:     "method call expression",
			input:    `func test() { player.move("north") }`,
			contains: []string{`(method-call player move (string "north"))`},
		},
		{
			name: "if else statements",
			input: `func test() { 
				if player.health() == 0 { return 1 } else { player.move("up") }
			}`,
			// Note: String() output might be nested, checking parts
			contains: []string{`(if (== (method-call player health) 0)`, `(return 1)`, `(else (method-call player move (string "up")))`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := NewParser(l)
			ast := p.Parse()

			if ast == nil {
				t.Fatal("Expected AST, got nil")
			}

			out := ast.String()

			for _, expectedStr := range tt.contains {
				if !strings.Contains(out, expectedStr) {
					t.Errorf("AST output missing expected string.\nExpected to contain:\n%s\nGot AST:\n%s", expectedStr, out)
				}
			}
		})
	}
}
