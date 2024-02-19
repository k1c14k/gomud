package vm

import (
	"bytes"
	"goMud/internal/gmsl/compiler"
	"goMud/internal/gmsl/lexer"
	"goMud/internal/gmsl/parser"
	"log"
	"os"
	"strconv"
)

type Class struct {
	name    string
	methods map[string]Method
}

func (c *Class) GetMethod(name string) Method {
	return c.methods[name]
}

func newClass(name string) *Class {
	log.Println("Loading class", name)

	b, err := os.ReadFile("mudlib/" + name + ".gms")
	if err != nil {
		log.Panicln("Error reading file:", err)
	}

	l := lexer.NewLexer(string(b))
	p := parser.NewParser(l)
	ast := p.Parse()
	aOut := compiler.NewCompiler(ast).Compile()

	return &Class{name: name, methods: NewMethodsFromAssembly(aOut)}
}

func NewEmptyClass(name string) *Class {
	return &Class{name: name, methods: make(map[string]Method)}
}

func (c *Class) RegisterInternalMethod(name string, argumentCount int, returnValueCount int, handle MethodHandler) {
	c.methods[name] = &internalMethod{argumentCount, returnValueCount, handle}

}

func (c *Class) String() string {
	buff := bytes.NewBufferString("Class[")
	buff.WriteString(c.name)
	buff.WriteString(", methods=")
	buff.WriteString(strconv.Itoa(len(c.methods)))
	buff.WriteString("]")
	return buff.String()
}
