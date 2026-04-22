package cli

import (
	"bytes"
	"testing"
)

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
