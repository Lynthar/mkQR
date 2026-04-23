package encoder

import (
	"strings"
	"testing"
)

func TestWriteFolded(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "short line no fold",
			in:   "SUMMARY:hi",
			want: "SUMMARY:hi\r\n",
		},
		{
			// Exactly 75 octets: the RFC allows ≤75, so no fold.
			name: "exactly 75 octets no fold",
			in:   strings.Repeat("a", 75),
			want: strings.Repeat("a", 75) + "\r\n",
		},
		{
			// 76 octets: fold once. Continuation line carries only 1 byte
			// after its leading SPACE (budget 74 ≥ 1).
			name: "76 octets folds once",
			in:   strings.Repeat("a", 76),
			want: strings.Repeat("a", 75) + "\r\n " + "a\r\n",
		},
		{
			// 150 octets: first line 75, then continuation 74 (budget 74
			// with leading SPACE), then continuation 1. Pins that the
			// helper keeps making progress and accounts for the SPACE cost.
			name: "150 octets folds twice",
			in:   strings.Repeat("a", 150),
			want: strings.Repeat("a", 75) + "\r\n " + strings.Repeat("a", 74) + "\r\n " + "a\r\n",
		},
		{
			// Boundary: 75 + 74 = 149 should fold into exactly two lines,
			// not three — the continuation line must fill up to 74 bytes
			// (full 75-octet physical line including the leading SPACE).
			name: "149 octets folds into two lines exactly",
			in:   strings.Repeat("a", 149),
			want: strings.Repeat("a", 75) + "\r\n " + strings.Repeat("a", 74) + "\r\n",
		},
		{
			// CJK: each "中" is 3 UTF-8 bytes. "SUMMARY:" is 8 bytes, so
			// the first line can hold (75-8)/3 = 22 full characters for
			// 8 + 66 = 74 bytes total (can't fit a 23rd because 74+3=77>75).
			name: "CJK breaks at rune boundary",
			in:   "SUMMARY:" + strings.Repeat("中", 30),
			want: "SUMMARY:" + strings.Repeat("中", 22) + "\r\n " + strings.Repeat("中", 8) + "\r\n",
		},
		{
			// A 3-byte rune straddling byte 73: the fold MUST happen before
			// it (cut at 73), not split the multi-byte sequence.
			name: "3-byte rune straddling octet 75 is pushed to next line",
			in:   strings.Repeat("a", 73) + "中b",
			want: strings.Repeat("a", 73) + "\r\n " + "中b\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b strings.Builder
			writeFolded(&b, tt.in)
			got := b.String()
			if got != tt.want {
				t.Errorf("writeFolded(%q):\n got: %q\nwant: %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestWriteFoldedPhysicalLinesWithinLimit(t *testing.T) {
	// Property: no physical line (between CRLFs) exceeds 75 octets.
	inputs := []string{
		strings.Repeat("a", 1000),
		"SUMMARY:" + strings.Repeat("中", 100),
		strings.Repeat("abc", 50) + strings.Repeat("中", 30) + strings.Repeat("x", 100),
	}
	for _, in := range inputs {
		var b strings.Builder
		writeFolded(&b, in)
		// Trim the terminal CRLF, then split by CRLF to get physical lines.
		out := strings.TrimSuffix(b.String(), "\r\n")
		for i, line := range strings.Split(out, "\r\n") {
			if len(line) > 75 {
				t.Errorf("input len=%d: physical line %d has %d octets (>75): %q",
					len(in), i, len(line), line)
			}
		}
	}
}
