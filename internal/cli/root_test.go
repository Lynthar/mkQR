package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

// Pins the closed-stdin path: os.Stdin.Stat() fails when fd 0 is closed (some
// daemons and CI runners do that), and a dropped error there is a nil deref.
func TestRootWithUnstattableStdinShowsHelp(t *testing.T) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatalf("open %s: %v", os.DevNull, err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close %s: %v", os.DevNull, err)
	}

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

// unlatchHelp clears the --help flag cobra leaves set on cmd after a help
// run; without it every later Execute of cmd prints help instead of running.
func unlatchHelp(t *testing.T, cmd *cobra.Command) {
	t.Cleanup(func() { _ = cmd.Flags().Set("help", "false") })
}

// silenceStderr routes os.Stderr to the null device until the test ends;
// generateQR writes its "Saved to:" line there rather than to the command.
func silenceStderr(t *testing.T) {
	t.Helper()
	null, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("open %s: %v", os.DevNull, err)
	}
	orig := os.Stderr
	os.Stderr = null
	t.Cleanup(func() {
		os.Stderr = orig
		_ = null.Close()
	})
}

// Pins --size validation on both generation paths: the single-shot commands
// and batch must reject a non-positive size with the same error.
func TestSizeMustBePositiveOnEveryPath(t *testing.T) {
	input := filepath.Join(t.TempDir(), "in.txt")
	if err := os.WriteFile(input, []byte("hello\n"), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	orig := outputSize
	t.Cleanup(func() { outputSize = orig })

	const want = "size must be a positive number, got 0"
	for _, args := range [][]string{
		{"text", "hi", "--size", "0"},
		{"batch", input, "--size", "0"},
	} {
		t.Run(args[0], func(t *testing.T) {
			var buf bytes.Buffer
			cmd := GetRootCmd()
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs(args)
			err := cmd.Execute()
			if err == nil || err.Error() != want {
				t.Errorf("%v returned %v, want %q", args, err, want)
			}
		})
	}
}
