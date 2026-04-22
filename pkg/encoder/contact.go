package encoder

import (
	"fmt"
	"net/url"
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
		params = append(params, "cc="+url.QueryEscape(e.CC))
	}
	if e.BCC != "" {
		params = append(params, "bcc="+url.QueryEscape(e.BCC))
	}
	if e.Subject != "" {
		params = append(params, "subject="+url.QueryEscape(e.Subject))
	}
	if e.Body != "" {
		params = append(params, "body="+url.QueryEscape(e.Body))
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
		result += "?body=" + url.QueryEscape(s.Body)
	}
	return result
}

// Geo encodes geographic location
type Geo struct {
	Latitude  float64
	Longitude float64
	Query     string // Optional location query/name
}

// Encode returns the geo: URL
func (g *Geo) Encode() string {
	if g.Query != "" {
		return fmt.Sprintf("geo:%f,%f?q=%s", g.Latitude, g.Longitude, url.QueryEscape(g.Query))
	}
	return fmt.Sprintf("geo:%f,%f", g.Latitude, g.Longitude)
}
