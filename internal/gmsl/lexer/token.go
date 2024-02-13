package lexer

import "bytes"

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
	IfToken
	ElseToken
	EqualToken
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
	IfToken:         "IfToken",
	ElseToken:       "ElseToken",
	EqualToken:      "EqualToken",
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
