package parser

import "bytes"

func (c *Class) PrettyPrint(tabs int) string {
	var buffer bytes.Buffer
	buffer.WriteString("# class ")
	buffer.WriteString(c.Name.String())
	buffer.WriteString("\n\n")
	for _, i := range c.Imports {
		buffer.WriteString(i.PrettyPrint(tabs))
		buffer.WriteString("\n")
	}
	for _, f := range c.Functions {
		buffer.WriteString(f.PrettyPrint(tabs))
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (s *SingleImportDeclaration) PrettyPrint(_ int) string {
	var buffer bytes.Buffer
	buffer.WriteString("import \"")
	buffer.WriteString(s.Name.String())
	buffer.WriteString("\"\n")
	return buffer.String()
}

func (i *ImportDeclarationList) PrettyPrint(_ int) string {
	var buffer bytes.Buffer
	buffer.WriteString("import (\n")
	for _, i := range i.Imports {
		buffer.WriteString("\"")
		buffer.WriteString(i.String())
		buffer.WriteString("\"\n")
	}
	buffer.WriteString(")\n")
	return buffer.String()
}

func (f *FunctionDeclaration) PrettyPrint(tabs int) string {
	var buffer bytes.Buffer
	buffer.WriteString("func ")
	buffer.WriteString(f.Name.String())
	buffer.WriteString(" (")
	for _, a := range f.Arguments {
		buffer.WriteString(a.PrettyPrint(tabs))
	}
	buffer.WriteString(") {\n")
	for _, s := range f.Statements {
		buffer.WriteString(s.PrettyPrint(tabs + 1))
	}
	buffer.WriteString("}\n")
	return buffer.String()
}

func (a *ArgumentDeclaration) PrettyPrint(_ int) string {
	var buffer bytes.Buffer
	buffer.WriteString(a.Name.String())
	buffer.WriteString(" ")
	buffer.WriteString(a.Typ.String())
	return buffer.String()
}

func (e *ExpressionStatement) PrettyPrint(tabs int) string {
	var buffer bytes.Buffer
	for i := 0; i < tabs; i++ {
		buffer.WriteString("\t")
	}
	buffer.WriteString(e.ExpressionValue.PrettyPrint(0))
	buffer.WriteString("\n")
	return buffer.String()
}

func (b *BinaryExpression) PrettyPrint(_ int) string {
	return b.Left.PrettyPrint(0) + " " + b.token.GetRawValue() + " " + b.Right.PrettyPrint(0)
}

func (m *MethodCallExpression) PrettyPrint(_ int) string {
	var buffer bytes.Buffer
	buffer.WriteString(m.ObjectName.String())
	buffer.WriteString(".")
	buffer.WriteString(m.MethodName.String())
	buffer.WriteString("(")
	for i, a := range m.Arguments {
		buffer.WriteString(a.PrettyPrint(0))
		if i < len(m.Arguments)-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString(")")
	return buffer.String()
}

func (s *StringLiteralExpression) PrettyPrint(_ int) string {
	var buffer bytes.Buffer
	buffer.WriteString("\"")
	buffer.WriteString(s.Value)
	buffer.WriteString("\"")
	return buffer.String()
}

func (i *IfStatement) PrettyPrint(tabs int) string {
	var buffer bytes.Buffer
	for i := 0; i < tabs; i++ {
		buffer.WriteString("\t")
	}
	buffer.WriteString("if ")
	buffer.WriteString(i.Condition.PrettyPrint(0))
	buffer.WriteString(" {\n")
	for _, s := range i.Statements {
		buffer.WriteString(s.PrettyPrint(tabs + 1))
	}
	for i := 0; i < tabs; i++ {
		buffer.WriteString("\t")
	}
	buffer.WriteString("} else {\n")
	for _, s := range i.ElseStatements {
		buffer.WriteString(s.PrettyPrint(tabs + 1))
	}
	for i := 0; i < tabs; i++ {
		buffer.WriteString("\t")
	}
	buffer.WriteString("}\n")
	return buffer.String()
}

func (i *IdentifierExpression) PrettyPrint(_ int) string {
	return i.Identifier.String()
}
