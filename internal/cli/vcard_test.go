package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

// Pins the "Contact:" notice: the vCard FN when any name is given, else the
// email, else the phone -- whatever the command actually has.
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
		{"phone only", []string{"-p", "+123"}, "Contact: +123"},
	}

	origFirst, origLast, origEmail, origPhone, origOut := vcardFirstName, vcardLastName, vcardEmail, vcardPhone, outputFile
	t.Cleanup(func() {
		vcardFirstName, vcardLastName, vcardEmail, vcardPhone, outputFile = origFirst, origLast, origEmail, origPhone, origOut
	})
	silenceStderr(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vcardFirstName, vcardLastName, vcardEmail, vcardPhone = "", "", "", ""
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
