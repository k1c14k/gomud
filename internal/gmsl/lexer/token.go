package lexer

type TokenType int

const (
	InvalidToken TokenType = iota
	EofToken
	PackageToken
	ImportToken
	FuncToken
	IdentifierToken
	OpenParenToken
	CloseParenToken
	OpenBraceToken
	CloseBraceToken
	StringToken
	AddToken
	MethodCallToken
	TypeToken
)

var tokenNames = map[TokenType]string{
	InvalidToken:    "InvalidToken",
	EofToken:        "EofToken",
	PackageToken:    "PackageToken",
	ImportToken:     "ImportToken",
	FuncToken:       "FuncToken",
	IdentifierToken: "IdentifierToken",
	OpenParenToken:  "OpenParenToken",
	CloseParenToken: "CloseParenToken",
	OpenBraceToken:  "OpenBraceToken",
	CloseBraceToken: "CloseBraceToken",
	StringToken:     "StringToken",
	AddToken:        "AddToken",
	MethodCallToken: "MethodCallToken",
	TypeToken:       "TypeToken",
}

func (t TokenType) String() string {
	return tokenNames[t]
}

type Token struct {
	Typ   TokenType
	Value string
}

func (t *Token) String() string {
	return tokenNames[t.Typ] + " " + t.Value
}
