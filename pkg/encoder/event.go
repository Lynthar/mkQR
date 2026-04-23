package encoder

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// Event encodes a calendar event in RFC 5545 (iCalendar) format.
//
// Zero-valued End means no end time is emitted. AllDay switches DTSTART/DTEND
// to VALUE=DATE (date-only) form; time-of-day is ignored in that case.
//
// UID and DTStamp are optional: when omitted, Encode generates a random UID
// and stamps the current time. Callers that need stable output (tests,
// re-publishing the same invitation) can set them explicitly.
type Event struct {
	Summary     string
	Start       time.Time
	End         time.Time
	Location    string
	Description string
	URL         string
	AllDay      bool
	UID         string
	DTStamp     time.Time
}

// Encode returns the iCalendar-formatted string wrapped in VCALENDAR.
// The VCALENDAR wrapper improves scanner compatibility — iOS accepts bare
// VEVENT but several Android scanners only recognize the full calendar object.
//
// Content lines are delimited by CRLF per RFC 5545 §3.1. UID (§3.8.4.7) and
// DTSTAMP (§3.8.7.2) are emitted because they are REQUIRED for a VEVENT;
// scanner apps are lenient, but Outlook and many CalDAV servers reject events
// that lack them.
func (e *Event) Encode() string {
	var b strings.Builder

	b.WriteString("BEGIN:VCALENDAR\r\n")
	b.WriteString("VERSION:2.0\r\n")
	b.WriteString("PRODID:-//mkQR//mkQR//EN\r\n")
	b.WriteString("BEGIN:VEVENT\r\n")

	uid := e.UID
	if uid == "" {
		uid = generateUID()
	}
	fmt.Fprintf(&b, "UID:%s\r\n", uid)

	dtstamp := e.DTStamp
	if dtstamp.IsZero() {
		dtstamp = time.Now()
	}
	fmt.Fprintf(&b, "DTSTAMP:%s\r\n", dtstamp.UTC().Format("20060102T150405Z"))

	if e.Summary != "" {
		fmt.Fprintf(&b, "SUMMARY:%s\r\n", escapeICal(e.Summary))
	}

	if !e.Start.IsZero() {
		if e.AllDay {
			// VALUE=DATE is a floating date per RFC 5545; formatting with the
			// stored location preserves the user's intended wall-clock date.
			// Converting through UTC first would shift the date for any
			// non-UTC timezone whose midnight crosses the UTC day boundary.
			fmt.Fprintf(&b, "DTSTART;VALUE=DATE:%s\r\n", e.Start.Format("20060102"))
		} else {
			fmt.Fprintf(&b, "DTSTART:%s\r\n", e.Start.UTC().Format("20060102T150405Z"))
		}
	}

	if !e.End.IsZero() {
		if e.AllDay {
			fmt.Fprintf(&b, "DTEND;VALUE=DATE:%s\r\n", e.End.Format("20060102"))
		} else {
			fmt.Fprintf(&b, "DTEND:%s\r\n", e.End.UTC().Format("20060102T150405Z"))
		}
	}

	if e.Location != "" {
		fmt.Fprintf(&b, "LOCATION:%s\r\n", escapeICal(e.Location))
	}
	if e.Description != "" {
		fmt.Fprintf(&b, "DESCRIPTION:%s\r\n", escapeICal(e.Description))
	}
	if e.URL != "" {
		// URLs are specified by RFC 5545 as being formatted per RFC 3986 and
		// not subject to text escaping.
		fmt.Fprintf(&b, "URL:%s\r\n", e.URL)
	}

	b.WriteString("END:VEVENT\r\n")
	b.WriteString("END:VCALENDAR")

	return b.String()
}

// generateUID produces a globally-unique identifier for a VEVENT.
// Format: <32-hex-chars>@mkqr. The literal domain avoids leaking the host's
// hostname — a VEVENT travels inside a QR a user may publish widely.
func generateUID() string {
	var buf [16]byte
	_, _ = rand.Read(buf[:])
	return hex.EncodeToString(buf[:]) + "@mkqr"
}

// escapeICal escapes the special characters listed in RFC 5545 §3.3.11.
func escapeICal(s string) string {
	r := strings.NewReplacer(
		`\`, `\\`,
		`;`, `\;`,
		`,`, `\,`,
		"\n", `\n`,
	)
	return r.Replace(s)
}
