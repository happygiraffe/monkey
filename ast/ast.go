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

// Program is the top-level program.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}
	return p.Statements[0].TokenLiteral()
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

type ReturnStatement struct {
	Token      token.Token
	RetenValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
