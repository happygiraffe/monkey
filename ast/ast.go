// Package ast is the Abstract Syntax Tree for monkey.
package ast

import (
	"monkey/token"
)

// Node is a single node in the AST.
type Node interface {
	TokenLiteral() string
}

// Statement is a node which can be evaluated but does not produce a value.
type Statement interface {
	Node
	statementNode()
}

// Expression is a node which can be evaluated to produce a value.
type Expression interface {
	Node
	expressionNode()
}

// LetStatement is a "let x = y" statement.
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// Identifier is a "let x = y" statement.
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
