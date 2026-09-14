package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
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

// captureStream routes *stream (os.Stdout or os.Stderr) to a file until the
// test ends and returns a function that reads back what was written.
func captureStream(t *testing.T, stream **os.File) func() []byte {
	t.Helper()
	path := filepath.Join(t.TempDir(), "captured")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create %s: %v", path, err)
	}
	orig := *stream
	*stream = f
	t.Cleanup(func() { *stream = orig })
	return func() []byte {
		if err := f.Close(); err != nil {
			t.Fatalf("close %s: %v", path, err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		return data
	}
}

// Pins `-o -`: the PNG goes to stdout, byte for byte what `-o file.png` writes,
// with and without --logo. It used to create a file literally named "-".
func TestOutputDashWritesPNGToStdout(t *testing.T) {
	gen, err := buildGenerator(io.Discard)
	if err != nil {
		t.Fatalf("buildGenerator: %v", err)
	}
	// Any PNG serves as a logo; a QR of its own is the cheapest one at hand.
	logoPNG, err := gen.GeneratePNG("logo")
	if err != nil {
		t.Fatalf("GeneratePNG: %v", err)
	}
	logo := filepath.Join(t.TempDir(), "logo.png")
	if err := os.WriteFile(logo, logoPNG, 0o644); err != nil {
		t.Fatalf("write logo: %v", err)
	}

	origOut, origLogo := outputFile, logoPath
	t.Cleanup(func() { outputFile, logoPath = origOut, origLogo })
	silenceStderr(t)

	for _, tt := range []struct {
		name string
		args []string
	}{
		{"plain", nil},
		{"logo", []string{"--logo", logo}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			logoPath = ""
			read := captureStream(t, &os.Stdout)

			var buf bytes.Buffer
			cmd := GetRootCmd()
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs(append([]string{"text", "hi", "-o", "-"}, tt.args...))
			if err := cmd.Execute(); err != nil {
				t.Fatalf("text -o - %v: %v", tt.args, err)
			}
			got := read()
			if _, err := os.Stat("-"); err == nil {
				_ = os.Remove("-")
				t.Error("-o - created a file named \"-\"")
			}

			want := filepath.Join(t.TempDir(), "want.png")
			cmd.SetArgs(append([]string{"text", "hi", "-o", want}, tt.args...))
			if err := cmd.Execute(); err != nil {
				t.Fatalf("text -o file %v: %v", tt.args, err)
			}
			wantBytes, err := os.ReadFile(want)
			if err != nil {
				t.Fatalf("read %s: %v", want, err)
			}
			if !bytes.Equal(got, wantBytes) {
				t.Errorf("stdout PNG (%d bytes) differs from the file PNG (%d bytes)", len(got), len(wantBytes))
			}
		})
	}
}

// Pins the small-size note on every PNG path: below the module grid the PNG is
// raised to it (this used to happen silently), and SVG — which honors any size
// — says nothing.
func TestSizeBelowGridIsNotedOnPNGPaths(t *testing.T) {
	input := filepath.Join(t.TempDir(), "in.txt")
	if err := os.WriteFile(input, []byte("hello\n"), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	const note = "Note: --size 10 is below the 29-pixel minimum for this code; using 29"

	origSize, origOut := outputSize, outputFile
	t.Cleanup(func() { outputSize, outputFile = origSize, origOut })

	for _, tt := range []struct {
		name     string
		args     []string
		wantNote bool
	}{
		{"png file", []string{"text", "hi", "-o", "OUT/x.png", "--size", "10"}, true},
		{"stdout", []string{"text", "hi", "-o", "-", "--size", "10"}, true},
		{"batch", []string{"batch", input, "-O", "OUT", "--size", "10"}, true},
		{"svg file", []string{"text", "hi", "-o", "OUT/x.svg", "--size", "10"}, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			out := t.TempDir()
			for i, a := range tt.args {
				tt.args[i] = strings.Replace(a, "OUT", out, 1)
			}
			readErr := captureStream(t, &os.Stderr)
			if tt.args[2] == "-" {
				captureStream(t, &os.Stdout)
			}

			var buf bytes.Buffer
			cmd := GetRootCmd()
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs(tt.args)
			if err := cmd.Execute(); err != nil {
				t.Fatalf("%v: %v", tt.args, err)
			}
			all := buf.String() + string(readErr())
			if got := strings.Contains(all, note); got != tt.wantNote {
				t.Errorf("%v: note present = %v, want %v; output was:\n%s", tt.args, got, tt.wantNote, all)
			}
		})
	}
}
