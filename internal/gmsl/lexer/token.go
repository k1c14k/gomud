package lexer

import "bytes"

type TokenType int

// TokenType constants represent the various lexical tokens recognized by the parser.
const (
	// InvalidToken indicates a meaningless or unrecognized sequence of characters
	InvalidToken TokenType = iota
	// EofToken signifies the end of the input file or string
	EofToken
	// ImportToken represents the "import" keyword
	ImportToken
	// FuncToken represents the "func" keyword
	FuncToken
	// IdentifierToken represents a user-defined name for variables, functions, etc.
	IdentifierToken
	// ContextToken represents a special context name (e.g., player, room, item)
	ContextToken
	// OpenParenToken represents an opening parenthesis "("
	OpenParenToken
	// CloseParenToken represents a closing parenthesis ")"
	CloseParenToken
	// OpenBraceToken represents an opening brace "{"
	OpenBraceToken
	// CloseBraceToken represents a closing brace "}"
	CloseBraceToken
	// StringToken represents a string literal enclosed in double quotes
	StringToken
	// NumericToken represents a numeric literal
	NumericToken
	// AddToken represents the addition operator "+"
	AddToken
	// SubtractToken represents the subtraction operator "-"
	SubtractToken
	// MultiplyToken represents the multiplication operator "*"
	MultiplyToken
	// DivideToken represents the division operator "/"
	DivideToken
	// ModuloToken represents the modulo operator "%"
	ModuloToken
	// MethodCallToken represents the method call operator "."
	MethodCallToken
	// TypeToken represents a built-in type name (e.g., int, string)
	TypeToken
	// IfToken represents the "if" keyword
	IfToken
	// ElseToken represents the "else" keyword
	ElseToken
	// EqualToken represents the equality operator "=="
	EqualToken
	// AssignToken represents the assignment operator "="
	AssignToken
	// CreateAndAssignToken represents the short variable declaration operator ":="
	CreateAndAssignToken
	// VarToken represents the "var" keyword
	VarToken
	// ReturnToken represents the "return" keyword
	ReturnToken
)

var tokenNames = map[TokenType]string{
	InvalidToken:         "InvalidToken",
	EofToken:             "EofToken",
	ImportToken:          "ImportToken",
	FuncToken:            "FuncToken",
	IdentifierToken:      "IdentifierToken",
	ContextToken:         "ContextToken",
	OpenParenToken:       "OpenParenToken",
	CloseParenToken:      "CloseParenToken",
	OpenBraceToken:       "OpenBraceToken",
	CloseBraceToken:      "CloseBraceToken",
	StringToken:          "StringToken",
	NumericToken:         "NumericToken",
	AddToken:             "AddToken",
	SubtractToken:        "SubtractToken",
	MultiplyToken:        "MultiplyToken",
	DivideToken:          "DivideToken",
	ModuloToken:          "ModuloToken",
	MethodCallToken:      "MethodCallToken",
	TypeToken:            "TypeToken",
	IfToken:              "IfToken",
	ElseToken:            "ElseToken",
	EqualToken:           "EqualToken",
	AssignToken:          "AssignToken",
	CreateAndAssignToken: "CreateAndAssignToken",
	VarToken:             "VarToken",
	ReturnToken:          "ReturnToken",
}

func (t TokenType) String() string {
	return tokenNames[t]
}

type Token struct {
	Typ      TokenType
	rawValue string
}

func (t *Token) String() string {
	return tokenNames[t.Typ] + " " + t.rawValue
}

func (t *Token) GetRawValue() string {
	return t.rawValue
}

func (t *Token) GetValueString() (string, error) {
	buffer := bytes.NewBufferString("")
	reader := bytes.NewReader([]byte(t.rawValue))
	for {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}
		switch b {
		case '\\':
			b, err = reader.ReadByte()
			if err != nil {
				return "", err
			}
			switch b {
			case 'n':
				buffer.WriteByte('\r')
				buffer.WriteByte('\n')
			case 't':
				buffer.WriteByte('\t')
			case 'r':
				// ignore
			default:
				buffer.WriteByte(b)
			}
		default:
			buffer.WriteByte(b)
		}
	}
	return buffer.String(), nil
}
