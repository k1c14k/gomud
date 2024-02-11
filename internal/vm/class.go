package vm

import (
	"goMud/internal/gmsl/compiler"
	"goMud/internal/gmsl/lexer"
	"goMud/internal/gmsl/parser"
	"log"
	"os"
)

type Class interface {
	GetStringPool() []string
	GetMethod(name string) Method
}

type vmClass struct {
	name       string
	stringPool []string
	methods    map[string]Method
}

func (c *vmClass) GetStringPool() []string {
	return c.stringPool
}

func (c *vmClass) GetMethod(name string) Method {
	return c.methods[name]
}

func NewClass(name string) Class {
	log.Println("Loading class", name)

	b, err := os.ReadFile("mudlib/" + name + ".gms")
	if err != nil {
		log.Panicln("Error reading file:", err)
	}

	l := lexer.NewLexer(string(b))
	p := parser.NewParser(l)
	ast := p.Parse()
	aOut := compiler.NewCompiler(ast).Compile()

	return &vmClass{name: name, stringPool: aOut.Consts, methods: NewMethodsFromAssembly(aOut)}
}
