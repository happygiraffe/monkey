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
		t.Errorf("intLit.Value = %d, want %d", got, want)
	}
	if got, want := intLit.TokenLiteral(), "5"; got != want {
		t.Errorf("intLit.TokenLiteral() = %q, want %q", got, want)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	tests := []struct {
		input, wantOp string
		wantInt       int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}
	for i, tc := range tests {
		p := New(lexer.New(tc.input))
		prog := p.Parse()
		checkParseErrors(t, p) // TODO: should show input‥

		if got, want := len(prog.Statements), 1; got != want {
			t.Fatalf("%d. len(prog.Statements) = %d, want %d", i, got, want)
		}

		stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is a %T, want *ast.ExpressionStatement", stmt)
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is a %T, want *ast.PrefixExpression", stmt.Expression)
		}
		if got, want := exp.Operator, tc.wantOp; got != want {
			t.Errorf("exp.Operator = %q, want %q", got, want)
		}
		if err := testIntegerLiteral(exp.Right, tc.wantInt); err != nil {
			t.Errorf("%d. IntegerLiteral(%q): %v", i, tc.input, err)
		}
	}
}

func testIntegerLiteral(exp ast.Expression, want int64) error {
	il, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		return fmt.Errorf("got a %T, want a *ast.IntegerLiteral")
	}
	if il.Value != want {
		return fmt.Errorf("got value %d, want %d", il.Value, want)
	}
	if got, want := il.TokenLiteral(), fmt.Sprintf("%d", want); got != want {
		return fmt.Errorf("got TokenLiteral() %q, want %q", got, want)
	}
	return nil
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		input string
		lval  int64
		op    string
		rval  int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}
	for i, tc := range tests {
		p := New(lexer.New(tc.input))
		prog := p.Parse()
		checkParseErrors(t, p)
		if got, want := len(prog.Statements), 1; got != want {
			t.Fatalf("%d. len(prog.Statements) = %d, want %d", i, got, want)
		}

		stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is a %T, want *ast.ExpressionStatement", stmt)
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is a %T, want *ast.InfixExpression", stmt.Expression)
		}
		if err := testIntegerLiteral(exp.Left, tc.lval); err != nil {
			t.Errorf("%d. IntegerLiteral(%q): %v", i, tc.input, err)
		}
		if got, want := exp.Operator, tc.op; got != want {
			t.Errorf("exp.Operator = %q, want %q", got, want)
		}
		if err := testIntegerLiteral(exp.Right, tc.rval); err != nil {
			t.Errorf("%d. IntegerLiteral(%q): %v", i, tc.input, err)
		}
	}
}
