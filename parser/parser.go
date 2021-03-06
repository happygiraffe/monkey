// Package parser implements a parser for the monkey programming language.
package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type prec int

const (
	LOWEST      prec = iota + 1
	EQUALS           // ==
	LESSGREATER      // > or <
	SUM              // +
	PRODUCT          // *
	PREFIX           // -X or !X
	CALL             // myFunc(X)
)

type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

var precedences = map[token.Type]prec{
	token.EQ:       EQUALS,
	token.NE:       EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

// Parser allows parsing the monkey language.
type Parser struct {
	l              *lexer.Lexer
	curTok         token.Token
	peekTok        token.Token
	errors         []string
	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		prefixParseFns: make(map[token.Type]prefixParseFn),
		infixParseFns:  make(map[token.Type]infixParseFn),
	}

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NE, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	// Read two tokens so curTok and peekTok are ready to use.
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPrefix(tt token.Type, fn prefixParseFn) {
	p.prefixParseFns[tt] = fn
}

func (p *Parser) registerInfix(tt token.Type, fn infixParseFn) {
	p.infixParseFns[tt] = fn
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
		return p.parseExpressionStatement()
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

	p.nextToken()

	st.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
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
	msg := fmt.Sprintf("expected token %v, got token %v (%q)", t, p.peekTok.Type, p.peekTok.Literal)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	st := &ast.ReturnStatement{Token: p.curTok}
	p.nextToken()

	st.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return st
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curTok}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %v found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence prec) ast.Expression {
	prefix, ok := p.prefixParseFns[p.curTok.Type]
	if !ok {
		p.noPrefixParseFnError(p.curTok.Type)
		return nil
	}
	leftExp := prefix()

	// The heart of the Pratt Parser…
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekTok.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curTok}
	val, err := strconv.ParseInt(p.curTok.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q: %v", p.curTok.Literal, err)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = val
	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.curTok,
		Operator: p.curTok.Literal,
	}
	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) peekPrecedence() prec {
	if p, ok := precedences[p.peekTok.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() prec {
	if p, ok := precedences[p.curTok.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.curTok,
		Operator: p.curTok.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)
	return expr
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curTok, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expr := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.curTok}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curTok}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curTok}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	// Special case: empty function parameter list
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return nil
	}

	p.nextToken()

	ids := []*ast.Identifier{{Token: p.curTok, Value: p.curTok.Literal}}
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ids = append(ids, &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal})
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return ids
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curTok, Function: fn}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return nil
	}

	var args []ast.Expression
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}
