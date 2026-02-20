package compiler

import (
	"goMud/internal/gmsl/lexer"
	"goMud/internal/gmsl/parser"
	"strings"
	"testing"
)

func TestCompilerBasic(t *testing.T) {
	input := `
func logic(name string) string {
    var x int
    x = 10
    y := "hello"
    if x == 10 {
        return y
    } else {
        return name
    }
}

func action() {
    room.TryMove("north")
    player.Send(room.GetDescription())
}
`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	ast := p.Parse()

	if ast == nil {
		t.Fatal("Failed to parse AST")
	}

	c := NewCompiler(ast)
	assembly := c.Compile()

	if assembly == nil {
		t.Fatal("Failed to compile to assembly")
	}

	funcs := assembly.GetFunctions()
	if len(funcs) != 2 {
		t.Fatalf("Expected 2 functions, got %d", len(funcs))
	}

	str := assembly.String()
	if !strings.Contains(str, "Function logic:") {
		t.Errorf("Assembly string missing Function logic:\n%s", str)
	}
	if !strings.Contains(str, "Function action:") {
		t.Errorf("Assembly string missing Function action:\n%s", str)
	}
}
