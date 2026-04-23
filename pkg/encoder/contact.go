package encoder

import (
	"fmt"
	"strconv"
	"strings"
)

// Email encodes email information
type Email struct {
	To      string
	CC      string
	BCC     string
	Subject string
	Body    string
}

// Encode returns the mailto: URL
func (e *Email) Encode() string {
	var params []string

	if e.CC != "" {
		params = append(params, "cc="+percentEscape(e.CC))
	}
	if e.BCC != "" {
		params = append(params, "bcc="+percentEscape(e.BCC))
	}
	if e.Subject != "" {
		params = append(params, "subject="+percentEscape(e.Subject))
	}
	if e.Body != "" {
		params = append(params, "body="+percentEscape(e.Body))
	}

	result := "mailto:" + e.To
	if len(params) > 0 {
		result += "?" + strings.Join(params, "&")
	}

	return result
}

// Phone encodes phone number for dialing
type Phone struct {
	Number string
}

// Encode returns the tel: URL
func (p *Phone) Encode() string {
	// Remove spaces and format
	number := strings.ReplaceAll(p.Number, " ", "")
	return "tel:" + number
}

// SMS encodes SMS message
type SMS struct {
	Number string
	Body   string
}

// Encode returns the sms: URL
func (s *SMS) Encode() string {
	number := strings.ReplaceAll(s.Number, " ", "")
	result := "sms:" + number
	if s.Body != "" {
		result += "?body=" + percentEscape(s.Body)
	}
	return result
}

// Geo encodes geographic location
type Geo struct {
	Latitude  float64
	Longitude float64
	Query     string // Optional location query/name
}

// Encode returns the geo: URL.
// Coordinates use FormatFloat with precision -1 so simple values don't pad
// trailing zeros ("40.0000" was a waste of QR payload) and high-precision
// GPS data isn't truncated at the %f-default 6 decimals.
func (g *Geo) Encode() string {
	lat := strconv.FormatFloat(g.Latitude, 'f', -1, 64)
	lng := strconv.FormatFloat(g.Longitude, 'f', -1, 64)
	if g.Query != "" {
		return fmt.Sprintf("geo:%s,%s?q=%s", lat, lng, percentEscape(g.Query))
	}
	return "geo:" + lat + "," + lng
}
