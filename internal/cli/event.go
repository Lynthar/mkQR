package cli

import (
	"fmt"
	"time"

	"github.com/Lynthar/mkQR/pkg/encoder"
	"github.com/spf13/cobra"
)

var (
	eventSummary     string
	eventStart       string
	eventEnd         string
	eventLocation    string
	eventDescription string
	eventURL         string
	eventAllDay      bool
)

var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "Generate QR code for a calendar event (iCalendar)",
	Long: `Generate a QR code containing a calendar event.

When scanned, the event can be added to the user's calendar. Outputs a
full VCALENDAR-wrapped VEVENT for maximum scanner compatibility.

Time formats accepted for --start and --end:
  2026-05-01T10:00:00Z        RFC 3339, UTC
  2026-05-01T10:00:00+08:00   RFC 3339, with explicit offset
  2026-05-01T10:00:00         Naive datetime (interpreted as local time)
  2026-05-01                  Date only (useful with --all-day)

Examples:
  mkqr event -s "Sync meeting" --start 2026-05-01T10:00:00Z --end 2026-05-01T11:00:00Z
  mkqr event -s "Holiday" --start 2026-05-01 --all-day
  mkqr event -s "Conference" --start 2026-05-10T09:00:00 --end 2026-05-12T18:00:00 -L "Shanghai"
  mkqr event -s "Call" --start 2026-05-01T10:00:00 --url "https://meet.example/123"`,
	RunE: runEvent,
}

func init() {
	eventCmd.Flags().StringVarP(&eventSummary, "summary", "s", "", "Event title/summary [required]")
	eventCmd.Flags().StringVar(&eventStart, "start", "", "Start time [required]")
	eventCmd.Flags().StringVar(&eventEnd, "end", "", "End time")
	// Note: `-l` is claimed by the root `--level` persistent flag; use `-L`.
	eventCmd.Flags().StringVarP(&eventLocation, "location", "L", "", "Location")
	eventCmd.Flags().StringVarP(&eventDescription, "description", "d", "", "Description / notes")
	eventCmd.Flags().StringVar(&eventURL, "url", "", "URL associated with the event")
	eventCmd.Flags().BoolVar(&eventAllDay, "all-day", false, "All-day event (DTSTART/DTEND emitted as dates)")

	eventCmd.MarkFlagRequired("summary")
	eventCmd.MarkFlagRequired("start")

	rootCmd.AddCommand(eventCmd)
}

func runEvent(cmd *cobra.Command, args []string) error {
	start, err := parseEventTime(eventStart)
	if err != nil {
		return fmt.Errorf("--start: %w", err)
	}

	var end time.Time
	if eventEnd != "" {
		end, err = parseEventTime(eventEnd)
		if err != nil {
			return fmt.Errorf("--end: %w", err)
		}
		if !end.After(start) {
			return fmt.Errorf("--end (%s) must be after --start (%s)", eventEnd, eventStart)
		}
	}

	ev := &encoder.Event{
		Summary:     eventSummary,
		Start:       start,
		End:         end,
		Location:    eventLocation,
		Description: eventDescription,
		URL:         eventURL,
		AllDay:      eventAllDay,
	}

	if !quiet {
		display := start.Format("2006-01-02 15:04 MST")
		if eventAllDay {
			display = start.Format("2006-01-02") + " (all-day)"
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "Event: %s @ %s\n", eventSummary, display)
	}

	return generateQR(ev.Encode())
}

// parseEventTime accepts several time formats for event start/end.
// Inputs with explicit timezone info (RFC 3339 Z or ±hh:mm) are honored;
// naive inputs are interpreted in the local timezone so users get the
// "type the wall-clock time without thinking" behavior.
func parseEventTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	localLayouts := []string{
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, layout := range localLayouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time %q (examples: 2026-05-01T10:00:00Z, 2026-05-01T10:00:00+08:00, 2026-05-01)", s)
}
