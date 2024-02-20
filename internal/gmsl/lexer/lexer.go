package lexer

import (
	"log"
	"strings"
)

type State func(*Lexer) State

type Lexer struct {
	input  string
	start  int
	pos    int
	tokens chan Token
	state  State
	peeked []*Token
}

func (l *Lexer) run() {
	state := defaultState
	for state != nil {
		state = state(l)
	}
	close(l.tokens)
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		tokens: make(chan Token, 2),
		state:  defaultState,
	}
	return l
}

func (l *Lexer) nextToken() *Token {
	for {
		select {
		case t := <-l.tokens:
			return &t
		default:
			l.state = l.state(l)
		}
	}
}

func (l *Lexer) ReadNext() *Token {
	switch {
	case len(l.peeked) > 0:
		t := l.peeked[0]
		l.peeked = l.peeked[1:]
		log.Println("Peeked token", t)
		return t
	default:
		token := l.nextToken()
		log.Println("Read token", token)
		return token
	}
}

var keywords = map[string]TokenType{
	"package": PackageToken,
	"import":  ImportToken,
	"func":    FuncToken,
	"if":      IfToken,
	"else":    ElseToken,
	"var":     VarToken,
}

func (l *Lexer) hasPrefix(m map[string]TokenType) bool {
	for k := range m {
		if strings.HasPrefix(l.input[l.pos:], k+" ") {
			return true
		}
	}
	return false
}

func (l *Lexer) nextRunes(n int) string {
	if l.pos+n > len(l.input) {
		return l.input[l.pos:]
	}
	return l.input[l.pos : l.pos+n]
}

var parenthesis = map[string]TokenType{
	"(": OpenParenToken,
	")": CloseParenToken,
	"{": OpenBraceToken,
	"}": CloseBraceToken,
}

func isParenthesis(r rune) bool {
	for k := range parenthesis {
		if rune(k[0]) == r {
			return true
		}
	}
	return false
}

func (l *Lexer) isParenthesis() bool {
	for k := range parenthesis {
		if strings.HasPrefix(l.input[l.pos:], k) {
			return true
		}
	}
	return false
}

var operator = map[string]TokenType{
	"+":  AddToken,
	".":  MethodCallToken,
	"==": EqualToken,
	"=":  AssignToken,
	":=": CreateAndAssignToken,
}

func isOperator(r rune) bool {
	for k := range operator {
		if rune(k[0]) == r {
			return true
		}
	}
	return false
}

func (l *Lexer) isOperator() bool {
	for k := range operator {
		if strings.HasPrefix(l.input[l.pos:], k) {
			return true
		}
	}
	return false
}

var types = [...]string{"int", "float", "string", "map", "object", "function", "mixed"}
var validIdentifier = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789"

func (l *Lexer) isType() bool {
	for _, t := range types {
		if strings.HasPrefix(l.input[l.pos:], t) {
			return !strings.ContainsAny(l.input[l.pos+len(t):l.pos+len(t)+1], validIdentifier)
		}
	}
	return false
}

func (l *Lexer) PeekSome(n int) []*Token {
	if len(l.peeked) >= n {
		return l.peeked[:n]
	}

	for i := len(l.peeked); i < n; i++ {
		l.peeked = append(l.peeked, l.nextToken())
	}

	return l.peeked
}

func (l *Lexer) Peek() *Token {
	return l.PeekSome(1)[0]
}

func (l *Lexer) invalidToken() {
	l.tokens <- Token{InvalidToken, "Invalid token near " + l.nextRunes(20)}
}

func defaultState(l *Lexer) State {
whitespaces:
	for {
		if l.pos >= len(l.input) {
			l.tokens <- Token{EofToken, ""}
			return nil
		}
		switch l.input[l.pos] {
		case ' ', '\t', '\n', '\r':
			l.pos++
			l.start++
		default:
			break whitespaces
		}
	}

	switch {
	case l.hasPrefix(keywords):
		return keywordState
	case l.isParenthesis():
		return parenthesisState
	case l.isOperator():
		return operatorState
	case l.input[l.pos] == '"':
		l.pos++
		l.start++
		return stringState
	case l.isType():
		return typeState
	default:
		return identifierState
	}
}

func keywordState(l *Lexer) State {
	for k, v := range keywords {
		if strings.HasPrefix(l.input[l.pos:], k+" ") {
			l.pos += len(k) + 1
			l.start = l.pos
			l.tokens <- Token{v, k}
			return defaultState
		}
	}
	l.invalidToken()
	return nil
}

func identifierState(l *Lexer) State {
	for {
		if l.pos >= len(l.input) {
			l.tokens <- Token{EofToken, ""}
			return nil
		}

		if isParenthesis(rune(l.input[l.pos])) || isOperator(rune(l.input[l.pos])) {
			l.tokens <- Token{IdentifierToken, l.input[l.start:l.pos]}
			l.start = l.pos
			return defaultState
		}

		switch l.input[l.pos] {
		case ' ', '\t', '\n', '\r':
			l.tokens <- Token{IdentifierToken, l.input[l.start:l.pos]}
			l.start = l.pos
			return defaultState
		default:
			l.pos++
		}
	}
}

func parenthesisState(l *Lexer) State {
	for k, v := range parenthesis {
		if strings.HasPrefix(l.input[l.pos:], k) {
			l.pos += len(k)
			l.start = l.pos
			l.tokens <- Token{v, k}
			return defaultState
		}
	}
	l.invalidToken()
	return nil
}

func stringState(l *Lexer) State {
	lastChar := rune(0)
	for {
		if l.pos >= len(l.input) {
			l.tokens <- Token{EofToken, ""}
			return nil
		}
		switch l.input[l.pos] {
		case '"':
			if lastChar != '\\' {
				l.tokens <- Token{StringToken, l.input[l.start:l.pos]}
				l.pos++
				l.start = l.pos
				return defaultState
			}
			l.pos++
		case '\r', '\n':
			l.invalidToken()
			return nil
		default:
			l.pos++
		}
		lastChar = rune(l.input[l.pos])
	}
}

func operatorState(l *Lexer) State {
	for k, v := range operator {
		if strings.HasPrefix(l.input[l.pos:], k) {
			l.pos += len(k)
			l.start = l.pos
			l.tokens <- Token{v, k}
			return defaultState
		}
	}
	l.invalidToken()
	return nil
}

func typeState(l *Lexer) State {
	for _, t := range types {
		if strings.HasPrefix(l.input[l.pos:], t) {
			l.pos += len(t)
			l.start = l.pos
			l.tokens <- Token{TypeToken, t}
			return defaultState
		}
	}
	l.invalidToken()
	return nil
}
