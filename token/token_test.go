package token

import "testing"

func TestLookup(t *testing.T) {
	tests := []struct {
		ident string
		want  Type
	}{
		{"fn", FUNCTION},
		{"let", LET},
		{"true", TRUE},
		{"false", FALSE},
		{"if", IF},
		{"else", ELSE},
		{"return", RETURN},
		{"xyzzy", IDENT},
	}
	for i, tc := range tests {
		if got := Lookup(tc.ident); got != tc.want {
			t.Errorf("%d. Lookup(%q) = %v, want %v", i, tc.ident, got, tc.want)
		}
	}
}
