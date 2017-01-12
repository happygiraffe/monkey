package repl

import (
	"bytes"
	"strings"
	"testing"
)

func TestRepl(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{
			input: "1;",
			want: []string{
				`INT("1")`,
				`;(";")`,
			},
		},
		{
			input: "let add = fn(x, y) { x + y; };",
			want: []string{
				`LET("let")`,
				`IDENT("add")`,
				`=("=")`,
				`FUNCTION("fn")`,
				`(("(")`,
				`IDENT("x")`,
				`,(",")`,
				`IDENT("y")`,
				`)(")")`,
				`{("{")`,
				`IDENT("x")`,
				`+("+")`,
				`IDENT("y")`,
				`;(";")`,
				`}("}")`,
				`;(";")`,
			},
		},
	}
	for i, tc := range tests {
		in := strings.NewReader(tc.input)
		var out bytes.Buffer
		Start(in, &out)

		want := strings.Join(tc.want, "\n") + "\n"
		want = Prompt + want + Prompt
		if got := out.String(); got != want {
			t.Errorf("%d. Start(%q) =\n%q\n, want\n%q", i, tc.input, got, want)
		}
	}
}
