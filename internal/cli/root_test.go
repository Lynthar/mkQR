package cli

import (
	"bytes"
	"os"
	"testing"
)

// Pins the closed-stdin path: os.Stdin.Stat() fails when fd 0 is closed (some
// daemons and CI runners do that), and a dropped error there is a nil deref.
func TestRootWithUnstattableStdinShowsHelp(t *testing.T) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatalf("open %s: %v", os.DevNull, err)
	}
	f.Close()

	orig := os.Stdin
	os.Stdin = f
	t.Cleanup(func() { os.Stdin = orig })

	var buf bytes.Buffer
	cmd := GetRootCmd()
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("root with unstattable stdin: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("root with unstattable stdin produced no help output")
	}
}
