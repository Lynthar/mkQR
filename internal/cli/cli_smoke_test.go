package cli

import (
	"bytes"
	"testing"
)

func TestEnsureHTTPScheme(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"https://example.com", "https://example.com"},
		{"http://example.com", "http://example.com"},
		{"HTTPS://EXAMPLE.COM", "HTTPS://EXAMPLE.COM"},
		{"example.com", "https://example.com"},
		{"example.com/path", "https://example.com/path"},
		// Regression: the previous `HasPrefix(..., "http")` check falsely
		// accepted domains that happened to start with "http", leaving them
		// without a scheme and producing invalid QR content.
		{"httpfoobar.com", "https://httpfoobar.com"},
		{"httpsfoo.example", "https://httpsfoo.example"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := ensureHTTPScheme(tt.in); got != tt.want {
				t.Errorf("ensureHTTPScheme(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestSubcommandHelpDoesNotPanic verifies every subcommand can render its
// --help without panicking.
//
// Cobra detects short-flag collisions between a subcommand flag and a root
// persistent flag only at lookup time (not compile time); the panic happens
// the first time the subcommand is resolved — even for `--help`. This class
// of bug has shipped to main twice (`mkqr geo -q` vs root `--quiet`, and
// `mkqr vcard -l` vs root `--level`). Running `--help` against every
// registered subcommand is the cheapest sentinel for the next regression.
//
// The currently-reserved short letters (case-sensitive) are o, l, q, v.
func TestSubcommandHelpDoesNotPanic(t *testing.T) {
	subcommands := []string{
		"wifi", "vcard", "otp", "email", "phone", "sms", "geo",
		"url", "text", "event", "batch",
	}
	for _, name := range subcommands {
		name := name
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := GetRootCmd()
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs([]string{name, "--help"})
			if err := cmd.Execute(); err != nil {
				t.Fatalf("%s --help returned error: %v", name, err)
			}
			if buf.Len() == 0 {
				t.Errorf("%s --help produced no output", name)
			}
		})
	}
}
