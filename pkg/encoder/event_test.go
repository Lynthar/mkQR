package encoder

import (
	"strings"
	"testing"
	"time"
)

func TestEventEncode(t *testing.T) {
	utc := time.UTC
	tests := []struct {
		name     string
		event    Event
		contains []string
	}{
		{
			name: "minimal timed event",
			event: Event{
				Summary: "Meeting",
				Start:   time.Date(2026, 5, 1, 10, 0, 0, 0, utc),
				End:     time.Date(2026, 5, 1, 11, 0, 0, 0, utc),
			},
			contains: []string{
				"BEGIN:VCALENDAR",
				"VERSION:2.0",
				"BEGIN:VEVENT",
				"SUMMARY:Meeting",
				"DTSTART:20260501T100000Z",
				"DTEND:20260501T110000Z",
				"END:VEVENT",
				"END:VCALENDAR",
			},
		},
		{
			name: "all-day event",
			event: Event{
				Summary: "Holiday",
				Start:   time.Date(2026, 5, 1, 0, 0, 0, 0, utc),
				AllDay:  true,
			},
			contains: []string{
				"DTSTART;VALUE=DATE:20260501",
			},
		},
		{
			name: "all-day range",
			event: Event{
				Summary: "Conference",
				Start:   time.Date(2026, 5, 10, 0, 0, 0, 0, utc),
				End:     time.Date(2026, 5, 12, 0, 0, 0, 0, utc),
				AllDay:  true,
			},
			contains: []string{
				"DTSTART;VALUE=DATE:20260510",
				"DTEND;VALUE=DATE:20260512",
			},
		},
		{
			name: "event with location and description",
			event: Event{
				Summary:     "Team sync",
				Start:       time.Date(2026, 5, 1, 10, 0, 0, 0, utc),
				Location:    "Conference Room A",
				Description: "Discuss Q3 plans",
			},
			contains: []string{
				"LOCATION:Conference Room A",
				"DESCRIPTION:Discuss Q3 plans",
			},
		},
		{
			name: "special characters escaped",
			event: Event{
				Summary:     "Meeting, take 2",
				Start:       time.Date(2026, 5, 1, 10, 0, 0, 0, utc),
				Description: "Line1\nLine2; with semicolons",
			},
			contains: []string{
				`SUMMARY:Meeting\, take 2`,
				`DESCRIPTION:Line1\nLine2\; with semicolons`,
			},
		},
		{
			name: "timezone is converted to UTC in output",
			event: Event{
				Summary: "TZ test",
				// 18:00 in UTC+8 → 10:00Z
				Start: time.Date(2026, 5, 1, 18, 0, 0, 0, time.FixedZone("CST", 8*3600)),
			},
			contains: []string{
				"DTSTART:20260501T100000Z",
			},
		},
		{
			name: "URL preserved without iCal escaping",
			event: Event{
				Summary: "WebEx",
				Start:   time.Date(2026, 5, 1, 10, 0, 0, 0, utc),
				URL:     "https://example.com/meet?id=abc,123",
			},
			contains: []string{
				"URL:https://example.com/meet?id=abc,123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.event.Encode()
			for _, exp := range tt.contains {
				if !strings.Contains(result, exp) {
					t.Errorf("Encode() missing %q in output:\n%s", exp, result)
				}
			}
		})
	}
}

func TestEscapeICal(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"simple", "simple"},
		{`with\backslash`, `with\\backslash`},
		{"with,comma", `with\,comma`},
		{"with;semi", `with\;semi`},
		{"with\nnewline", `with\nnewline`},
		{`all\;,chars\nplus`, `all\\\;\,chars\\nplus`},
	}
	for _, tt := range tests {
		if got := escapeICal(tt.in); got != tt.want {
			t.Errorf("escapeICal(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
