package parser

import (
	"bytes"
	"goMud/internal/gmsl/lexer"
)

type AstNode interface {
	GetToken() *lexer.Token
	String() string
	PrettyPrint(tabs int) string
}

type Identifier struct {
	token *lexer.Token
	Value string
}

type ImportDeclaration interface {
	AstNode
}

type SingleImportDeclaration struct {
	token *lexer.Token
	Name  Identifier
}

type ImportDeclarationList struct {
	token   *lexer.Token
	Imports []Identifier
}

type Type struct {
	token *lexer.Token
	Name  string
}

type ArgumentDeclaration struct {
	token *lexer.Token
	Name  Identifier
	Typ   Type
}

type Statement interface {
	AstNode
}

type Expression interface {
	AstNode
}

type BinaryExpression struct {
	token *lexer.Token
	Left  Expression
	Right Expression
}

type StringLiteralExpression struct {
	token *lexer.Token
	Value string
}

type MethodCallExpression struct {
	token      *lexer.Token
	ObjectName Identifier
	MethodName Identifier
	Arguments  []Expression
}

type ExpressionStatement struct {
	token           *lexer.Token
	ExpressionValue Expression
}

type FunctionDeclaration struct {
	token      *lexer.Token
	Name       Identifier
	Arguments  []ArgumentDeclaration
	Statements []Statement
}

type Class struct {
	token     *lexer.Token
	Name      Identifier
	Imports   []ImportDeclaration
	Functions []FunctionDeclaration
}

type IfStatement struct {
	token          *lexer.Token
	Condition      Expression
	Statements     []Statement
	ElseStatements []Statement
}

type IdentifierExpression struct {
	token      *lexer.Token
	Identifier Identifier
}

func (c *Class) GetToken() *lexer.Token {
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

func (s *SingleImportDeclaration) GetToken() *lexer.Token {
	return s.token
}

func (s *SingleImportDeclaration) String() string {
	var buf bytes.Buffer
	buf.WriteString("(import ")
	buf.WriteString(s.Name.String())
	buf.WriteString(")")
	return buf.String()
}

func (i *ImportDeclarationList) GetToken() *lexer.Token {
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

func (i *Identifier) GetToken() *lexer.Token {
	return i.token
}

func (i *Identifier) String() string {
	return i.Value
}

func (b *BinaryExpression) GetToken() *lexer.Token {
	return b.token
}

func (b *BinaryExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString("(")
	buf.WriteString(b.token.GetRawValue())
	buf.WriteString(" ")
	buf.WriteString(b.Left.String())
	buf.WriteString(" ")
	buf.WriteString(b.Right.String())
	buf.WriteString(")")
	return buf.String()
}

func (m *MethodCallExpression) GetToken() *lexer.Token {
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

func (s *StringLiteralExpression) GetToken() *lexer.Token {
	return s.token
}

func (s *StringLiteralExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString("(string \"")
	buf.WriteString(s.Value)
	buf.WriteString("\")")
	return buf.String()
}

func (f *FunctionDeclaration) GetToken() *lexer.Token {
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

func (a *ArgumentDeclaration) GetToken() *lexer.Token {
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

func (e *ExpressionStatement) GetToken() *lexer.Token {
	return e.token
}

func (e *ExpressionStatement) String() string {
	return e.ExpressionValue.String()
}

func (t *Type) GetToken() *lexer.Token {
	return t.token
}

func (t *Type) String() string {
	return t.Name
}

func (i *IfStatement) GetToken() *lexer.Token {
	return i.token
}

func (i *IfStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("(if ")
	buf.WriteString(i.Condition.String())
	for _, s := range i.Statements {
		buf.WriteString(" ")
		buf.WriteString(s.String())
	}
	if len(i.ElseStatements) > 0 {
		buf.WriteString(" (else")
		for _, s := range i.ElseStatements {
			buf.WriteString(" ")
			buf.WriteString(s.String())
		}
		buf.WriteString(")")
	}
	buf.WriteString(")")
	return buf.String()
}

func (i *IdentifierExpression) GetToken() *lexer.Token {
	return i.token
}

func (i *IdentifierExpression) String() string {
	return i.Identifier.String()
}
