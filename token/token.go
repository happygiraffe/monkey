// Package token represents all the possible tokens that the lexer can use.
package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers & literals.
	IDENT TokenType = "IDENT" // add, foobar, x, y, ‥
	INT   TokenType = "INT"   // 123456

	// Operators
	ASSIGN TokenType = "="
	PLUS   TokenType = "+"

	// DELIMITERS
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"

	LPAREN TokenType = "("
	RPAREN TokenType = ")"

	LBRACE TokenType = "{"
	RBRACE TokenType = "}"

	// Keywords
	FUNCTION TokenType = "FUNCTION"
	LET      TokenType = "LET"
)
