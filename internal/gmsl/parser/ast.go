package parser

import (
	"bytes"
	"goMud/internal/gmsl"
)

type AstNode interface {
	GetToken() *gmsl.Token
	String() string
}

type Identifier struct {
	AstNode
	token *gmsl.Token
	Value string
}

type ImportDeclaration interface {
	AstNode
}

type SingleImportDeclaration struct {
	AstNode
	ImportDeclaration
	token *gmsl.Token
	Name  *Identifier
}

type ImportDeclarationList struct {
	AstNode
	ImportDeclaration
	token   *gmsl.Token
	Imports []*Identifier
}

type Type struct {
	AstNode
	token *gmsl.Token
	Name  string
}

type ArgumentDeclaration struct {
	AstNode
	token *gmsl.Token
	Name  *Identifier
	Typ   *Type
}

type Statement interface {
	AstNode
}

type Expression interface {
	AstNode
}

type BinaryExpression struct {
	AstNode
	Expression
	token *gmsl.Token
	Left  Expression
	Right Expression
}

type StringLiteralExpression struct {
	AstNode
	Expression
	token *gmsl.Token
	Value string
}

type MethodCallExpression struct {
	AstNode
	Expression
	token      *gmsl.Token
	ObjectName *Identifier
	MethodName *Identifier
	Arguments  []Expression
}

type ExpressionStatement struct {
	AstNode
	Expression
	token           *gmsl.Token
	ExpressionValue Expression
}

type FunctionDeclaration struct {
	AstNode
	token      *gmsl.Token
	Name       *Identifier
	Arguments  []ArgumentDeclaration
	Statements []Statement
}

type Class struct {
	AstNode
	token     *gmsl.Token
	Name      *Identifier
	Imports   []ImportDeclaration
	Functions []FunctionDeclaration
}

func (c *Class) GetToken() *gmsl.Token {
	return c.token
}

func (c *Class) String() string {
	var buf bytes.Buffer
	buf.WriteString("(class ")
	buf.WriteString(c.Name.String())
	for _, i := range c.Imports {
		buf.WriteString(" ")
		buf.WriteString(i.String())
	}
	for _, f := range c.Functions {
		buf.WriteString(" ")
		buf.WriteString(f.String())
	}
	buf.WriteString(")")
	return buf.String()
}

func (s *SingleImportDeclaration) GetToken() *gmsl.Token {
	return s.token
}

func (s *SingleImportDeclaration) String() string {
	var buf bytes.Buffer
	buf.WriteString("(import ")
	buf.WriteString(s.Name.String())
	buf.WriteString(")")
	return buf.String()
}

func (i *ImportDeclarationList) GetToken() *gmsl.Token {
	return i.token
}

func (i *ImportDeclarationList) String() string {
	var buf bytes.Buffer
	buf.WriteString("(import")
	for _, i := range i.Imports {
		buf.WriteString(" ")
		buf.WriteString(i.String())
	}
	buf.WriteString(")")
	return buf.String()
}

func (i *Identifier) GetToken() *gmsl.Token {
	return i.token
}

func (i *Identifier) String() string {
	return i.Value
}

func (b *BinaryExpression) GetToken() *gmsl.Token {
	return b.token
}

func (b *BinaryExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString("(")
	buf.WriteString(b.token.Value)
	buf.WriteString(" ")
	buf.WriteString(b.Left.String())
	buf.WriteString(" ")
	buf.WriteString(b.Right.String())
	buf.WriteString(")")
	return buf.String()
}

func (m *MethodCallExpression) GetToken() *gmsl.Token {
	return m.token
}

func (m *MethodCallExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString("(method-call ")
	buf.WriteString(m.ObjectName.String())
	buf.WriteString(" ")
	buf.WriteString(m.MethodName.String())
	for _, a := range m.Arguments {
		buf.WriteString(" ")
		buf.WriteString(a.String())
	}
	buf.WriteString(")")
	return buf.String()
}

func (s *StringLiteralExpression) GetToken() *gmsl.Token {
	return s.token
}

func (s *StringLiteralExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString("(string \"")
	buf.WriteString(s.Value)
	buf.WriteString("\")")
	return buf.String()
}

func (f *FunctionDeclaration) GetToken() *gmsl.Token {
	return f.token
}

func (f *FunctionDeclaration) String() string {
	var buf bytes.Buffer
	buf.WriteString("(func ")
	buf.WriteString(f.Name.String())
	for _, a := range f.Arguments {
		buf.WriteString(" ")
		buf.WriteString(a.String())
	}
	for _, s := range f.Statements {
		buf.WriteString(" ")
		buf.WriteString(s.String())
	}
	buf.WriteString(")")
	return buf.String()
}

func (a *ArgumentDeclaration) GetToken() *gmsl.Token {
	return a.token
}

func (a *ArgumentDeclaration) String() string {
	var buf bytes.Buffer
	buf.WriteString("(arg ")
	buf.WriteString(a.Name.String())
	buf.WriteString(" ")
	buf.WriteString(a.Typ.String())
	buf.WriteString(")")
	return buf.String()
}

func (e *ExpressionStatement) GetToken() *gmsl.Token {
	return e.token
}

func (e *ExpressionStatement) String() string {
	return e.ExpressionValue.String()
}

func (t *Type) GetToken() *gmsl.Token {
	return t.token
}

func (t *Type) String() string {
	return t.Name
}
