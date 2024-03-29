package parser

import (
	"fmt"
	"goMud/internal/gmsl/lexer"
	"log"
)

type Parser struct {
	lexer *lexer.Lexer
}

func NewParser(l *lexer.Lexer) *Parser {
	return &Parser{l}
}

func (p *Parser) Parse() *Class {
	class := p.parseClass()
	fmt.Println(class.PrettyPrint(0))
	return class
}

func (p *Parser) parseClass() *Class {
	log.Println("Parsing class")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.PackageToken {
		panic("Expected class")
	}

	name := p.parseIdentifier()
	class := newClass(name, token)
	for {
		peeked := p.lexer.Peek()
		switch peeked.Typ {
		case lexer.ImportToken:
			imports := p.parseImportDeclarations()
			class.Imports = append(class.Imports, imports...)
		case lexer.FuncToken:
			functions := p.parseFunctionDeclarations()
			class.Functions = append(class.Functions, functions...)
		case lexer.EofToken:
			return class
		default:
			p.unexpectedToken(peeked)
		}
	}
}

func (p *Parser) parseIdentifier() *Identifier {
	log.Println("Parsing identifier")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.IdentifierToken {
		log.Panicln("Expected identifier, got", token.String())
	}

	return newIdentifier(token)
}

func (p *Parser) parseImportDeclarations() []ImportDeclaration {
	log.Println("Parsing import declarations")
	token := p.lexer.Peek()
	tokenDeclarations := make([]ImportDeclaration, 0)
	if token.Typ == lexer.ImportToken {
		declaration := p.parseImportDeclaration()
		tokenDeclarations = append(tokenDeclarations, declaration)
	}

	return tokenDeclarations
}

func (p *Parser) parseImportDeclaration() ImportDeclaration {
	log.Println("Parsing import declaration")
	tokens := p.lexer.PeekSome(2)
	if len(tokens) < 2 {
		log.Panicln("Expected import declaration")
	}

	switch tokens[1].Typ {
	case lexer.IdentifierToken:
		return p.parseSingleImportDeclaration()
	case lexer.OpenParenToken:
		return p.parseImportDeclarationList()
	default:
		p.unexpectedToken(tokens[1])
	}
	return nil
}

func (p *Parser) parseSingleImportDeclaration() ImportDeclaration {
	log.Println("Parsing single import declaration")
	token := p.lexer.ReadNext()
	name := p.parseStringValue()

	return newSingleImportDeclaration(name, token)
}

func (p *Parser) parseImportDeclarationList() ImportDeclaration {
	log.Println("Parsing import declaration list")
	token := p.lexer.ReadNext()
	imports := make([]Identifier, 0)
	skip := p.lexer.ReadNext()
	if skip.Typ != lexer.OpenParenToken {
		p.unexpectedTokenExpected(lexer.OpenParenToken, skip)
	}
	for {
		token := p.lexer.Peek()
		if token.Typ == lexer.CloseParenToken {
			p.lexer.ReadNext()
			break
		}
		imports = append(imports, *p.parseStringValue())
	}

	return newImportDeclarationList(&imports, token)
}

func (p *Parser) parseStringValue() *Identifier {
	log.Println("Parsing string value")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.StringToken {
		log.Panicln("Expected string value, got", token.String())
	}

	return newIdentifier(token)

}

func (p *Parser) parseFunctionDeclarations() []FunctionDeclaration {
	log.Println("Parsing function declarations")
	var functions []FunctionDeclaration
	for {
		token := p.lexer.Peek()
		if token.Typ == lexer.FuncToken {
			functions = append(functions, p.parseFunctionDeclaration())
		} else {
			break
		}
	}
	return functions
}

func (p *Parser) parseFunctionDeclaration() FunctionDeclaration {
	log.Println("Parsing function declaration")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.FuncToken {
		log.Panicln("Expected FuncToken, got", token.String())
	}

	name := p.parseIdentifier()
	arguments := p.parseArgumentDeclarations()
	returnTypes := make([]Type, 0)
	switch p.lexer.Peek().Typ {
	case lexer.TypeToken:
		returnTypes = append(returnTypes, *p.parseType())
	case lexer.OpenBraceToken:
		break
	default:
		p.unexpectedToken(p.lexer.Peek())
	}
	statements := p.parseStatements()

	declaration := newFunctionDeclaration(name, &arguments, returnTypes, &statements, token)
	return *declaration
}

func (p *Parser) parseArgumentDeclarations() []ArgumentDeclaration {
	log.Println("Parsing arguments")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.OpenParenToken {
		p.unexpectedTokenExpected(lexer.OpenParenToken, token)
	}

	arguments := make([]ArgumentDeclaration, 0)
	for {
		token := p.lexer.Peek()
		if token.Typ == lexer.CloseParenToken {
			p.lexer.ReadNext()
			break
		}

		arguments = append(arguments, p.parseArgumentDeclaration())
	}

	return arguments
}

func (p *Parser) parseArgumentDeclaration() ArgumentDeclaration {
	log.Println("Parsing argument")
	name := p.parseIdentifier()
	typ := p.parseType()
	return *newArgumentDeclaration(name, typ, name.token)
}

func (p *Parser) parseType() *Type {
	log.Println("Parsing type")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.TypeToken {
		log.Panicln("Expected TypeToken, got", token.String())
	}

	return newType(token)
}

func (p *Parser) parseStatements() []Statement {
	log.Println("Parsing statements")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.OpenBraceToken {
		p.unexpectedTokenExpected(lexer.OpenBraceToken, token)
	}
	statements := make([]Statement, 0)
	for {
		token := p.lexer.Peek()
		if token.Typ == lexer.CloseBraceToken {
			p.lexer.ReadNext()
			break
		}
		statements = append(statements, p.parseStatement())
	}
	return statements
}

func (p *Parser) parseStatement() Statement {
	log.Println("Parsing statement")
	peeked := p.lexer.PeekSome(2)

	switch peeked[0].Typ {
	case lexer.VarToken:
		return p.parseVariableDeclarationStatement()
	case lexer.IdentifierToken:
		switch peeked[1].Typ {
		case lexer.AssignToken:
			return p.parseVariableAssignmentStatement()
		case lexer.CreateAndAssignToken:
			return p.parseVariableCreateAndAssignStatement()
		case lexer.MethodCallToken:
			return p.parseExpressionStatement()
		default:
			p.unexpectedToken(peeked[1])
		}
	case lexer.IfToken:
		return p.parseIfStatement()
	case lexer.ReturnToken:
		return p.parseReturnStatement()
	default:
		p.unexpectedToken(peeked[0])
	}

	return nil
}

func (p *Parser) parseExpressionStatement() Statement {
	log.Println("Parsing ExpressionValue statement")
	token := p.lexer.Peek()
	expression := p.parseExpression()
	return newExpressionStatement(&expression, token)
}

func (p *Parser) parseExpression() Expression {
	log.Println("Parsing ExpressionValue")
	peeked := p.lexer.PeekSome(2)

	tree := NewExpressionTree()

	switch peeked[0].Typ {
	case lexer.IdentifierToken, lexer.StringToken, lexer.NumericToken:
	default:
		p.unexpectedToken(peeked[0])
	}

	for {
		peeked := p.lexer.PeekSome(2)
		expression, done := p.tryAddExpression(peeked, tree)
		if done {
			return expression
		}
	}
}

func (p *Parser) tryAddExpression(peeked []*lexer.Token, tree *ExpressionTree) (Expression, bool) {
	switch {
	case peeked[1].Typ == lexer.MethodCallToken:
		if !tree.CanAddLeaf() {
			return tree.GetExpression(), true
		}
		tree.AddExpression(p.parseMethodCallExpression())

	case peeked[0].Typ == lexer.StringToken:
		if !tree.CanAddLeaf() {
			return tree.GetExpression(), true
		}
		tree.AddExpression(p.parseStringLiteralExpression())
	case peeked[0].Typ == lexer.NumericToken:
		if !tree.CanAddLeaf() {
			return tree.GetExpression(), true
		}
		tree.AddExpression(p.parseNumericLiteralExpression())
	case peeked[0].Typ == lexer.IdentifierToken:
		if !tree.CanAddLeaf() {
			return tree.GetExpression(), true
		}
		tree.AddExpression(p.parseIdentifierExpression())
	case peeked[0].Typ == lexer.AddToken, peeked[0].Typ == lexer.EqualToken, peeked[0].Typ == lexer.DivideToken, peeked[0].Typ == lexer.MultiplyToken, peeked[0].Typ == lexer.SubtractToken, peeked[0].Typ == lexer.ModuloToken:
		if !tree.CanAddBranch() {
			p.unexpectedToken(peeked[0])
		}
		tree.AddExpression(newBinaryExpression(p.lexer.ReadNext()))
	case !tree.CanAddLeaf():
		return tree.GetExpression(), true
	default:
		p.unexpectedToken(peeked[0])
	}
	return nil, false
}

func (p *Parser) parseMethodCallExpression() Expression {
	log.Println("Parsing method call ExpressionValue")
	token := p.lexer.Peek()
	p.unexpectedTokenExpected(lexer.IdentifierToken, token)

	objectName := p.parseIdentifier()

	token = p.lexer.ReadNext()
	if token.Typ != lexer.MethodCallToken {
		log.Panicln("Expected MethodCallToken, got", token.String())
	}
	methodName := p.parseIdentifier()
	arguments := p.parseArguments()

	return newMethodCallExpression(objectName, methodName, &arguments, token)
}

func (p *Parser) parseArguments() []Expression {
	log.Println("Parsing arguments")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.OpenParenToken {
		p.unexpectedTokenExpected(lexer.OpenParenToken, token)
	}

	arguments := make([]Expression, 0)
	for {
		token := p.lexer.Peek()
		if token.Typ == lexer.CloseParenToken {
			p.lexer.ReadNext()
			break
		}

		arguments = append(arguments, p.parseExpression())
	}

	return arguments
}

func (p *Parser) parseStringLiteralExpression() Expression {
	log.Println("Parsing string literal ExpressionValue")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.StringToken {
		log.Panicln("Expected StringToken, got", token.String())
	}

	return newStringLiteralExpression(token)
}

func (p *Parser) unexpectedToken(token *lexer.Token) {
	log.Panicln("Unexpected token", token.String())
}

func (p *Parser) unexpectedTokenExpected(expected lexer.TokenType, actual *lexer.Token) {
	if actual.Typ != expected {
		log.Panicln("Unexpected token", actual.String(), "expected", expected)
	}
}

func (p *Parser) parseIfStatement() Statement {
	log.Println("Parsing if statement")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.IfToken {
		log.Panicln("Expected IfToken, got", token.String())
	}

	condition := p.parseExpression()

	token = p.lexer.Peek()
	if token.Typ != lexer.OpenBraceToken {
		p.unexpectedTokenExpected(lexer.OpenBraceToken, token)
	}

	statements := p.parseStatements()

	var elseStatements []Statement
	token = p.lexer.Peek()
	if token.Typ == lexer.ElseToken {
		p.lexer.ReadNext() // consume the 'else' token

		token = p.lexer.Peek()
		if token.Typ != lexer.OpenBraceToken {
			p.unexpectedTokenExpected(lexer.OpenBraceToken, token)
		}

		elseStatements = p.parseStatements()
	}

	return newIfStatement(&condition, &statements, &elseStatements, token)
}

func (p *Parser) parseIdentifierExpression() Expression {
	log.Println("Parsing identifier ExpressionValue")
	token := p.lexer.Peek()
	p.unexpectedTokenExpected(lexer.IdentifierToken, token)

	return newIdentifierExpression(p.parseIdentifier(), token)
}

func (p *Parser) parseVariableDeclarationStatement() Statement {
	log.Println("Parsing variable declaration statement")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.VarToken {
		log.Panicln("Expected VarToken, got", token.String())
	}

	name := p.parseIdentifier()
	typ := p.parseType()

	return newVariableDeclarationStatement(name, typ, token)
}

func (p *Parser) parseVariableAssignmentStatement() Statement {
	log.Println("Parsing variable assignment statement")
	token := p.lexer.Peek()
	p.unexpectedTokenExpected(lexer.IdentifierToken, token)

	name := p.parseIdentifier()
	token = p.expect(lexer.AssignToken, "AssignToken")
	expression := p.parseExpression()

	return newVariableAssignmentStatement(name, &expression, token)
}

func (p *Parser) parseVariableCreateAndAssignStatement() Statement {
	log.Println("Parsing variable create and assign statement")
	token := p.lexer.Peek()
	if token.Typ != lexer.IdentifierToken {
		p.unexpectedTokenExpected(lexer.IdentifierToken, token)
	}

	name := p.parseIdentifier()
	token = p.expect(lexer.CreateAndAssignToken, "CreateAndAssignToken")
	expression := p.parseExpression()

	return newVariableCreateAndAssignStatement(name, &expression, token)
}

func (p *Parser) expect(token lexer.TokenType, s string) *lexer.Token {
	read := p.lexer.ReadNext()
	if token != read.Typ {
		log.Panicln("Expected", s, "got", token.String())
	}
	return read
}

func (p *Parser) parseNumericLiteralExpression() Expression {
	log.Println("Parsing numeric literal ExpressionValue")
	token := p.lexer.ReadNext()
	p.unexpectedTokenExpected(lexer.NumericToken, token)

	return newNumericLiteralExpression(token)
}

func (p *Parser) parseReturnStatement() Statement {
	log.Println("Parsing return statement")
	token := p.expect(lexer.ReturnToken, "ReturnToken")

	expression := p.parseExpression()

	return newReturnStatement(&expression, token)
}
