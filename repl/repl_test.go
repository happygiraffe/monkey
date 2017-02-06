package repl

import (
	"bytes"
	"strings"
	"testing"
)

func TestRepl(t *testing.T) {
	tests := []struct {
		input     string
		wantLines []string
	}{
		{
			input:     "1;",
			wantLines: []string{"1;"},
		},
		{
			input: "let add = fn(x, y) { x + y; };",
			wantLines: []string{
				"let add = fn(x, y) {",
				"(x + y);",
				"};",
			},
		},
		{
			input: "let y 5 9;",
			wantLines: []string{`	expected token =, got token INT ("5")`},
		},
	}
	for i, tc := range tests {
		in := strings.NewReader(tc.input)
		var out bytes.Buffer
		Start(in, &out)

		want := Prompt + strings.Join(tc.wantLines, "\n") + "\n" + Prompt
		if got := out.String(); got != want {
			t.Errorf("%d. Start(%q) =\n%q\n, want\n%q", i, tc.input, got, want)
		}
	}
}
