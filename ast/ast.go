// Package ast is the Abstract Syntax Tree for monkey.
package ast

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
