package compiler

import (
	"bytes"
	"goMud/internal/gmsl/parser"
	"log"
	"strconv"
)

type RegisterType int

const (
	StringRegister RegisterType = iota
)

type RegisterReference struct {
	Typ   RegisterType
	Index int
}

type FunctionContext struct {
	variableToRegister map[string]RegisterReference
}

type Compiler struct {
	ast             *parser.AstNode
	strings         []string
	assemblyEntries []AssemblyEntry
	methodPositions map[string]int
}

func NewCompiler(ast parser.AstNode) *Compiler {
	return &Compiler{
		&ast,
		make([]string, 0),
		make([]AssemblyEntry, 0),
		make(map[string]int),
	}
}

type Assembly struct {
	Consts  []string
	Entries []AssemblyEntry
}

func (a *Assembly) String() string {
	var b bytes.Buffer
	b.WriteString("Data:\n")
	for n, c := range a.Consts {
		b.WriteString("[")
		b.WriteString(strconv.Itoa(n))
		b.WriteString("] ")
		b.WriteString(c)
		b.WriteString("\n")
	}
	b.WriteString("Program:\n")
	for _, e := range a.Entries {
		b.WriteString(e.String())
		b.WriteString("\n")
	}
	return b.String()
}

func (c *Compiler) Compile() *Assembly {
	c.processNode(c.ast)
	return &Assembly{c.strings, c.assemblyEntries}
}

func (c *Compiler) processNode(node *parser.AstNode) {
	switch n := (*node).(type) {
	case *parser.Class:
		c.processClass(n)
	case *parser.ImportDeclarationList:
		c.processImportDeclarationList(n)
	case *parser.SingleImportDeclaration:
		c.processSingleImportDeclaration(n)
	case *parser.FunctionDeclaration:
		c.processFunctionDeclaration(n)
	default:
		log.Panicln("Unknown node type", n.String())
	}
}

func (c *Compiler) processClass(n *parser.Class) {
	c.strings = append(c.strings, n.Name.Value)
	c.assemblyEntries = append(c.assemblyEntries, NewLabelEntry(".class_name", len(c.strings)-1, *n.GetToken()))

	for _, i := range n.Imports {
		a := i.(parser.AstNode)
		c.processNode(&a)
	}

	for _, f := range n.Functions {
		var a parser.AstNode = &f
		c.processNode(&a)
	}
}

func (c *Compiler) processImportDeclarationList(n *parser.ImportDeclarationList) {
	for _, i := range n.Imports {
		c.strings = append(c.strings, i.Value)
		c.assemblyEntries = append(c.assemblyEntries, NewLabelEntry(".import_name", len(c.strings)-1, *i.GetToken()))
	}
}

func (c *Compiler) processSingleImportDeclaration(n *parser.SingleImportDeclaration) {
	c.strings = append(c.strings, n.Name.Value)
	c.assemblyEntries = append(c.assemblyEntries, NewLabelEntry(".import_name", len(c.strings)-1, *n.GetToken()))
}

func (c *Compiler) processFunctionDeclaration(n *parser.FunctionDeclaration) {
	c.strings = append(c.strings, n.Name.Value)
	c.assemblyEntries = append(c.assemblyEntries, NewLabelEntry(".function_name", len(c.strings)-1, *n.GetToken()))
	c.methodPositions[n.Name.Value] = len(c.assemblyEntries) - 1

	ctx := NewFunctionContext()
	for _, a := range n.Arguments {
		c.processArgumentDeclaration(&a, ctx)
	}
	for _, s := range n.Statements {
		c.processStatement(&s, ctx)
	}
	c.assemblyEntries = append(c.assemblyEntries, NewReturnEntry(*n.GetToken()))
}

func NewFunctionContext() *FunctionContext {
	return &FunctionContext{
		variableToRegister: make(map[string]RegisterReference),
	}
}

var typToRTyp = map[string]RegisterType{
	"string": StringRegister,
}

func (c *Compiler) processArgumentDeclaration(argumentDeclaration *parser.ArgumentDeclaration, context *FunctionContext) {
	typ := argumentDeclaration.Typ
	rType := typToRTyp[typ.Name]
	index := 0
	for _, r := range context.variableToRegister {
		if r.Typ == rType {
			index = max(index, r.Index)
		}
	}
	index++
	r := RegisterReference{rType, index}
	context.variableToRegister[argumentDeclaration.Name.Value] = r
	c.assemblyEntries = append(c.assemblyEntries, NewPopToRegisterEntry(r, *argumentDeclaration.GetToken()))
}

func (c *Compiler) processStatement(s *parser.Statement, ctx *FunctionContext) {
	switch n := (*s).(type) {
	case *parser.ExpressionStatement:
		c.processExpressionStatement(n, ctx)
	default:
		log.Panicln("Unknown statement type", n.String())
	}
}

func (c *Compiler) processExpressionStatement(statement *parser.ExpressionStatement, ctx *FunctionContext) {
	c.assemblyEntries = append(c.assemblyEntries, c.processExpression(&statement.ExpressionValue, ctx)...)
}

func isContextName(name string) bool {
	return name == "player" || name == "room" || name == "item"
}

func (c *Compiler) processExpression(expression *parser.Expression, ctx *FunctionContext) []AssemblyEntry {
	var result []AssemblyEntry
	switch (*expression).(type) {
	case *parser.MethodCallExpression:
		for _, a := range (*expression).(*parser.MethodCallExpression).Arguments {
			result = append(result, c.processExpression(&a, ctx)...)
		}
		methodName := (*expression).(*parser.MethodCallExpression).MethodName
		c.strings = append(c.strings, methodName.Value)
		result = append(result, NewLabelEntry(".method_name", len(c.strings)-1, *methodName.GetToken()))
		objectName := (*expression).(*parser.MethodCallExpression).ObjectName
		c.strings = append(c.strings, objectName.Value)
		if isContextName(objectName.Value) {
			result = append(result, NewPushContextEntry(len(c.strings)-1, *objectName.GetToken()))
		} else {
			result = append(result, NewLabelEntry(".object_name", len(c.strings)-1, *objectName.GetToken()))
		}
		n := (*expression).(*parser.MethodCallExpression).GetToken()
		result = append(result, NewMethodCallEntry(*n))
	case *parser.BinaryExpression:
		result = append(result, c.processExpression(&(*expression).(*parser.BinaryExpression).Left, ctx)...)
		result = append(result, c.processExpression(&(*expression).(*parser.BinaryExpression).Right, ctx)...)
		result = append(result, NewOperationEntry(*(*expression).(*parser.BinaryExpression).GetToken()))
	case *parser.StringLiteralExpression:
		c.strings = append(c.strings, (*expression).(*parser.StringLiteralExpression).Value)
		result = append(result, NewLabelEntry(".string", len(c.strings)-1, *(*expression).(*parser.StringLiteralExpression).GetToken()))
	default:
		log.Panicln("Unknown expression type", (*expression).String())
	}

	return result
}
