package cli

import (
	"bytes"
	"strings"
	"testing"
)

// Pins the Key URI Format defaults as the user sees them: an authenticator
// that reads a URL without algorithm/digits/period assumes exactly these.
func TestOTPHelpShowsKeyURIDefaults(t *testing.T) {
	unlatchHelp(t, otpCmd)
	var buf bytes.Buffer
	cmd := GetRootCmd()
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"otp", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("otp --help: %v", err)
	}
	for _, want := range []string{
		`--algorithm string   Hash algorithm (SHA1/SHA256/SHA512) (default "SHA1")`,
		`--digits int         Number of digits (6 or 8) (default 6)`,
		`--period int         Time period in seconds (TOTP) (default 30)`,
	} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("otp --help lacks %q, output was:\n%s", want, buf.String())
		}
	}
}
