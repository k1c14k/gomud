package lexer

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// State represents a transition function for the lexical analyzer.
type State func(*Lexer) State

// Lexer breaks down an input string into a stream of tokens for the parser.
type Lexer struct {
	input  string
	start  int
	pos    int
	tokens chan Token
	state  State
	peeked []*Token
}

// NewLexer initializes a new Lexer with the given input string.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		tokens: make(chan Token, 2),
		state:  defaultState,
	}
}

func (l *Lexer) run() {
	state := defaultState
	for state != nil {
		state = state(l)
	}
	close(l.tokens)
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

// ReadNext consumes and returns the next token from the lexer pipeline.
func (l *Lexer) ReadNext() *Token {
	if len(l.peeked) > 0 {
		t := l.peeked[0]
		l.peeked = l.peeked[1:]
		logrus.Debug("Peeked token", t)
		return t
	}

	token := l.nextToken()
	logrus.Debug("Read token", token)
	return token
}

// PeekSome retrieves the next 'n' tokens without consuming them from the lexer.
func (l *Lexer) PeekSome(n int) []*Token {
	for i := len(l.peeked); i < n; i++ {
		l.peeked = append(l.peeked, l.nextToken())
	}
	return l.peeked[:n]
}

// Peek returns the next token immediately available without consuming it.
func (l *Lexer) Peek() *Token {
	return l.PeekSome(1)[0]
}

var (
	keywords = map[string]TokenType{
		"import": ImportToken,
		"func":   FuncToken,
		"if":     IfToken,
		"else":   ElseToken,
		"var":    VarToken,
		"return": ReturnToken,
	}

	contextNames = map[string]TokenType{
		"player": ContextToken,
		"room":   ContextToken,
		"item":   ContextToken,
	}

	parenthesis = map[string]TokenType{
		"(": OpenParenToken,
		")": CloseParenToken,
		"{": OpenBraceToken,
		"}": CloseBraceToken,
	}

	operator = map[string]TokenType{
		"+":  AddToken,
		"-":  SubtractToken,
		"*":  MultiplyToken,
		"/":  DivideToken,
		"%":  ModuloToken,
		".":  MethodCallToken,
		"==": EqualToken,
		"=":  AssignToken,
		":=": CreateAndAssignToken,
	}

	types = []string{"int", "string"}
)

const validIdentifier = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789"

func isParenthesis(r rune) bool {
	switch r {
	case '(', ')', '{', '}':
		return true
	}
	return false
}

func (l *Lexer) isParenthesis() bool {
	if l.pos >= len(l.input) {
		return false
	}
	return isParenthesis(rune(l.input[l.pos]))
}

func isOperator(r rune) bool {
	switch r {
	case '+', '-', '*', '/', '%', '.', '=', ':':
		return true
	}
	return false
}

func (l *Lexer) isOperator() bool {
	if l.pos >= len(l.input) {
		return false
	}
	return isOperator(rune(l.input[l.pos]))
}

func (l *Lexer) hasPrefix(m map[string]TokenType) bool {
	for k := range m {
		if strings.HasPrefix(l.input[l.pos:], k+" ") {
			return true
		}
	}
	return false
}

func (l *Lexer) isType() bool {
	for _, t := range types {
		if strings.HasPrefix(l.input[l.pos:], t) {
			if l.pos+len(t) < len(l.input) {
				nextChar := rune(l.input[l.pos+len(t)])
				if strings.ContainsRune(validIdentifier, nextChar) {
					continue
				}
			}
			return true
		}
	}
	return false
}

func (l *Lexer) isNumeric() bool {
	if l.pos >= len(l.input) {
		return false
	}
	return l.input[l.pos] >= '0' && l.input[l.pos] <= '9'
}

func (l *Lexer) invalidToken() {
	var context string
	if l.pos+20 > len(l.input) {
		context = l.input[l.pos:]
	} else {
		context = l.input[l.pos : l.pos+20]
	}
	l.tokens <- Token{InvalidToken, "Invalid token near " + context}
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
	case l.isNumeric():
		return numberState
	default:
		return identifierState
	}
}

func (l *Lexer) emit(typ TokenType, val string, advance int) {
	l.tokens <- Token{typ, val}
	l.pos += advance
	l.start = l.pos
}

func numberState(l *Lexer) State {
	for {
		if l.pos >= len(l.input) || !l.isNumeric() {
			val := l.input[l.start:l.pos]
			l.tokens <- Token{NumericToken, val}
			l.start = l.pos
			return defaultState
		}
		l.pos++
	}
}

func keywordState(l *Lexer) State {
	for k, v := range keywords {
		if strings.HasPrefix(l.input[l.pos:], k+" ") {
			l.emit(v, k, len(k)+1)
			return defaultState
		}
	}
	l.invalidToken()
	return nil
}

func identifierState(l *Lexer) State {
loop:
	for {
		if l.pos >= len(l.input) {
			break loop
		}

		char := rune(l.input[l.pos])
		if isParenthesis(char) || isOperator(char) {
			break loop
		}

		switch char {
		case ' ', '\t', '\n', '\r':
			break loop
		default:
			l.pos++
			continue loop
		}
	}

	val := l.input[l.start:l.pos]
	if tok, ok := contextNames[val]; ok {
		l.tokens <- Token{tok, val}
	} else {
		l.tokens <- Token{IdentifierToken, val}
	}
	l.start = l.pos
	return defaultState
}

func parenthesisState(l *Lexer) State {
	for k, v := range parenthesis {
		if strings.HasPrefix(l.input[l.pos:], k) {
			l.emit(v, k, len(k))
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
				val := l.input[l.start:l.pos]
				l.emit(StringToken, val, 1) // Advance past closing quote
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
	switch l.input[l.pos] {
	case '=':
		if l.pos+1 < len(l.input) && l.input[l.pos+1] == '=' {
			l.emit(EqualToken, "==", 2)
			return defaultState
		}
		l.emit(AssignToken, "=", 1)
		return defaultState
	case ':':
		if l.pos+1 < len(l.input) && l.input[l.pos+1] == '=' {
			l.emit(CreateAndAssignToken, ":=", 2)
			return defaultState
		}
	case '.', '+', '-', '*', '/', '%':
		char := l.input[l.pos : l.pos+1]
		l.emit(operator[char], char, 1)
		return defaultState
	}

	l.invalidToken()
	return nil
}

func typeState(l *Lexer) State {
	for _, t := range types {
		if strings.HasPrefix(l.input[l.pos:], t) {
			l.emit(TypeToken, t, len(t))
			return defaultState
		}
	}
	l.invalidToken()
	return nil
}
