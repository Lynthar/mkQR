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

// Encode returns the vCard format string (version 3.0)
func (v *VCard) Encode() string {
	var b strings.Builder

	b.WriteString("BEGIN:VCARD\n")
	b.WriteString("VERSION:3.0\n")

	// Name
	if v.FirstName != "" || v.LastName != "" {
		b.WriteString(fmt.Sprintf("N:%s;%s;;;\n", escapeVCard(v.LastName), escapeVCard(v.FirstName)))
		b.WriteString(fmt.Sprintf("FN:%s %s\n", escapeVCard(v.FirstName), escapeVCard(v.LastName)))
	}

	// Organization
	if v.Organization != "" {
		b.WriteString(fmt.Sprintf("ORG:%s\n", escapeVCard(v.Organization)))
	}

	// Title
	if v.Title != "" {
		b.WriteString(fmt.Sprintf("TITLE:%s\n", escapeVCard(v.Title)))
	}

	// Phone numbers
	if v.Phone != "" {
		b.WriteString(fmt.Sprintf("TEL;TYPE=HOME:%s\n", v.Phone))
	}
	if v.PhoneWork != "" {
		b.WriteString(fmt.Sprintf("TEL;TYPE=WORK:%s\n", v.PhoneWork))
	}
	if v.PhoneMobile != "" {
		b.WriteString(fmt.Sprintf("TEL;TYPE=CELL:%s\n", v.PhoneMobile))
	}

	// Email
	if v.Email != "" {
		b.WriteString(fmt.Sprintf("EMAIL;TYPE=HOME:%s\n", v.Email))
	}
	if v.EmailWork != "" {
		b.WriteString(fmt.Sprintf("EMAIL;TYPE=WORK:%s\n", v.EmailWork))
	}

	// Website
	if v.Website != "" {
		b.WriteString(fmt.Sprintf("URL:%s\n", v.Website))
	}

	// Address
	if v.Address != "" {
		b.WriteString(fmt.Sprintf("ADR:;;%s;;;;\n", escapeVCard(v.Address)))
	}

	// Note
	if v.Note != "" {
		b.WriteString(fmt.Sprintf("NOTE:%s\n", escapeVCard(v.Note)))
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
