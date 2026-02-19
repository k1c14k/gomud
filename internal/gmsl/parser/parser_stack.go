package parser

import (
	"goMud/internal/gmsl/lexer"
	"log"
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

type Parser struct {
	lexer     *lexer.Lexer
	stack     []*Frame
	nodeStack []AstNode
}

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
	log.Panicln("Unexpected token", token.String())
}

func (p *Parser) unexpectedTokenExpected(expected lexer.TokenType, actual *lexer.Token) {
	if actual.Typ != expected {
		log.Panicln("Unexpected token", actual.String(), "expected", expected)
	}
}

func (p *Parser) expect(expected lexer.TokenType, name string) *lexer.Token {
	read := p.lexer.ReadNext()
	if expected != read.Typ {
		log.Panicln("Expected", name, "got", read.String())
	}
	return read
}

func (p *Parser) parseIdentifier() *Identifier {
	log.Println("Parsing identifier")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.IdentifierToken {
		log.Panicln("Expected identifier, got", token.String())
	}
	return newIdentifier(token)
}

func (p *Parser) parseStringValue() *Identifier {
	log.Println("Parsing string value")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.StringToken {
		log.Panicln("Expected string value, got", token.String())
	}
	return newIdentifier(token)
}

func (p *Parser) parseType() *Type {
	log.Println("Parsing type")
	token := p.lexer.ReadNext()
	if token.Typ != lexer.TypeToken {
		log.Panicln("Expected TypeToken, got", token.String())
	}
	return newType(token)
}

// Parse orchestrates the state machine loop
func (p *Parser) Parse() *Class {
	p.push(&Frame{State: StateClass})

	var resultClass *Class

	for len(p.stack) > 0 {
		frame := p.pop()

		switch frame.State {

		case StateClass:
			log.Println("Parsing class")
			token := p.lexer.ReadNext()
			if token.Typ != lexer.PackageToken {
				panic("Expected class")
			}
			name := p.parseIdentifier()
			frame.Class = newClass(name, token)

			frame.State = StateClassBody
			p.push(frame)

		case StateClassBody:
			// Absorb returned nodes
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
					log.Panicln("Unexpected node in StateClassBody:", n)
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
				resultClass = frame.Class
			default:
				p.unexpectedToken(peeked)
			}

		case StateImportDecl:
			log.Println("Parsing import declarations")
			token := p.lexer.Peek()
			if token.Typ == lexer.ImportToken {
				log.Println("Parsing import declaration")
				tokens := p.lexer.PeekSome(2)
				if len(tokens) < 2 {
					log.Panicln("Expected import declaration")
				}
				switch tokens[1].Typ {
				case lexer.IdentifierToken:
					log.Println("Parsing single import declaration")
					t2 := p.lexer.ReadNext()
					name := p.parseStringValue()
					p.pushNode(newSingleImportDeclaration(name, t2))
				case lexer.OpenParenToken:
					log.Println("Parsing import declaration list")
					t2 := p.lexer.ReadNext()
					skip := p.lexer.ReadNext()
					if skip.Typ != lexer.OpenParenToken {
						p.unexpectedTokenExpected(lexer.OpenParenToken, skip)
					}
					// Launch list loop
					newFrame := &Frame{State: StateImportDeclList, Token: t2}
					p.push(newFrame)
				default:
					p.unexpectedToken(tokens[1])
				}
			}

		case StateImportDeclList:
			token := p.lexer.Peek()
			if token.Typ == lexer.CloseParenToken {
				p.lexer.ReadNext()
				p.pushNode(newImportDeclarationList(&frame.Imports, frame.Token))
			} else {
				frame.Imports = append(frame.Imports, *p.parseStringValue())
				p.push(frame) // loop
			}

		case StateFuncDecl:
			log.Println("Parsing function declaration")
			token := p.lexer.ReadNext()
			if token.Typ != lexer.FuncToken {
				log.Panicln("Expected FuncToken, got", token.String())
			}
			name := p.parseIdentifier()
			frame.FuncDecl = &FunctionDeclaration{
				token:       token,
				Name:        *name,
				ReturnTypes: make([]Type, 0),
			}

			log.Println("Parsing arguments")
			t2 := p.lexer.ReadNext()
			if t2.Typ != lexer.OpenParenToken {
				p.unexpectedTokenExpected(lexer.OpenParenToken, t2)
			}

			frame.State = StateFuncDeclArgs
			p.push(frame)

		case StateFuncDeclArgs:
			token := p.lexer.Peek()
			if token.Typ == lexer.CloseParenToken {
				p.lexer.ReadNext()
				frame.State = StateFuncDeclReturn
				p.push(frame)
			} else {
				log.Println("Parsing argument")
				argName := p.parseIdentifier()
				argType := p.parseType()
				arg := newArgumentDeclaration(argName, argType, argName.token)
				frame.FuncDecl.Arguments = append(frame.FuncDecl.Arguments, *arg)
				p.push(frame) // loop
			}

		case StateFuncDeclReturn:
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

		case StateStatements:
			if frame.FuncDecl != nil {
				// From Function Declaration
				log.Println("Parsing statements")
				token := p.lexer.ReadNext()
				if token.Typ != lexer.OpenBraceToken {
					p.unexpectedTokenExpected(lexer.OpenBraceToken, token)
				}
				frame.State = StateStatementsLoop
				p.push(frame)
			} else if frame.IfStmt != nil {
				// From If Statement (already consumed { token inside StateIfStatement)
				frame.State = StateStatementsLoop
				p.push(frame)
			} else {
				// General use statements list
				log.Println("Parsing statements")
				token := p.lexer.ReadNext()
				if token.Typ != lexer.OpenBraceToken {
					p.unexpectedTokenExpected(lexer.OpenBraceToken, token)
				}
				frame.State = StateStatementsLoop
				p.push(frame)
			}

		case StateStatementsLoop:
			// Absorb returned statement
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
					// Check if we were doing the primary statements or the else statements
					if frame.IfStmt.Statements == nil {
						frame.IfStmt.Statements = frame.Stmts
						frame.Stmts = nil // clear for potential 'else'

						token = p.lexer.Peek()
						if token.Typ == lexer.ElseToken {
							p.lexer.ReadNext() // consume else
							token = p.lexer.Peek()
							if token.Typ != lexer.OpenBraceToken {
								p.unexpectedTokenExpected(lexer.OpenBraceToken, token)
							}
							p.lexer.ReadNext() // consume open brace for else

							// Continue loop for else body
							p.push(frame)
						} else {
							// Done with If
							p.pushNode(frame.IfStmt)
						}
					} else {
						// We already had primary statements, this must be else statements ending
						frame.IfStmt.ElseStatements = frame.Stmts
						p.pushNode(frame.IfStmt)
					}
				} else {
					// We were just parsing a generic block (not currently used directly in standard AST except IF/FUNC, but maybe return a NodeList later. However in goMud Ast, slices are passed natively. Wait, the original code returns []Statement.
					// Since Go cannot pass slice to interface AstNode cleanly, we must attach it to the parent logic!
					// But we only ever call `parseStatements()` from func or if. So we embedded those loops inside their frames!)
				}
			} else {
				p.push(frame) // suspend loop
				p.push(&Frame{State: StateStatement})
			}

		case StateStatement:
			log.Println("Parsing statement")
			peeked := p.lexer.PeekSome(2)
			switch peeked[0].Typ {
			case lexer.VarToken:
				p.push(&Frame{State: StateVariableDeclStmt})
			case lexer.IdentifierToken:
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

		case StateVariableDeclStmt:
			log.Println("Parsing variable declaration statement")
			token := p.lexer.ReadNext()
			if token.Typ != lexer.VarToken {
				log.Panicln("Expected VarToken, got", token.String())
			}
			name := p.parseIdentifier()
			typ := p.parseType()
			p.pushNode(newVariableDeclarationStatement(name, typ, token))

		case StateVariableAssignStmt:
			if frame.Identifier == nil {
				log.Println("Parsing variable assignment statement")
				token := p.lexer.Peek()
				p.unexpectedTokenExpected(lexer.IdentifierToken, token)

				name := p.parseIdentifier()
				assignToken := p.expect(lexer.AssignToken, "AssignToken")
				frame.Identifier = name
				frame.Token = assignToken

				// Parse inner Expression
				p.push(frame)
				p.push(&Frame{State: StateExpression})
			} else {
				// Expression returned
				expr := p.popNode().(Expression)
				p.pushNode(newVariableAssignmentStatement(frame.Identifier, &expr, frame.Token))
			}

		case StateVariableCreateAssignStmt:
			if frame.Identifier == nil {
				log.Println("Parsing variable create and assign statement")
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

		case StateExpressionStmt:
			if frame.ExprTree == nil {
				// First pass
				p.push(frame)
				p.push(&Frame{State: StateExpression})
				// We attach arbitrary non-nil object to recognize second pass
				frame.ExprTree = NewExpressionTree() // just a marker
			} else {
				expr := p.popNode().(Expression)
				p.pushNode(newExpressionStatement(&expr, frame.Token))
			}

		case StateIfStatement:
			if frame.IfStmt == nil {
				log.Println("Parsing if statement")
				token := p.lexer.ReadNext()
				if token.Typ != lexer.IfToken {
					log.Panicln("Expected IfToken, got", token.String())
				}
				frame.IfStmt = &IfStatement{token: token}

				p.push(frame)
				p.push(&Frame{State: StateExpression})
			} else {
				// Expression parsed
				expr := p.popNode().(Expression)
				frame.IfStmt.Condition = expr

				token := p.lexer.Peek()
				if token.Typ != lexer.OpenBraceToken {
					p.unexpectedTokenExpected(lexer.OpenBraceToken, token)
				}
				p.lexer.ReadNext() // consume {

				frame.State = StateStatementsLoop
				p.push(frame)
			}

		case StateReturnStmt:
			if frame.ReturnStmt == nil {
				log.Println("Parsing return statement")
				token := p.expect(lexer.ReturnToken, "ReturnToken")
				frame.Token = token
				frame.ReturnStmt = &ReturnStatement{} // dummy to track pass

				p.push(frame)
				p.push(&Frame{State: StateExpression})
			} else {
				expr := p.popNode().(Expression)
				p.pushNode(newReturnStatement(&expr, frame.Token))
			}

		case StateExpression:
			log.Println("Parsing ExpressionValue")
			peeked := p.lexer.PeekSome(2)
			switch peeked[0].Typ {
			case lexer.IdentifierToken, lexer.StringToken, lexer.NumericToken:
			default:
				p.unexpectedToken(peeked[0])
			}

			frame.ExprTree = NewExpressionTree()
			frame.State = StateExpressionLoop
			p.push(frame)

		case StateExpressionLoop:
			// Apply completed inner component
			if len(p.nodeStack) > 0 {
				expr := p.popNode().(Expression)
				frame.ExprTree.AddExpression(expr)
			}

			peeked := p.lexer.PeekSome(2)

			// TryAdd logic extracted into sequential checks
			if len(peeked) > 1 && peeked[1].Typ == lexer.MethodCallToken {
				if !frame.ExprTree.CanAddLeaf() {
					p.pushNode(frame.ExprTree.GetExpression())
					continue
				}
				p.push(frame)
				p.push(&Frame{State: StateMethodCall})
				continue
			}

			if peeked[0].Typ == lexer.StringToken {
				if !frame.ExprTree.CanAddLeaf() {
					p.pushNode(frame.ExprTree.GetExpression())
					continue
				}
				log.Println("Parsing string literal ExpressionValue")
				token := p.lexer.ReadNext()
				if token.Typ != lexer.StringToken {
					log.Panicln("Expected StringToken, got", token.String())
				}
				frame.ExprTree.AddExpression(newStringLiteralExpression(token))
				p.push(frame)
				continue
			}

			if peeked[0].Typ == lexer.NumericToken {
				if !frame.ExprTree.CanAddLeaf() {
					p.pushNode(frame.ExprTree.GetExpression())
					continue
				}
				log.Println("Parsing numeric literal ExpressionValue")
				token := p.lexer.ReadNext()
				p.unexpectedTokenExpected(lexer.NumericToken, token)
				frame.ExprTree.AddExpression(newNumericLiteralExpression(token))
				p.push(frame)
				continue
			}

			if peeked[0].Typ == lexer.IdentifierToken {
				if !frame.ExprTree.CanAddLeaf() {
					p.pushNode(frame.ExprTree.GetExpression())
					continue
				}
				log.Println("Parsing identifier ExpressionValue")
				token := p.lexer.Peek()
				p.unexpectedTokenExpected(lexer.IdentifierToken, token)
				id := p.parseIdentifier()
				frame.ExprTree.AddExpression(newIdentifierExpression(id, token))
				p.push(frame)
				continue
			}

			if peeked[0].Typ == lexer.AddToken || peeked[0].Typ == lexer.EqualToken || peeked[0].Typ == lexer.DivideToken || peeked[0].Typ == lexer.MultiplyToken || peeked[0].Typ == lexer.SubtractToken || peeked[0].Typ == lexer.ModuloToken {
				if !frame.ExprTree.CanAddBranch() {
					p.unexpectedToken(peeked[0])
				}
				frame.ExprTree.AddExpression(newBinaryExpression(p.lexer.ReadNext()))
				p.push(frame)
				continue
			}

			if !frame.ExprTree.CanAddLeaf() {
				p.pushNode(frame.ExprTree.GetExpression())
				continue
			}

			p.unexpectedToken(peeked[0])

		case StateMethodCall:
			if frame.Identifier == nil {
				log.Println("Parsing method call ExpressionValue")
				token := p.lexer.Peek()
				p.unexpectedTokenExpected(lexer.IdentifierToken, token)

				objectName := p.parseIdentifier()

				mtoken := p.lexer.ReadNext()
				if mtoken.Typ != lexer.MethodCallToken {
					log.Panicln("Expected MethodCallToken, got", mtoken.String())
				}
				methodName := p.parseIdentifier()

				frame.Identifier = objectName
				frame.MethodCall = &MethodCallExpression{
					token:      mtoken,
					ObjectName: *objectName,
					MethodName: *methodName,
					Arguments:  make([]Expression, 0),
				}

				log.Println("Parsing arguments")
				t2 := p.lexer.ReadNext()
				if t2.Typ != lexer.OpenParenToken {
					p.unexpectedTokenExpected(lexer.OpenParenToken, t2)
				}

				frame.State = StateMethodCallArgs
				p.push(frame)
			}

		case StateMethodCallArgs:
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
	}

	return resultClass
}
