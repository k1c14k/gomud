package parser

import (
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
	return p.parseClass()
}

func (p *Parser) parseClass() *Class {
	log.Println("Parsing class")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.PackageToken {
		panic("Expected class")
	}

	name := p.parseIdentifier()
	class := Class{token: token, Name: name, Imports: make([]ImportDeclaration, 0)}
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
			return &class
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

	return &Identifier{token: token, Value: token.GetRawValue()}
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

	return &SingleImportDeclaration{token: token, Name: name}
}

func (p *Parser) parseImportDeclarationList() ImportDeclaration {
	log.Println("Parsing import declaration list")
	token := p.lexer.ReadNext()
	imports := make([]*Identifier, 0)
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
		imports = append(imports, p.parseStringValue())
	}

	return &ImportDeclarationList{token: token, Imports: imports}
}

func (p *Parser) parseStringValue() *Identifier {
	log.Println("Parsing string value")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.StringToken {
		log.Panicln("Expected string value, got", token.String())
	}

	return &Identifier{token: token, Value: token.GetRawValue()}

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
	statements := p.parseStatements()

	declaration := FunctionDeclaration{token: token, Name: name, Arguments: arguments, Statements: statements}
	return declaration
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
	return ArgumentDeclaration{token: name.token, Name: name, Typ: typ}
}

func (p *Parser) parseType() *Type {
	log.Println("Parsing type")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.TypeToken {
		log.Panicln("Expected TypeToken, got", token.String())
	}

	return &Type{token: token, Name: token.GetRawValue()}
}

func (p *Parser) parseStatements() []Statement {
	log.Println("Parsing statements")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.OpenBraceToken {
		log.Panicln("Expected OpenBraceToken, got", token.String())
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

	if peeked[0].Typ == lexer.IdentifierToken {
		switch peeked[1].Typ {
		case lexer.MethodCallToken:
			return p.parseExpressionStatement()
		default:
			p.unexpectedToken(peeked[1])
		}
	} else {
		p.unexpectedToken(peeked[0])
	}
	return nil
}

func (p *Parser) parseExpressionStatement() Statement {
	log.Println("Parsing ExpressionValue statement")
	token := p.lexer.Peek()
	return &ExpressionStatement{token: token, ExpressionValue: p.parseExpression()}
}

func (p *Parser) parseExpression() Expression {
	log.Println("Parsing ExpressionValue")
	peeked := p.lexer.PeekSome(2)

	tree := NewExpressionTree()

	switch peeked[0].Typ {
	case lexer.IdentifierToken, lexer.StringToken:
	default:
		p.unexpectedToken(peeked[0])

	}

	for {
		peeked := p.lexer.PeekSome(2)
		switch {
		case peeked[1].Typ == lexer.MethodCallToken:
			if !tree.CanAddLeaf() {
				return tree.GetExpression()
			}
			tree.AddExpression(p.parseMethodCallExpression())

		case peeked[0].Typ == lexer.StringToken:
			if !tree.CanAddLeaf() {
				return tree.GetExpression()
			}
			tree.AddExpression(p.parseStringLiteralExpression())
		case peeked[0].Typ == lexer.AddToken:
			if !tree.CanAddBranch() {
				p.unexpectedToken(peeked[0])
			}
			tree.AddExpression(&BinaryExpression{token: p.lexer.ReadNext()})
		case peeked[0].Typ == lexer.CloseParenToken || peeked[0].Typ == lexer.CloseBraceToken:
			return tree.GetExpression()
		default:
			p.unexpectedToken(peeked[0])
		}
	}
}

func (p *Parser) parseMethodCallExpression() Expression {
	log.Println("Parsing method call ExpressionValue")
	var result MethodCallExpression
	token := p.lexer.Peek()
	if token.Typ != lexer.IdentifierToken {
		log.Panicln("Expected IdentifierToken, got", token.String())
	}

	result.ObjectName = p.parseIdentifier()

	token = p.lexer.ReadNext()
	if token.Typ != lexer.MethodCallToken {
		log.Panicln("Expected MethodCallToken, got", token.String())
	}
	result.token = token
	result.MethodName = p.parseIdentifier()
	result.Arguments = p.parseArguments()

	return &result
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

	valueString, err := token.GetValueString()
	if err != nil {
		log.Panicln("Error parsing string value", err)
	}

	return &StringLiteralExpression{token: token, Value: valueString}
}

func (p *Parser) unexpectedToken(token *lexer.Token) {
	log.Panicln("Unexpected token", token.String())
}

func (p *Parser) unexpectedTokenExpected(expected lexer.TokenType, actual *lexer.Token) {
	log.Panicln("Unexpected token", actual.String(), "expected", expected)
}
