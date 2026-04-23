package encoder

import (
	"fmt"
	"strings"
)

// VCard encodes contact information in vCard format
type VCard struct {
	FirstName    string
	LastName     string
	Organization string
	Title        string
	Phone        string
	PhoneWork    string
	PhoneMobile  string
	Email        string
	EmailWork    string
	Website      string
	Address      string
	Note         string
}

// Encode returns the vCard format string (version 3.0).
// Content lines are CRLF-delimited per RFC 2426 §2.4.2.
func (v *VCard) Encode() string {
	var b strings.Builder

	b.WriteString("BEGIN:VCARD\r\n")
	b.WriteString("VERSION:3.0\r\n")

	// Name
	if v.FirstName != "" || v.LastName != "" {
		b.WriteString(fmt.Sprintf("N:%s;%s;;;\r\n", escapeVCard(v.LastName), escapeVCard(v.FirstName)))
		// Build FN (formatted name) properly to avoid extra spaces
		var fn string
		switch {
		case v.FirstName != "" && v.LastName != "":
			fn = v.FirstName + " " + v.LastName
		case v.FirstName != "":
			fn = v.FirstName
		default:
			fn = v.LastName
		}
		b.WriteString(fmt.Sprintf("FN:%s\r\n", escapeVCard(fn)))
	}

	// Organization
	if v.Organization != "" {
		b.WriteString(fmt.Sprintf("ORG:%s\r\n", escapeVCard(v.Organization)))
	}

	// Title
	if v.Title != "" {
		b.WriteString(fmt.Sprintf("TITLE:%s\r\n", escapeVCard(v.Title)))
	}

	// Phone numbers
	if v.Phone != "" {
		b.WriteString(fmt.Sprintf("TEL;TYPE=HOME:%s\r\n", v.Phone))
	}
	if v.PhoneWork != "" {
		b.WriteString(fmt.Sprintf("TEL;TYPE=WORK:%s\r\n", v.PhoneWork))
	}
	if v.PhoneMobile != "" {
		b.WriteString(fmt.Sprintf("TEL;TYPE=CELL:%s\r\n", v.PhoneMobile))
	}

	// Email
	if v.Email != "" {
		b.WriteString(fmt.Sprintf("EMAIL;TYPE=HOME:%s\r\n", v.Email))
	}
	if v.EmailWork != "" {
		b.WriteString(fmt.Sprintf("EMAIL;TYPE=WORK:%s\r\n", v.EmailWork))
	}

	// Website
	if v.Website != "" {
		b.WriteString(fmt.Sprintf("URL:%s\r\n", v.Website))
	}

	// Address
	if v.Address != "" {
		b.WriteString(fmt.Sprintf("ADR:;;%s;;;;\r\n", escapeVCard(v.Address)))
	}

	// Note
	if v.Note != "" {
		b.WriteString(fmt.Sprintf("NOTE:%s\r\n", escapeVCard(v.Note)))
	}

	b.WriteString("END:VCARD")

	return b.String()
}

// escapeVCard escapes special characters for vCard format
func escapeVCard(s string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`,`, `\,`,
		`;`, `\;`,
		"\n", `\n`,
	)
	return replacer.Replace(s)
}
