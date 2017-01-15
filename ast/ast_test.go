package ast

import (
	"testing"

	"monkey/token"
)

func tok(typ token.Type, lit string) token.Token {
	return token.Token{
		Type:    typ,
		Literal: lit,
	}
}

func TestString(t *testing.T) {
	prog := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: tok(token.LET, "let"),
				Name: &Identifier{
					Token: tok(token.IDENT, "myVar"),
					Value: "myVar",
				},
				Value: &Identifier{
					Token: tok(token.IDENT, "anotherVar"),
					Value: "anotherVar",
				},
			},
		},
	}
	if got, want := prog.String(), "let myVar = anotherVar;"; got != want {
		t.Errorf("prog.String() = %q, want %q", got, want)
	}
}
