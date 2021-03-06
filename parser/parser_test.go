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
	tests := []struct {
		input   string
		wantID  string
		wantVal interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = 10;", "y", 10},
		{"let foobar = 836383;", "foobar", 836383},
	}
	for _, tc := range tests {
		p := New(lexer.New(tc.input))
		prog := p.Parse()
		checkParseErrors(t, p)
		if got, want := len(prog.Statements), 1; got != want {
			t.Fatalf("len(prog.Statements) = %d, want %d", got, want)
		}

		stmt := prog.Statements[0]
		if err := testLetStatement(t, stmt, tc.wantID); err != nil {
			t.Errorf("%q: let statement: %v", tc.input, err)
		}

		val := stmt.(*ast.LetStatement).Value
		if err := testLiteralExpression(val, tc.wantVal); err != nil {
			t.Errorf("%q: value: %v", tc.input, err)
		}
	}

	input := `
		let x = 5;
		let y = 10;
		let foobar = 836383;
	`
	p := New(lexer.New(input))
	prog := p.Parse()
	checkParseErrors(t, p)

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
	tests := []struct {
		input   string
		wantVal interface{}
	}{
		{"return 5;", 5},
		{"return 10;", 10},
		{"return 993322;", 993322},
	}
	for i, tc := range tests {
		p := New(lexer.New(tc.input))
		prog := p.Parse()
		checkParseErrors(t, p)
		if got, want := len(prog.Statements), 1; got != want {
			t.Fatalf("len(prog.Statements) = %d, want %d", got, want)
		}

		stmt := prog.Statements[0]
		retStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("%d. %q: stmt is a %T, want *ast.ReturnStatement", i, tc.input, stmt)
			continue
		}
		if got, want := retStmt.TokenLiteral(), "return"; got != want {
			t.Errorf("%d. %q: stmt.TokenLiteral() = %q, want %q", i, tc.input, got, want)
		}
		if err := testLiteralExpression(retStmt.ReturnValue, tc.wantVal); err != nil {
			t.Errorf("%d. %q: stmt.ReturnValue: %v", i, tc.input, err)
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
		return fmt.Errorf("got a %T, want a *ast.IntegerLiteral", exp)
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
		return fmt.Errorf("got a %T, want a *ast.Identifier", exp)
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
		return fmt.Errorf("got a %T, want a *ast.Boolean", exp)
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
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
	}
	for i, tc := range tests {
		p := New(lexer.New(tc.input))
		prog := p.Parse()
		checkParseErrors(t, p)

		if got, want := prog.String(), tc.want; got != want {
			t.Errorf("%d. Parse(%q) = %q, want: %q", i, tc.input, got, want)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	p := New(lexer.New(input))
	prog := p.Parse()
	checkParseErrors(t, p)

	if got, want := len(prog.Statements), 1; got != want {
		t.Fatalf("Parse() got %d statements, want %d", got, want)
	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is a %T, want *ast.ExpressionStatement", stmt)
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is a %t, want a *ast.IfExpression", stmt.Expression)
	}

	// Condition

	if err := testInfixExpression(exp.Condition, "x", "<", "y"); err != nil {
		t.Fatal(err)
	}

	// Consequence

	if got, want := len(exp.Consequence.Statements), 1; got != want {
		t.Fatalf("exp.Consequence got %d statements, want %d", got, want)
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence.Statements[0] is a %t, want a *ast.IfExpression", exp.Consequence.Statements[0])
	}

	if err := testIdentifier(consequence.Expression, "x"); err != nil {
		t.Fatal(err)
	}

	// Alternative

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative = %v, want nil", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	p := New(lexer.New(input))
	prog := p.Parse()
	checkParseErrors(t, p)

	if got, want := len(prog.Statements), 1; got != want {
		t.Fatalf("Parse() got %d statements, want %d", got, want)
	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is a %T, want *ast.ExpressionStatement", stmt)
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is a %t, want a *ast.IfExpression", stmt.Expression)
	}

	// Condition

	if err := testInfixExpression(exp.Condition, "x", "<", "y"); err != nil {
		t.Fatal(err)
	}

	// Consequence

	if got, want := len(exp.Consequence.Statements), 1; got != want {
		t.Fatalf("exp.Consequence got %d statements, want %d", got, want)
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence.Statements[0] is a %t, want a *ast.IfExpression", exp.Consequence.Statements[0])
	}

	if err := testIdentifier(consequence.Expression, "x"); err != nil {
		t.Fatal(err)
	}

	// Alternative

	if got, want := len(exp.Alternative.Statements), 1; got != want {
		t.Fatalf("exp.Alternative got %d statements, want %d", got, want)
	}

	alt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative.Statements[0] is a %t, want a *ast.IfExpression", exp.Alternative.Statements[0])
	}

	if err := testIdentifier(alt.Expression, "y"); err != nil {
		t.Fatal(err)
	}
}

func testInfixExpression(exp ast.Expression, left interface{}, op string, right interface{}) error {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		return fmt.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
	}

	if err := testLiteralExpression(opExp.Left, left); err != nil {
		return err
	}

	if opExp.Operator != op {
		return fmt.Errorf("exp.Operator is not '%s'. got=%q", op, opExp.Operator)
	}

	if err := testLiteralExpression(opExp.Right, right); err != nil {
		return err
	}

	return nil
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	p := New(lexer.New(input))
	prog := p.Parse()
	checkParseErrors(t, p)

	if got, want := len(prog.Statements), 1; got != want {
		t.Fatalf("got %d statements, want %d", got, want)
	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("prog.Statements[0] is a %t, want a *ast.ExpressionStatement", prog.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatal("stmt.Expression is a %t, want a *ast.FunctionLiteral", stmt.Expression)
	}

	// Parameters

	if got, want := len(fn.Parameters), 2; got != want {
		t.Fatalf("got %d parameters, want %d", got, want)
	}

	if err := testLiteralExpression(fn.Parameters[0], "x"); err != nil {
		t.Fatal(err)
	}
	if err := testLiteralExpression(fn.Parameters[1], "y"); err != nil {
		t.Fatal(err)
	}

	// Body

	if got, want := len(fn.Body.Statements), 1; got != want {
		t.Fatalf("got %d body statements, want %d", got, want)
	}

	bodyStmt, ok := fn.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("fn.Body.Statements[0 is a %T, want a *ast.ExpressionStatement", fn.Body.Statements[0])
	}
	if err := testInfixExpression(bodyStmt.Expression, "x", "+", "y"); err != nil {
		t.Fatal(err)
	}
}

func TestParseFunctionParameters(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"fn() {};", nil},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}
	for _, tc := range tests {
		p := New(lexer.New(tc.input))
		prog := p.Parse()
		checkParseErrors(t, p)
		// Just let this blow up at runtime if wrong‥
		stmt := prog.Statements[0].(*ast.ExpressionStatement)
		fn := stmt.Expression.(*ast.FunctionLiteral)
		if got, want := len(fn.Parameters), len(tc.want); got != want {
			t.Fatalf("%q has %d params, want %d", tc.input, got, want)
		}
		for i, ident := range tc.want {
			if err := testLiteralExpression(fn.Parameters[i], ident); err != nil {
				t.Errorf("%q param[%d]: %v", tc.input, i, err)
			}
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := `add(1, 2*3, 4+5);`
	p := New(lexer.New(input))
	prog := p.Parse()
	checkParseErrors(t, p)
	if got, want := len(prog.Statements), 1; got != want {
		t.Fatalf("len(prog.Statements) = %d, want %d", got, want)
	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("prog.Statements[0] is a %t, want a *ast.ExpressionStatement", prog.Statements[0])
	}

	ce, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatal("stmt.Expression is a %t, want a *ast.CallExpression", stmt.Expression)
	}

	if err := testIdentifier(ce.Function, "add"); err != nil {
		t.Fatalf("function name: %v", err)
	}

	if got, want := len(ce.Arguments), 3; got != want {
		t.Fatalf("len(ce.Arguments) = %d, want %d", got, want)
	}
	if err := testLiteralExpression(ce.Arguments[0], 1); err != nil {
		t.Errorf("arg[0]: %v", err)
	}
	if err := testInfixExpression(ce.Arguments[1], 2, "*", 3); err != nil {
		t.Errorf("arg[1]: %v", err)
	}
	if err := testInfixExpression(ce.Arguments[2], 4, "+", 5); err != nil {
		t.Errorf("arg[2]: %v", err)
	}
}
