package parser

import (
	"fmt"
	"reflect"
	"testing"

	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
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

	// TODO: finish parsing let expressions so we get values stored and can compare structs directly.
	want := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{token.LET, "let"},
				Name:  &ast.Identifier{Token: token.Token{token.IDENT, "x"}, Value: "x"},
				Value: &ast.IntegerLiteral{Token: token.Token{token.INT, "5"}, Value: 5},
			},
			&ast.LetStatement{
				Token: token.Token{token.LET, "let"},
				Name:  &ast.Identifier{Token: token.Token{token.IDENT, "y"}, Value: "y"},
				Value: &ast.IntegerLiteral{Token: token.Token{token.INT, "10"}, Value: 10},
			},
			&ast.LetStatement{
				Token: token.Token{token.LET, "let"},
				Name:  &ast.Identifier{Token: token.Token{token.IDENT, "foobar"}, Value: "foobar"},
				Value: &ast.IntegerLiteral{Token: token.Token{token.INT, "836383"}, Value: 836383},
			},
		},
	}
	if !reflect.DeepEqual(prog, want) {
		t.Logf("Parse(%q) =\n%v\nwant:\n%v", input, prog, want)
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
	if err := testIntegerLiteral(stmt.Expression, 5); err != nil {
		t.Error(err)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	tests := []struct {
		input, wantOp string
		wantVal       interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
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
		if err := testLiteralExpression(exp.Right, tc.wantVal); err != nil {
			t.Errorf("%d. %q: Right: %v", i, tc.input, err)
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

func testIdentifier(exp ast.Expression, want string) error {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		return fmt.Errorf("got a %T, want a *ast.Identifier")
	}
	if ident.Value != want {
		return fmt.Errorf("got value %q, want %q", ident.Value, want)
	}
	if got := ident.TokenLiteral(); got != want {
		return fmt.Errorf("got TokenLiteral() %q, want %q", got, want)
	}
	return nil
}

func testBooleanLiteral(exp ast.Expression, want bool) error {
	bl, ok := exp.(*ast.Boolean)
	if !ok {
		return fmt.Errorf("got a %T, want a *ast.Boolean")
	}
	if bl.Value != want {
		return fmt.Errorf("got value %d, want %d", bl.Value, want)
	}
	if got, want := bl.TokenLiteral(), fmt.Sprintf("%t", want); got != want {
		return fmt.Errorf("got TokenLiteral() %q, want %q", got, want)
	}
	return nil
}

func testLiteralExpression(exp ast.Expression, want interface{}) error {
	switch v := want.(type) {
	case int:
		return testIntegerLiteral(exp, int64(v))
	case int64:
		return testIntegerLiteral(exp, v)
	case string:
		return testIdentifier(exp, v)
	case bool:
		return testBooleanLiteral(exp, v)
	default:
		return fmt.Errorf("type %T not handled (for %T)", want, exp)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		input string
		lval  interface{}
		op    string
		rval  interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
		{"false == false;", false, "==", false},
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
		if err := testLiteralExpression(exp.Left, tc.lval); err != nil {
			t.Errorf("%d. %q: Left: %v", i, tc.input, err)
		}
		if got, want := exp.Operator, tc.op; got != want {
			t.Errorf("exp.Operator = %q, want %q", got, want)
		}
		if err := testLiteralExpression(exp.Right, tc.rval); err != nil {
			t.Errorf("%d. %q: Right: %v", i, tc.input, err)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
	}
	for i, tc := range tests {
		p := New(lexer.New(tc.input))
		prog := p.Parse()
		checkParseErrors(t, p)

		if got, want := prog.String(), tc.want; got != want {
			t.Errorf("%d. Parse(%q) = %q, want %q", i, tc.input, got, want)
		}
	}
}
