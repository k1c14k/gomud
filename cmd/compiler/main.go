package main

import (
	"flag"
	"fmt"
	"goMud/internal/gmsl/compiler"
	"goMud/internal/gmsl/lexer"
	"goMud/internal/gmsl/parser"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	source := flag.String("source", "", "Path to the source file")
	flag.Parse()
	b, err := os.ReadFile(*source)
	if err != nil {
		logrus.Error("Error reading file:", err)
		return
	}

	l := lexer.NewLexer(string(b))

	p := parser.NewParser(l)
	ast := p.Parse()
	aout := compiler.NewCompiler(ast).Compile()
	fmt.Println(aout)
}
