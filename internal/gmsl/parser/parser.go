package parser

import (
	"goMud/internal/gmsl"
	"log"
)

type Parser struct {
	lexer *gmsl.Lexer
}

func NewParser(l *gmsl.Lexer) *Parser {
	return &Parser{l}
}

func (p *Parser) Parse() *Class {
	return p.parseClass()
}

func (p *Parser) parseClass() *Class {
	log.Println("Parsing class")
	token := p.lexer.ReadNext()
	if token.Typ != gmsl.PackageToken {
		panic("Expected class")
	}

	name := p.parseIdentifier()
	class := Class{token: token, Name: name, Imports: make([]ImportDeclaration, 0)}
	for {
		peeked := p.lexer.Peek()
		switch peeked.Typ {
		case gmsl.ImportToken:
			imports := p.parseImportDeclarations()
			class.Imports = append(class.Imports, imports...)
		case gmsl.FuncToken:
			functions := p.parseFunctionDeclarations()
			class.Functions = append(class.Functions, functions...)
		case gmsl.EofToken:
			return &class
		default:
			log.Panicln("Unexpected token", peeked.String())
		}
	}
}

func (p *Parser) parseIdentifier() *Identifier {
	log.Println("Parsing identifier")
	token := p.lexer.ReadNext()
	if token.Typ != gmsl.IdentifierToken {
		log.Panicln("Expected identifier, got", token.String())
	}

	return &Identifier{token: token, Value: token.Value}
}

func (p *Parser) parseImportDeclarations() []ImportDeclaration {
	log.Println("Parsing import declarations")
	token := p.lexer.Peek()
	tokenDeclarations := make([]ImportDeclaration, 0)
	if token.Typ == gmsl.ImportToken {
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
	case gmsl.IdentifierToken:
		return p.parseSingleImportDeclaration()
	case gmsl.OpenParenToken:
		return p.parseImportDeclarationList()
	default:
		log.Panicln("Unexpected token", tokens[1].String())
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
	if skip.Typ != gmsl.OpenParenToken {
		log.Panicln("Expected OpenParenToken, got", skip.String())
	}
	for {
		token := p.lexer.Peek()
		if token.Typ == gmsl.CloseParenToken {
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
	if token.Typ != gmsl.StringToken {
		log.Panicln("Expected string value, got", token.String())
	}

	return &Identifier{token: token, Value: token.Value}

}

func (p *Parser) parseFunctionDeclarations() []FunctionDeclaration {
	log.Println("Parsing function declarations")
	var functions []FunctionDeclaration
	for {
		token := p.lexer.Peek()
		if token.Typ == gmsl.FuncToken {
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
	if token.Typ != gmsl.FuncToken {
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
	if token.Typ != gmsl.OpenParenToken {
		log.Panicln("Expected OpenParenToken, got", token.String())
	}

	arguments := make([]ArgumentDeclaration, 0)
	for {
		token := p.lexer.Peek()
		if token.Typ == gmsl.CloseParenToken {
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
	if token.Typ != gmsl.TypeToken {
		log.Panicln("Expected TypeToken, got", token.String())
	}

	return &Type{token: token, Name: token.Value}
}

func (p *Parser) parseStatements() []Statement {
	log.Println("Parsing statements")
	token := p.lexer.ReadNext()
	if token.Typ != gmsl.OpenBraceToken {
		log.Panicln("Expected OpenBraceToken, got", token.String())
	}
	statements := make([]Statement, 0)
	for {
		token := p.lexer.Peek()
		if token.Typ == gmsl.CloseBraceToken {
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

	if peeked[0].Typ == gmsl.IdentifierToken {
		switch peeked[1].Typ {
		case gmsl.MethodCallToken:
			return p.parseExpressionStatement()
		default:
			log.Panicln("Unexpected token", peeked[1].String())
		}
	} else {
		log.Panicln("Unexpected token", peeked[0].String())
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
	var expression Expression

	switch peeked[0].Typ {
	case gmsl.IdentifierToken, gmsl.StringToken:
	default:
		log.Panicln("Unexpected token", peeked[0].String())

	}

	for {
		peeked := p.lexer.PeekSome(2)
		switch {
		case peeked[1].Typ == gmsl.MethodCallToken:
			if expression == nil {
				expression = p.parseMethodCallExpression()
			} else {
				switch expression.(type) {
				case *BinaryExpression:
					if expression.(*BinaryExpression).Right == nil {
						expression.(*BinaryExpression).Right = p.parseMethodCallExpression()
					} else {
						return expression
					}
				}
			}

		case peeked[0].Typ == gmsl.StringToken:
			if expression == nil {
				expression = p.parseStringLiteralExpression()
			} else {
				switch expression.(type) {
				case *BinaryExpression:
					if expression.(*BinaryExpression).Right == nil {
						expression.(*BinaryExpression).Right = p.parseStringLiteralExpression()
					} else {
						return expression
					}
				}
			}
		case peeked[0].Typ == gmsl.AddToken:
			token := p.lexer.ReadNext()
			if expression == nil {
				log.Panicln("Unexpected token", peeked[0].String())
			} else {
				expression = &BinaryExpression{token: token, Left: expression}
			}
		case peeked[0].Typ == gmsl.CloseParenToken:
			return expression
		case peeked[0].Typ == gmsl.SemicolonToken:
			p.lexer.ReadNext()
			return expression
		default:
			log.Panicln("Unexpected token", peeked[1].String())
		}
	}

	return expression
}

func (p *Parser) parseMethodCallExpression() Expression {
	log.Println("Parsing method call ExpressionValue")
	var result MethodCallExpression
	token := p.lexer.Peek()
	if token.Typ != gmsl.IdentifierToken {
		log.Panicln("Expected IdentifierToken, got", token.String())
	}

	result.ObjectName = p.parseIdentifier()

	token = p.lexer.ReadNext()
	if token.Typ != gmsl.MethodCallToken {
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
	if token.Typ != gmsl.OpenParenToken {
		log.Panicln("Expected OpenParenToken, got", token.String())
	}

	arguments := make([]Expression, 0)
	for {
		token := p.lexer.Peek()
		if token.Typ == gmsl.CloseParenToken {
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
	if token.Typ != gmsl.StringToken {
		log.Panicln("Expected StringToken, got", token.String())
	}

	return &StringLiteralExpression{token: token, Value: token.Value}
}
