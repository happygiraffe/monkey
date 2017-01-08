// Package token represents all the possible tokens that the lexer can use.
package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	ILLEGAL Type = "ILLEGAL"
	EOF     Type = "EOF"

	// Identifiers & literals.
	IDENT Type = "IDENT" // add, foobar, x, y, ‥
	INT   Type = "INT"   // 123456

	// Operators
	ASSIGN Type = "="
	PLUS   Type = "+"

	// DELIMITERS
	COMMA     Type = ","
	SEMICOLON Type = ";"

	LPAREN Type = "("
	RPAREN Type = ")"

	LBRACE Type = "{"
	RBRACE Type = "}"

	// Keywords
	FUNCTION Type = "FUNCTION"
	LET      Type = "LET"
)
