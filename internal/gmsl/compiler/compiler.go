package compiler

import (
	"goMud/internal/gmsl/parser"
	"log"
	"strconv"
)

type Compiler struct {
	ast    *parser.AstNode
	result Assembly
}

func NewCompiler(ast parser.AstNode) *Compiler {
	return &Compiler{
		ast:    &ast,
		result: *newAssembly(),
	}
}

func (c *Compiler) Compile() *Assembly {
	c.processNode(c.ast)
	return &c.result
}

func (c *Compiler) processNode(node *parser.AstNode) {
	switch n := (*node).(type) {
	case *parser.Class:
		c.processClass(n)
	case *parser.FunctionDeclaration:
		c.result.addFunction(c.processFunctionDeclaration(n))
	default:
		log.Panicln("Unknown node type", n.String())
	}
}

func (c *Compiler) processClass(n *parser.Class) {
	for _, f := range n.Functions {
		var a parser.AstNode = &f
		c.processNode(&a)
	}
}

func (c *Compiler) processFunctionDeclaration(n *parser.FunctionDeclaration) *FunctionInfo {
	result := newFunctionInfo(n.Name.Value)
	for _, a := range n.Arguments {
		c.processArgumentDeclaration(&a, result)
	}

	for _, r := range n.ReturnTypes {
		result.addReturnType(typToType[r.Name])
	}

	for _, s := range n.Statements {
		c.processStatement(&s, result)
	}
	result.addEntry(*NewReturnEntry(nil, *n.GetToken()))
	return result
}

var typToType = map[string]Type{
	"string": StringType,
}

func (c *Compiler) processArgumentDeclaration(argumentDeclaration *parser.ArgumentDeclaration, function *FunctionInfo) {
	function.addArgument(argumentDeclaration.Name.Value, typToType[argumentDeclaration.Typ.Name])
	function.addEntry(*NewPopToRegisterEntry(nil, function.getRegisterOf(argumentDeclaration.Name.Value), *argumentDeclaration.GetToken()))
}

func (c *Compiler) processStatement(s *parser.Statement, f *FunctionInfo) {
	switch n := (*s).(type) {
	case *parser.ExpressionStatement:
		c.processExpressionStatement(n, f)
	case *parser.IfStatement:
		c.processIfStatement(n, f)
	case *parser.VariableDeclarationStatement:
		c.processVariableDeclarationStatement(n, f)
	case *parser.VariableAssignmentStatement:
		c.processVariableAssignmentStatement(n, f)
	case *parser.VariableCreateAndAssignStatement:
		c.processVariableCreateAndAssignStatement(n, f)
	case *parser.ReturnStatement:
		c.processReturnStatement(n, f)
	default:
		log.Panicln("Unknown statement type", n.String())
	}
}

func (c *Compiler) processExpressionStatement(statement *parser.ExpressionStatement, f *FunctionInfo) {
	f.addEntries(c.processExpression(&statement.ExpressionValue, f))
}

func isContextName(name string) bool {
	return name == "player" || name == "room" || name == "item"
}

func (c *Compiler) processExpression(expression *parser.Expression, f *FunctionInfo) []AssemblyEntry {
	var result []AssemblyEntry
	switch (*expression).(type) {
	case *parser.MethodCallExpression:
		for _, a := range (*expression).(*parser.MethodCallExpression).Arguments {
			result = append(result, c.processExpression(&a, f)...)
		}
		methodName := (*expression).(*parser.MethodCallExpression).MethodName
		nameIdx := f.addString(methodName.Value)
		result = append(result, *NewPushStringEntry(nil, nameIdx, *methodName.GetToken()))
		objectName := (*expression).(*parser.MethodCallExpression).ObjectName
		objectIdx := f.addString(objectName.Value)
		if isContextName(objectName.Value) {
			result = append(result, *NewPushContextEntry(nil, objectIdx, *objectName.GetToken()))
		} else {
			result = append(result, *NewPushStringEntry(nil, objectIdx, *objectName.GetToken()))
		}
		n := (*expression).(*parser.MethodCallExpression).GetToken()
		result = append(result, *NewCallEntry(nil, *n))
	case *parser.BinaryExpression:
		result = append(result, c.processExpression(&(*expression).(*parser.BinaryExpression).Left, f)...)
		result = append(result, c.processExpression(&(*expression).(*parser.BinaryExpression).Right, f)...)
		result = append(result, *NewOperationEntry(nil, *(*expression).(*parser.BinaryExpression).GetToken()))
	case *parser.StringLiteralExpression:
		stringIdx := f.addString((*expression).(*parser.StringLiteralExpression).Value)
		result = append(result, *NewPushStringEntry(nil, stringIdx, *(*expression).(*parser.StringLiteralExpression).GetToken()))
	case *parser.NumericLiteralExpression:
		e := (*expression).(*parser.NumericLiteralExpression)
		result = append(result, *NewPushNumberEntry(nil, e.GetValue(), *e.GetToken()))
	case *parser.IdentifierExpression:
		result = append(result, c.processIdentifierExpression((*expression).(*parser.IdentifierExpression), f))
	default:
		log.Panicln("Unknown expression type", (*expression).String())
	}

	return result
}

func (c *Compiler) processIfStatement(statement *parser.IfStatement, f *FunctionInfo) {
	// Process the condition expression
	f.addEntries(c.processExpression(&statement.Condition, f))
	jumpLabelName := ".if_jump_" + strconv.Itoa(f.nextEntryPost())
	f.addEntry(*NewJumpIfFalseEntry(nil, jumpLabelName, *statement.GetToken()))

	// Process the statements in the 'if' block
	for _, s := range statement.Statements {
		c.processStatement(&s, f)
	}

	jumpToEndLabelName := ".if_jump_end_" + strconv.Itoa(f.nextEntryPost())
	f.addEntry(*NewJumpEntry(nil, jumpToEndLabelName, *statement.GetToken()))
	f.setNextLabel(&jumpLabelName)

	// Process the statements in the 'else' block, if it exists
	if statement.ElseStatements != nil {
		for _, s := range statement.ElseStatements {
			c.processStatement(&s, f)
		}
	}

	f.setNextLabel(&jumpToEndLabelName)
	f.addEntry(*NewNoOpEntry(nil, *statement.GetToken()))
}

func (c *Compiler) processIdentifierExpression(expression *parser.IdentifierExpression, f *FunctionInfo) AssemblyEntry {
	return *NewPushFromRegisterEntry(nil, f.getRegisterOf(expression.Identifier.Value), *expression.GetToken())
}

func (c *Compiler) processVariableDeclarationStatement(statement *parser.VariableDeclarationStatement, f *FunctionInfo) {
	f.addIdentifier(statement.GetVariableName(), typToType[statement.GetType().Name])
}

func (c *Compiler) processVariableAssignmentStatement(statement *parser.VariableAssignmentStatement, f *FunctionInfo) {
	f.addEntries(c.processExpression(statement.GetExpression(), f))
	f.addEntry(*NewPopToRegisterEntry(nil, f.getRegisterOf(statement.GetVariableName()), *statement.GetToken()))
}

func (c *Compiler) processVariableCreateAndAssignStatement(statement *parser.VariableCreateAndAssignStatement, f *FunctionInfo) {
	f.addIdentifier(statement.GetVariableName(), StringType)
	f.addEntries(c.processExpression(statement.GetExpression(), f))
	f.addEntry(*NewPopToRegisterEntry(nil, f.getRegisterOf(statement.GetVariableName()), *statement.GetToken()))
}

func (c *Compiler) processReturnStatement(statement *parser.ReturnStatement, functionInfo *FunctionInfo) {
	functionInfo.addEntries(c.processExpression(statement.GetValue(), functionInfo))
	functionInfo.addEntry(*NewReturnEntry(nil, *statement.GetToken()))
}
