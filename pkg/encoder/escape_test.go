package encoder

import "testing"

func TestEscapeText(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"simple", "simple"},
		{`with\backslash`, `with\\backslash`},
		{"with,comma", `with\,comma`},
		{"with;semicolon", `with\;semicolon`},
		{"with\nnewline", `with\nnewline`},
		{`all\;,chars\nplus`, `all\\\;\,chars\\nplus`},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := escapeText(tt.in); got != tt.want {
				t.Errorf("escapeText(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
