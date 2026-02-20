package lexer

import (
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "context tokens",
			input: "player room item player.room",
			expected: []Token{
				{ContextToken, "player"},
				{ContextToken, "room"},
				{ContextToken, "item"},
				{ContextToken, "player"},
				{MethodCallToken, "."},
				{ContextToken, "room"},
				{EofToken, ""},
			},
		},
		{
			name:  "identifiers",
			input: "myVar foo_bar",
			expected: []Token{
				{IdentifierToken, "myVar"},
				{IdentifierToken, "foo_bar"},
				{EofToken, ""},
			},
		},
		{
			name:  "keywords and operators",
			input: "func if == := \"str\"",
			expected: []Token{
				{FuncToken, "func"},
				{IfToken, "if"},
				{EqualToken, "=="},
				{CreateAndAssignToken, ":="},
				{StringToken, "str"},
				{EofToken, ""},
			},
		},
		{
			name:  "mixed parsing",
			input: "if player.health == 100 { return }",
			expected: []Token{
				{IfToken, "if"},
				{ContextToken, "player"},
				{MethodCallToken, "."},
				{IdentifierToken, "health"},
				{EqualToken, "=="},
				{NumericToken, "100"},
				{OpenBraceToken, "{"},
				{ReturnToken, "return"},
				{CloseBraceToken, "}"},
				{EofToken, ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)

			var tokens []Token
			for {
				tok := lexer.ReadNext()
				tokens = append(tokens, *tok)
				if tok.Typ == EofToken || tok.Typ == InvalidToken {
					break
				}
			}

			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d. Tokens: %v", len(tt.expected), len(tokens), tokens)
			}

			for i, expectedTok := range tt.expected {
				if tokens[i].Typ != expectedTok.Typ || tokens[i].GetRawValue() != expectedTok.GetRawValue() {
					t.Errorf("token %d: expected %+v, got %+v", i, expectedTok, tokens[i])
				}
			}
		})
	}
}
