package parser

import (
	"goMud/internal/gmsl/lexer"
	"log"
)

func newClass(name *Identifier, token *lexer.Token) *Class {
	return &Class{token: token, Name: *name, Imports: make([]ImportDeclaration, 0)}
}

func newSingleImportDeclaration(name *Identifier, token *lexer.Token) *SingleImportDeclaration {
	return &SingleImportDeclaration{token: token, Name: *name}
}

func newImportDeclarationList(imports *[]Identifier, token *lexer.Token) *ImportDeclarationList {
	return &ImportDeclarationList{token: token, Imports: *imports}
}

func newFunctionDeclaration(name *Identifier, args *[]ArgumentDeclaration, statements *[]Statement, token *lexer.Token) *FunctionDeclaration {
	return &FunctionDeclaration{token: token, Name: *name, Arguments: *args, Statements: *statements}
}

func newArgumentDeclaration(name *Identifier, typ *Type, token *lexer.Token) *ArgumentDeclaration {
	return &ArgumentDeclaration{token: token, Name: *name, Typ: *typ}
}

func newMethodCallExpression(objectName *Identifier, methodName *Identifier, args *[]Expression, token *lexer.Token) *MethodCallExpression {
	return &MethodCallExpression{token: token, ObjectName: *objectName, Arguments: *args, MethodName: *methodName}
}

func newIdentifierExpression(name *Identifier, token *lexer.Token) *IdentifierExpression {
	return &IdentifierExpression{token: token, Identifier: *name}
}

func newIdentifier(token *lexer.Token) *Identifier {
	return &Identifier{token: token, Value: token.GetRawValue()}
}

func newType(token *lexer.Token) *Type {
	return &Type{token: token, Name: token.GetRawValue()}
}

func newExpressionStatement(expression *Expression, token *lexer.Token) *ExpressionStatement {
	return &ExpressionStatement{token: token, ExpressionValue: *expression}
}

func newBinaryExpression(token *lexer.Token) *BinaryExpression {
	return &BinaryExpression{token: token}
}

func newStringLiteralExpression(token *lexer.Token) *StringLiteralExpression {
	valueString, err := token.GetValueString()
	if err != nil {
		log.Panicln("Error parsing string value", err)
	}
	return &StringLiteralExpression{token: token, Value: valueString}
}

func newIfStatement(condition *Expression, consequence *[]Statement, alternative *[]Statement, token *lexer.Token) *IfStatement {
	return &IfStatement{token: token, Condition: *condition, Statements: *consequence, ElseStatements: *alternative}
}

func newVariableAssignmentStatement(name *Identifier, expression *Expression, token *lexer.Token) *VariableAssignmentStatement {
	return &VariableAssignmentStatement{
		token: token,
		name:  *name,
		value: *expression,
	}
}

func newVariableCreateAndAssignStatement(name *Identifier, expression *Expression, token *lexer.Token) *VariableCreateAndAssignStatement {
	return &VariableCreateAndAssignStatement{
		token: token,
		name:  *name,
		value: *expression,
	}
}

func newVariableDeclarationStatement(name *Identifier, typ *Type, token *lexer.Token) *VariableDeclarationStatement {
	return &VariableDeclarationStatement{
		token: token,
		name:  *name,
		typ:   *typ,
	}
}

func newNumericLiteralExpression(token *lexer.Token) *NumericLiteralExpression {
	return &NumericLiteralExpression{token: token}
}
