package main

import (
	"fmt"
	"goMud/internal/gmsl"
	"goMud/internal/gmsl/compiler"
	"goMud/internal/gmsl/parser"
	"os"
)

func main() {
	// read mudlib/player_handler.gms into string
	b, err := os.ReadFile("mudlib/player_handler.gms")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	l := gmsl.NewLexer(string(b))

	p := parser.NewParser(l)
	ast := p.Parse()
	aout := compiler.NewCompiler(ast).Compile()
	fmt.Println(aout)
}
