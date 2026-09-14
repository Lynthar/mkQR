package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Pins two batch contracts: a line that fails to generate must not exit 0, and
// the output filename must carry the input line number.
func TestBatchPartialFailure(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "in.txt")
	// Line 2 is past the QR byte capacity at level M, so it is the only failure.
	content := "example.com\n" + strings.Repeat("x", 3000) + "\nhello\n"
	if err := os.WriteFile(input, []byte(content), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	out := filepath.Join(dir, "out")

	var buf bytes.Buffer
	cmd := GetRootCmd()
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"batch", input, "-O", out})
	if err := cmd.Execute(); err == nil {
		t.Error("batch with a failing line returned nil error")
	}

	for _, name := range []string{"qr_0001.png", "qr_0003.png"} {
		if _, err := os.Stat(filepath.Join(out, name)); err != nil {
			t.Errorf("expected %s: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(out, "qr_0002.png")); err == nil {
		t.Error("qr_0002.png exists; line 2 failed, so its number must stay unused")
	}
}

// Pins the -o notice: batch writes PNG into --output-dir and ignores -o, which
// used to happen without a word to the user.
func TestBatchNotesThatOutputFileIsIgnored(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "in.txt")
	if err := os.WriteFile(input, []byte("example.com\n"), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	out := filepath.Join(dir, "out")

	orig := outputFile
	t.Cleanup(func() { outputFile = orig })

	var buf bytes.Buffer
	cmd := GetRootCmd()
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"batch", input, "-O", out, "-o", "ignored.svg"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("batch: %v", err)
	}

	if !strings.Contains(buf.String(), "-o is ignored") {
		t.Errorf("batch with -o gave no notice, output was:\n%s", buf.String())
	}
	if _, err := os.Stat(filepath.Join(out, "qr_0001.png")); err != nil {
		t.Errorf("expected qr_0001.png: %v", err)
	}
}

// Pins the summary line to the files: a --prefix with path components moves
// the output, and the line used to keep naming --output-dir anyway.
func TestBatchSummaryNamesTheDirectoryWritten(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "in.txt")
	if err := os.WriteFile(input, []byte("example.com\n"), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	out := filepath.Join(dir, "out")

	origPrefix := batchPrefix
	t.Cleanup(func() { batchPrefix = origPrefix })

	var buf bytes.Buffer
	cmd := GetRootCmd()
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"batch", input, "-O", out, "--prefix", "../escaped_"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("batch: %v", err)
	}

	written := filepath.Join(dir, "escaped_0001.png")
	if _, err := os.Stat(written); err != nil {
		t.Fatalf("expected %s: %v", written, err)
	}
	// The line must end at dir: "in <dir>" is also a prefix of "in <dir>/out".
	if want := "Generated 1 QR codes in " + dir + "\n"; !strings.Contains(buf.String(), want) {
		t.Errorf("summary lacks %q, output was:\n%s", want, buf.String())
	}
}
