package parser

import (
	"fmt"
	"testing"

	"monkey/ast"
	"monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 836383;
	`
	p := New(lexer.New(input))
	prog := p.Parse()
	checkParseErrors(t, p)
	if prog == nil {
		t.Fatal("Parse() returned nil")
	}
	if got, want := len(prog.Statements), 3; got != want {
		t.Fatalf("Parse() got %d statements, want %d", got, want)
	}

	tests := []struct {
		wantIdent string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tc := range tests {
		if err := testLetStatement(t, prog.Statements[i], tc.wantIdent); err != nil {
			t.Fatalf("%d. testLetStatement: %v", i, err)
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) error {
	if got, want := s.TokenLiteral(), "let"; got != want {
		return fmt.Errorf("s.TokenLiteral = %q, want %q", got, want)
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		return fmt.Errorf("s is a %T, want LetStatement", s)
	}
	if got, want := letStmt.Name.Value, name; got != want {
		return fmt.Errorf("letStmt.Name.Value = %q, want %q", got, want)
	}
	if got, want := letStmt.Name.TokenLiteral(), name; got != want {
		return fmt.Errorf("letStmt.Name.TokenLiteral() = %q, want %q", got, want)
	}
	return nil
}

func checkParseErrors(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	for i, msg := range errs {
		t.Errorf("parse error %d: %v", i, msg)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	input := `
    return 5;
    return 10;
	  return 993322;
	`
	l := lexer.New(input)
	p := New(l)
	prog := p.Parse()
	checkParseErrors(t, p)
	if got, want := len(prog.Statements), 3; got != want {
		t.Fatalf("len(prog.Statements) = %d, want %d", got, want)
	}
	for _, stmt := range prog.Statements {
		retStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt is a %T, want *ast.ReturnStatement", stmt)
			continue
		}
		if got, want := retStmt.TokenLiteral(), "return"; got != want {
			t.Errorf("stmt.TokenLiteral() = %q, want %q", got, want)
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	prog := p.Parse()
	checkParseErrors(t, p)

	if got, want := len(prog.Statements), 1; got != want {
		t.Fatalf("len(prog.Statements) = %d, want %d", got, want)
	}
	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is a %T, want *ast.ExpressionStatement", stmt)
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Expression is a %T, want *ast.Identifier", stmt.Expression)
	}
	if got, want := ident.Value, "foobar"; got != want {
		t.Errorf("ident.Value = %q, want %q", got, want)
	}
	if got, want := ident.TokenLiteral(), "foobar"; got != want {
		t.Errorf("ident.TokenLiteral() = %q, want %q", got, want)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	prog := p.Parse()
	checkParseErrors(t, p)

	if got, want := len(prog.Statements), 1; got != want {
		t.Fatalf("len(prog.Statements) = %d, want %d", got, want)
	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is a %T, want *ast.ExpressionStatement", stmt)
	}
	intLit, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is a %T, want *ast.IntegerLiteral", stmt.Expression)
	}
	if got, want := intLit.Value, int64(5); got != want {
		t.Errorf("intLit.Value = %q, want %q", got, want)
	}
	if got, want := intLit.TokenLiteral(), "5"; got != want {
		t.Errorf("intLit.TokenLiteral() = %q, want %q", got, want)
	}
}
