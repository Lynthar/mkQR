package encoder

import (
	"fmt"
	"strings"
	"time"
)

// Event encodes a calendar event in RFC 5545 (iCalendar) format.
//
// Zero-valued End means no end time is emitted. AllDay switches DTSTART/DTEND
// to VALUE=DATE (date-only) form; time-of-day is ignored in that case.
type Event struct {
	Summary     string
	Start       time.Time
	End         time.Time
	Location    string
	Description string
	URL         string
	AllDay      bool
}

// Encode returns the iCalendar-formatted string wrapped in VCALENDAR.
// The VCALENDAR wrapper improves scanner compatibility — iOS accepts bare
// VEVENT but several Android scanners only recognize the full calendar object.
func (e *Event) Encode() string {
	var b strings.Builder

	b.WriteString("BEGIN:VCALENDAR\n")
	b.WriteString("VERSION:2.0\n")
	b.WriteString("PRODID:-//mkQR//mkQR//EN\n")
	b.WriteString("BEGIN:VEVENT\n")

	if e.Summary != "" {
		fmt.Fprintf(&b, "SUMMARY:%s\n", escapeICal(e.Summary))
	}

	if !e.Start.IsZero() {
		if e.AllDay {
			fmt.Fprintf(&b, "DTSTART;VALUE=DATE:%s\n", e.Start.UTC().Format("20060102"))
		} else {
			fmt.Fprintf(&b, "DTSTART:%s\n", e.Start.UTC().Format("20060102T150405Z"))
		}
	}

	if !e.End.IsZero() {
		if e.AllDay {
			fmt.Fprintf(&b, "DTEND;VALUE=DATE:%s\n", e.End.UTC().Format("20060102"))
		} else {
			fmt.Fprintf(&b, "DTEND:%s\n", e.End.UTC().Format("20060102T150405Z"))
		}
	}

	if e.Location != "" {
		fmt.Fprintf(&b, "LOCATION:%s\n", escapeICal(e.Location))
	}
	if e.Description != "" {
		fmt.Fprintf(&b, "DESCRIPTION:%s\n", escapeICal(e.Description))
	}
	if e.URL != "" {
		// URLs are specified by RFC 5545 as being formatted per RFC 3986 and
		// not subject to text escaping.
		fmt.Fprintf(&b, "URL:%s\n", e.URL)
	}

	b.WriteString("END:VEVENT\n")
	b.WriteString("END:VCALENDAR")

	return b.String()
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
