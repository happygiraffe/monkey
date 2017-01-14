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
