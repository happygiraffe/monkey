package lexer

import (
	"monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		want    token.Type
		wantLit string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lex := New(input)
	for i, tc := range tests {
		tok := lex.NextToken()
		if tok.Type != tc.want {
			t.Fatalf("%d. token type = %v, want %v", i, tok.Type, tc.want)
		}
		if tok.Literal != tc.wantLit {
			t.Fatalf("%d. token literal = %q, want %q", i, tok.Literal, tc.wantLit)
		}
	}
}
