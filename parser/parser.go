// Package parser implements a parser for the monkey programming language.
package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

// Parser allows parsing the monkey language.
type Parser struct {
	l       *lexer.Lexer
	curTok  token.Token
	peekTok token.Token
	errors  []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Read two tokens so curTok and peekTok are ready to use.
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.l.NextToken()
}

func (p *Parser) Parse() *ast.Program {
	prog := &ast.Program{}
	for p.curTok.Type != token.EOF {
		if st := p.parseStatement(); st != nil {
			prog.Statements = append(prog.Statements, st)
		}
		p.nextToken()
	}
	return prog
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curTok.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	st := &ast.LetStatement{Token: p.curTok}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	st.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: skipping expressions until we hit a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return st
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curTok.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekTok.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected token %v, got token %v", t, p.peekTok.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	st := &ast.ReturnStatement{Token: p.curTok}
	p.nextToken()

	// TODO: skipping expressions until we hit a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return st
}
