package parser

import (
	"goMud/internal/gmsl/lexer"

	"github.com/sirupsen/logrus"
)

type ParseState int

const (
	StateClass ParseState = iota
	StateClassBody
	StateImportDecl
	StateImportDeclList
	StateFuncDecl
	StateFuncDeclArgs
	StateFuncDeclReturn
	StateStatements
	StateStatementsLoop
	StateStatement
	StateIfStatement
	StateIfStatementBody
	StateIfStatementElse
	StateVariableDeclStmt
	StateVariableAssignStmt
	StateVariableCreateAssignStmt
	StateReturnStmt
	StateExpressionStmt
	StateExpression
	StateExpressionLoop
	StateMethodCall
	StateMethodCallArgs
)

type Frame struct {
	State ParseState
	Token *lexer.Token

	// Context trackers
	Class           *Class
	FuncDecl        *FunctionDeclaration
	IfStmt          *IfStatement
	MethodCall      *MethodCallExpression
	ExprTree        *ExpressionTree
	VarDecl         *VariableDeclarationStatement
	VarAssign       *VariableAssignmentStatement
	VarCreateAssign *VariableCreateAndAssignStatement
	ReturnStmt      *ReturnStatement

	// Aggregators
	Stmts   []Statement
	Exprs   []Expression
	Args    []ArgumentDeclaration
	Imports []Identifier

	// Temporaries
	Identifier *Identifier
	Type       *Type
}

// Parser processes a stream of lexer tokens into an Abstract Syntax Tree (AST).
type Parser struct {
	lexer     *lexer.Lexer
	stack     []*Frame
	nodeStack []AstNode
}

// NewParser creates a new Parser with the provided Lexer instance.
func NewParser(l *lexer.Lexer) *Parser {
	return &Parser{
		lexer:     l,
		stack:     make([]*Frame, 0),
		nodeStack: make([]AstNode, 0),
	}
}

func (p *Parser) push(f *Frame) {
	p.stack = append(p.stack, f)
}

func (p *Parser) pop() *Frame {
	if len(p.stack) == 0 {
		return nil
	}
	f := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return f
}

func (p *Parser) pushNode(n AstNode) {
	p.nodeStack = append(p.nodeStack, n)
}

func (p *Parser) popNode() AstNode {
	if len(p.nodeStack) == 0 {
		return nil
	}
	n := p.nodeStack[len(p.nodeStack)-1]
	p.nodeStack = p.nodeStack[:len(p.nodeStack)-1]
	return n
}

func (p *Parser) unexpectedToken(token *lexer.Token) {
	logrus.Panic("Unexpected token ", token.String())
}

func (p *Parser) unexpectedTokenExpected(expected lexer.TokenType, actual *lexer.Token) {
	if actual.Typ != expected {
		logrus.Panic("Unexpected token ", actual.String(), " expected ", expected)
	}
}

func (p *Parser) expect(expected lexer.TokenType, name string) *lexer.Token {
	read := p.lexer.ReadNext()
	if expected != read.Typ {
		logrus.Panic("Expected ", name, " got ", read.String())
	}
	return read
}

func (p *Parser) parseIdentifier() *Identifier {
	logrus.Debug("Parsing identifier")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.IdentifierToken && token.Typ != lexer.ContextToken {
		logrus.Panic("Expected identifier, got ", token.String())
	}
	return newIdentifier(token)
}

func (p *Parser) parseStringValue() *Identifier {
	logrus.Debug("Parsing string value")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.StringToken {
		logrus.Panic("Expected string value, got ", token.String())
	}
	return newIdentifier(token)
}

func (p *Parser) parseType() *Type {
	logrus.Debug("Parsing type")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.TypeToken {
		logrus.Panic("Expected TypeToken, got ", token.String())
	}
	return newType(token)
}

// Parse orchestrates the state machine loop to consume tokens and construct the Class AST node.
func (p *Parser) Parse() *Class {
	p.push(&Frame{State: StateClass})

	var resultClass *Class

	for len(p.stack) > 0 {
		frame := p.pop()

		switch frame.State {
		case StateClass:
			p.handleStateClass(frame)
		case StateClassBody:
			if c := p.handleStateClassBody(frame); c != nil {
				resultClass = c
			}
		case StateImportDecl:
			p.handleStateImportDecl(frame)
		case StateImportDeclList:
			p.handleStateImportDeclList(frame)
		case StateFuncDecl:
			p.handleStateFuncDecl(frame)
		case StateFuncDeclArgs:
			p.handleStateFuncDeclArgs(frame)
		case StateFuncDeclReturn:
			p.handleStateFuncDeclReturn(frame)
		case StateStatements:
			p.handleStateStatements(frame)
		case StateStatementsLoop:
			p.handleStateStatementsLoop(frame)
		case StateStatement:
			p.handleStateStatement(frame)
		case StateVariableDeclStmt:
			p.handleStateVariableDeclStmt(frame)
		case StateVariableAssignStmt:
			p.handleStateVariableAssignStmt(frame)
		case StateVariableCreateAssignStmt:
			p.handleStateVariableCreateAssignStmt(frame)
		case StateExpressionStmt:
			p.handleStateExpressionStmt(frame)
		case StateIfStatement:
			p.handleStateIfStatement(frame)
		case StateReturnStmt:
			p.handleStateReturnStmt(frame)
		case StateExpression:
			p.handleStateExpression(frame)
		case StateExpressionLoop:
			p.handleStateExpressionLoop(frame)
		case StateMethodCall:
			p.handleStateMethodCall(frame)
		case StateMethodCallArgs:
			p.handleStateMethodCallArgs(frame)
		}
	}

	return resultClass
}

func (p *Parser) handleStateClass(frame *Frame) {
	logrus.Debug("Parsing class")
	token := p.lexer.Peek()
	frame.Class = newClass(token)
	frame.State = StateClassBody
	p.push(frame)
}

func (p *Parser) handleStateClassBody(frame *Frame) *Class {
	for len(p.nodeStack) > 0 {
		node := p.popNode()
		switch n := node.(type) {
		case *FunctionDeclaration:
			frame.Class.Functions = append(frame.Class.Functions, *n)
		case *SingleImportDeclaration:
			frame.Class.Imports = append(frame.Class.Imports, n)
		case *ImportDeclarationList:
			frame.Class.Imports = append(frame.Class.Imports, n)
		default:
			logrus.Panic("Unexpected node in StateClassBody: ", n)
		}
	}

	peeked := p.lexer.Peek()
	switch peeked.Typ {
	case lexer.ImportToken:
		p.push(frame)
		p.push(&Frame{State: StateImportDecl})
	case lexer.FuncToken:
		p.push(frame)
		p.push(&Frame{State: StateFuncDecl})
	case lexer.EofToken:
		return frame.Class
	default:
		p.unexpectedToken(peeked)
	}
	return nil
}

func (p *Parser) handleStateImportDecl(frame *Frame) {
	logrus.Debug("Parsing import declarations")
	token := p.lexer.Peek()
	if token.Typ == lexer.ImportToken {
		logrus.Debug("Parsing import declaration")
		tokens := p.lexer.PeekSome(2)
		if len(tokens) < 2 {
			logrus.Panic("Expected import declaration")
		}
		switch tokens[1].Typ {
		case lexer.StringToken:
			logrus.Debug("Parsing single import declaration")
			t2 := p.lexer.ReadNext()
			name := p.parseStringValue()
			p.pushNode(newSingleImportDeclaration(name, t2))
		case lexer.OpenParenToken:
			logrus.Debug("Parsing import declaration list")
			t2 := p.lexer.ReadNext()
			skip := p.lexer.ReadNext()
			if skip.Typ != lexer.OpenParenToken {
				p.unexpectedTokenExpected(lexer.OpenParenToken, skip)
			}
			newFrame := &Frame{State: StateImportDeclList, Token: t2}
			p.push(newFrame)
		default:
			p.unexpectedToken(tokens[1])
		}
	}
}

func (p *Parser) handleStateImportDeclList(frame *Frame) {
	token := p.lexer.Peek()
	if token.Typ == lexer.CloseParenToken {
		p.lexer.ReadNext()
		p.pushNode(newImportDeclarationList(&frame.Imports, frame.Token))
	} else {
		frame.Imports = append(frame.Imports, *p.parseStringValue())
		p.push(frame)
	}
}

func (p *Parser) handleStateFuncDecl(frame *Frame) {
	logrus.Debug("Parsing function declaration")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.FuncToken {
		logrus.Panic("Expected FuncToken, got ", token.String())
	}
	name := p.parseIdentifier()
	frame.FuncDecl = &FunctionDeclaration{
		token:       token,
		Name:        *name,
		ReturnTypes: make([]Type, 0),
	}

	logrus.Debug("Parsing arguments")
	t2 := p.lexer.ReadNext()
	if t2.Typ != lexer.OpenParenToken {
		p.unexpectedTokenExpected(lexer.OpenParenToken, t2)
	}

	frame.State = StateFuncDeclArgs
	p.push(frame)
}

func (p *Parser) handleStateFuncDeclArgs(frame *Frame) {
	token := p.lexer.Peek()
	if token.Typ == lexer.CloseParenToken {
		p.lexer.ReadNext()
		frame.State = StateFuncDeclReturn
		p.push(frame)
	} else {
		logrus.Debug("Parsing argument")
		argName := p.parseIdentifier()
		argType := p.parseType()
		arg := newArgumentDeclaration(argName, argType, argName.token)
		frame.FuncDecl.Arguments = append(frame.FuncDecl.Arguments, *arg)
		p.push(frame)
	}
}

func (p *Parser) handleStateFuncDeclReturn(frame *Frame) {
	peeked := p.lexer.Peek()
	switch peeked.Typ {
	case lexer.TypeToken:
		frame.FuncDecl.ReturnTypes = append(frame.FuncDecl.ReturnTypes, *p.parseType())
		// The original parser handles return types in a switch but does not loop.
		// Next is Statements.
		fallthrough
	case lexer.OpenBraceToken:
		frame.State = StateStatements
		p.push(frame)
	default:
		p.unexpectedToken(peeked)
	}
}

func (p *Parser) handleStateStatements(frame *Frame) {
	if frame.FuncDecl != nil {
		logrus.Debug("Parsing statements")
		token := p.lexer.ReadNext()
		if token.Typ != lexer.OpenBraceToken {
			p.unexpectedTokenExpected(lexer.OpenBraceToken, token)
		}
		frame.State = StateStatementsLoop
		p.push(frame)
	} else if frame.IfStmt != nil {
		frame.State = StateStatementsLoop
		p.push(frame)
	} else {
		logrus.Debug("Parsing statements")
		token := p.lexer.ReadNext()
		if token.Typ != lexer.OpenBraceToken {
			p.unexpectedTokenExpected(lexer.OpenBraceToken, token)
		}
		frame.State = StateStatementsLoop
		p.push(frame)
	}
}

func (p *Parser) handleStateStatementsLoop(frame *Frame) {
	for len(p.nodeStack) > 0 {
		stmt := p.popNode().(Statement)
		frame.Stmts = append(frame.Stmts, stmt)
	}

	token := p.lexer.Peek()
	if token.Typ == lexer.CloseBraceToken {
		p.lexer.ReadNext()
		if frame.FuncDecl != nil {
			frame.FuncDecl.Statements = frame.Stmts
			p.pushNode(frame.FuncDecl)
		} else if frame.IfStmt != nil {
			if frame.IfStmt.Statements == nil {
				frame.IfStmt.Statements = frame.Stmts
				frame.Stmts = nil

				token = p.lexer.Peek()
				if token.Typ == lexer.ElseToken {
					p.lexer.ReadNext()
					token = p.lexer.Peek()
					if token.Typ != lexer.OpenBraceToken {
						p.unexpectedTokenExpected(lexer.OpenBraceToken, token)
					}
					p.lexer.ReadNext()
					p.push(frame)
				} else {
					p.pushNode(frame.IfStmt)
				}
			} else {
				frame.IfStmt.ElseStatements = frame.Stmts
				p.pushNode(frame.IfStmt)
			}
		}
	} else {
		p.push(frame)
		p.push(&Frame{State: StateStatement})
	}
}

func (p *Parser) handleStateStatement(frame *Frame) {
	logrus.Debug("Parsing statement")
	peeked := p.lexer.PeekSome(2)
	switch peeked[0].Typ {
	case lexer.VarToken:
		p.push(&Frame{State: StateVariableDeclStmt})
	case lexer.IdentifierToken, lexer.ContextToken:
		switch peeked[1].Typ {
		case lexer.AssignToken:
			p.push(&Frame{State: StateVariableAssignStmt})
		case lexer.CreateAndAssignToken:
			p.push(&Frame{State: StateVariableCreateAssignStmt})
		case lexer.MethodCallToken:
			p.push(&Frame{State: StateExpressionStmt, Token: p.lexer.Peek()})
		default:
			p.unexpectedToken(peeked[1])
		}
	case lexer.IfToken:
		p.push(&Frame{State: StateIfStatement})
	case lexer.ReturnToken:
		p.push(&Frame{State: StateReturnStmt})
	default:
		p.unexpectedToken(peeked[0])
	}
}

func (p *Parser) handleStateVariableDeclStmt(frame *Frame) {
	logrus.Debug("Parsing variable declaration statement")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.VarToken {
		logrus.Panic("Expected VarToken, got ", token.String())
	}
	name := p.parseIdentifier()
	typ := p.parseType()
	p.pushNode(newVariableDeclarationStatement(name, typ, token))
}

func (p *Parser) handleStateVariableAssignStmt(frame *Frame) {
	if frame.Identifier == nil {
		logrus.Debug("Parsing variable assignment statement")
		token := p.lexer.Peek()
		p.unexpectedTokenExpected(lexer.IdentifierToken, token)

		name := p.parseIdentifier()
		assignToken := p.expect(lexer.AssignToken, "AssignToken")
		frame.Identifier = name
		frame.Token = assignToken

		p.push(frame)
		p.push(&Frame{State: StateExpression})
	} else {
		expr := p.popNode().(Expression)
		p.pushNode(newVariableAssignmentStatement(frame.Identifier, &expr, frame.Token))
	}
}

func (p *Parser) handleStateVariableCreateAssignStmt(frame *Frame) {
	if frame.Identifier == nil {
		logrus.Debug("Parsing variable create and assign statement")
		token := p.lexer.Peek()
		if token.Typ != lexer.IdentifierToken {
			p.unexpectedTokenExpected(lexer.IdentifierToken, token)
		}

		name := p.parseIdentifier()
		assignToken := p.expect(lexer.CreateAndAssignToken, "CreateAndAssignToken")
		frame.Identifier = name
		frame.Token = assignToken

		p.push(frame)
		p.push(&Frame{State: StateExpression})
	} else {
		expr := p.popNode().(Expression)
		p.pushNode(newVariableCreateAndAssignStatement(frame.Identifier, &expr, frame.Token))
	}
}

func (p *Parser) handleStateExpressionStmt(frame *Frame) {
	if frame.ExprTree == nil {
		p.push(frame)
		p.push(&Frame{State: StateExpression})
		frame.ExprTree = NewExpressionTree()
	} else {
		expr := p.popNode().(Expression)
		p.pushNode(newExpressionStatement(&expr, frame.Token))
	}
}

func (p *Parser) handleStateIfStatement(frame *Frame) {
	if frame.IfStmt == nil {
		logrus.Debug("Parsing if statement")
		token := p.lexer.ReadNext()
		if token.Typ != lexer.IfToken {
			logrus.Panic("Expected IfToken, got ", token.String())
		}
		frame.IfStmt = &IfStatement{token: token}

		p.push(frame)
		p.push(&Frame{State: StateExpression})
	} else {
		expr := p.popNode().(Expression)
		frame.IfStmt.Condition = expr

		token := p.lexer.Peek()
		if token.Typ != lexer.OpenBraceToken {
			p.unexpectedTokenExpected(lexer.OpenBraceToken, token)
		}
		p.lexer.ReadNext()

		frame.State = StateStatementsLoop
		p.push(frame)
	}
}

func (p *Parser) handleStateReturnStmt(frame *Frame) {
	if frame.ReturnStmt == nil {
		logrus.Debug("Parsing return statement")
		token := p.expect(lexer.ReturnToken, "ReturnToken")
		frame.Token = token
		frame.ReturnStmt = &ReturnStatement{}

		p.push(frame)
		p.push(&Frame{State: StateExpression})
	} else {
		expr := p.popNode().(Expression)
		p.pushNode(newReturnStatement(&expr, frame.Token))
	}
}

func (p *Parser) handleStateExpression(frame *Frame) {
	logrus.Debug("Parsing ExpressionValue")
	peeked := p.lexer.PeekSome(2)
	switch peeked[0].Typ {
	case lexer.IdentifierToken, lexer.ContextToken, lexer.StringToken, lexer.NumericToken:
	default:
		p.unexpectedToken(peeked[0])
	}

	frame.ExprTree = NewExpressionTree()
	frame.State = StateExpressionLoop
	p.push(frame)
}

func (p *Parser) handleStateExpressionLoop(frame *Frame) {
	if len(p.nodeStack) > 0 {
		expr := p.popNode().(Expression)
		frame.ExprTree.AddExpression(expr)
	}

	peeked := p.lexer.PeekSome(2)

	if len(peeked) > 1 && peeked[1].Typ == lexer.MethodCallToken {
		if !frame.ExprTree.CanAddLeaf() {
			p.pushNode(frame.ExprTree.GetExpression())
			return
		}
		p.push(frame)
		p.push(&Frame{State: StateMethodCall})
		return
	}

	if peeked[0].Typ == lexer.StringToken {
		if !frame.ExprTree.CanAddLeaf() {
			p.pushNode(frame.ExprTree.GetExpression())
			return
		}
		logrus.Debug("Parsing string literal ExpressionValue")
		token := p.lexer.ReadNext()
		if token.Typ != lexer.StringToken {
			logrus.Panic("Expected StringToken, got ", token.String())
		}
		frame.ExprTree.AddExpression(newStringLiteralExpression(token))
		p.push(frame)
		return
	}

	if peeked[0].Typ == lexer.NumericToken {
		if !frame.ExprTree.CanAddLeaf() {
			p.pushNode(frame.ExprTree.GetExpression())
			return
		}
		logrus.Debug("Parsing numeric literal ExpressionValue")
		token := p.lexer.ReadNext()
		p.unexpectedTokenExpected(lexer.NumericToken, token)
		frame.ExprTree.AddExpression(newNumericLiteralExpression(token))
		p.push(frame)
		return
	}

	if peeked[0].Typ == lexer.IdentifierToken || peeked[0].Typ == lexer.ContextToken {
		if !frame.ExprTree.CanAddLeaf() {
			p.pushNode(frame.ExprTree.GetExpression())
			return
		}
		logrus.Debug("Parsing identifier ExpressionValue")
		token := p.lexer.Peek()
		if token.Typ != lexer.IdentifierToken && token.Typ != lexer.ContextToken {
			p.unexpectedTokenExpected(lexer.IdentifierToken, token)
		}
		id := p.parseIdentifier()
		if token.Typ == lexer.ContextToken {
			frame.ExprTree.AddExpression(newContextExpression(id, token))
		} else {
			frame.ExprTree.AddExpression(newIdentifierExpression(id, token))
		}
		p.push(frame)
		return
	}

	if peeked[0].Typ == lexer.AddToken || peeked[0].Typ == lexer.EqualToken || peeked[0].Typ == lexer.DivideToken || peeked[0].Typ == lexer.MultiplyToken || peeked[0].Typ == lexer.SubtractToken || peeked[0].Typ == lexer.ModuloToken {
		if !frame.ExprTree.CanAddBranch() {
			p.unexpectedToken(peeked[0])
		}
		frame.ExprTree.AddExpression(newBinaryExpression(p.lexer.ReadNext()))
		p.push(frame)
		return
	}

	if !frame.ExprTree.CanAddLeaf() {
		p.pushNode(frame.ExprTree.GetExpression())
		return
	}

	p.unexpectedToken(peeked[0])
}

func (p *Parser) handleStateMethodCall(frame *Frame) {
	if frame.Identifier == nil {
		logrus.Debug("Parsing method call ExpressionValue")
		token := p.lexer.Peek()
		if token.Typ != lexer.IdentifierToken && token.Typ != lexer.ContextToken {
			p.unexpectedTokenExpected(lexer.IdentifierToken, token)
		}

		objectName := p.parseIdentifier()

		mtoken := p.lexer.ReadNext()
		if mtoken.Typ != lexer.MethodCallToken {
			logrus.Panic("Expected MethodCallToken, got ", mtoken.String())
		}
		methodName := p.parseIdentifier()

		frame.Identifier = objectName
		frame.MethodCall = &MethodCallExpression{
			token:      mtoken,
			ObjectName: *objectName,
			MethodName: *methodName,
			Arguments:  make([]Expression, 0),
		}

		logrus.Debug("Parsing arguments")
		t2 := p.lexer.ReadNext()
		if t2.Typ != lexer.OpenParenToken {
			p.unexpectedTokenExpected(lexer.OpenParenToken, t2)
		}

		frame.State = StateMethodCallArgs
		p.push(frame)
	}
}

func (p *Parser) handleStateMethodCallArgs(frame *Frame) {
	for len(p.nodeStack) > 0 {
		expr := p.popNode().(Expression)
		frame.Exprs = append(frame.Exprs, expr)
	}

	token := p.lexer.Peek()
	if token.Typ == lexer.CloseParenToken {
		p.lexer.ReadNext()
		frame.MethodCall.Arguments = frame.Exprs
		p.pushNode(frame.MethodCall)
	} else {
		p.push(frame)
		p.push(&Frame{State: StateExpression})
	}
}
