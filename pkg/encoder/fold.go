package encoder

import (
	"strings"
	"unicode/utf8"
)

// writeFolded writes a single iCalendar / vCard content line to b with the
// line-folding transformation from RFC 5545 §3.1 / RFC 2426 §2.6 applied,
// followed by a trailing CRLF.
//
// Lines ≤75 octets are written unchanged. Longer lines are broken into
// physical lines of at most 75 octets; each continuation line begins with
// a single SPACE (which counts toward the 75-octet cap, so continuation
// lines carry at most 74 octets of payload). Breaks happen only at UTF-8
// rune boundaries — no multi-byte sequence is ever split.
//
// Callers pass the logical content line WITHOUT a trailing CRLF. Mixing
// raw CR/LF into the input would corrupt the output; upstream escaping
// (escapeICal / escapeVCard) already handles user-supplied newlines.
func writeFolded(b *strings.Builder, line string) {
	const limit = 75

	if len(line) <= limit {
		b.WriteString(line)
		b.WriteString("\r\n")
		return
	}

	first := true
	for len(line) > 0 {
		budget := limit
		if !first {
			// Continuation lines start with SPACE, consuming one octet of
			// this physical line's 75-byte budget.
			b.WriteByte(' ')
			budget = limit - 1
		}

		if len(line) <= budget {
			b.WriteString(line)
			b.WriteString("\r\n")
			return
		}

		// Longest UTF-8-aligned prefix whose byte length is ≤ budget.
		// With budget ≥ 74 and max rune size 4 bytes, cut is guaranteed
		// to advance at least one rune per iteration — no infinite loop.
		cut := 0
		for i, r := range line {
			runeSize := utf8.RuneLen(r)
			if i+runeSize > budget {
				break
			}
			cut = i + runeSize
		}
		b.WriteString(line[:cut])
		b.WriteString("\r\n")
		line = line[cut:]
		first = false
	}
}
