package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

// Pins the "Contact:" notice: the vCard FN when any name is given, the email
// otherwise.
func TestVCardNoticeShowsName(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"both names", []string{"-f", "John", "--last", "Doe"}, "Contact: John Doe"},
		{"first only", []string{"-f", "John"}, "Contact: John"},
		{"last only", []string{"--last", "Doe"}, "Contact: Doe"},
		{"email only", []string{"-e", "jane@example.com"}, "Contact: jane@example.com"},
	}

	origFirst, origLast, origEmail, origOut := vcardFirstName, vcardLastName, vcardEmail, outputFile
	t.Cleanup(func() {
		vcardFirstName, vcardLastName, vcardEmail, outputFile = origFirst, origLast, origEmail, origOut
	})
	silenceStderr(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vcardFirstName, vcardLastName, vcardEmail = "", "", ""
			out := filepath.Join(t.TempDir(), "got.png")

			var buf bytes.Buffer
			cmd := GetRootCmd()
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs(append([]string{"vcard", "-o", out}, tt.args...))
			if err := cmd.Execute(); err != nil {
				t.Fatalf("vcard %v: %v", tt.args, err)
			}
			if !strings.Contains(buf.String(), tt.want+"\n") {
				t.Errorf("notice lacks %q, output was:\n%s", tt.want, buf.String())
			}
		})
	}
}
