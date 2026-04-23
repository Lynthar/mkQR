package encoder

import (
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
// Content lines are CRLF-delimited and folded at 75 octets per RFC 2426 §2.6
// via writeFolded.
func (v *VCard) Encode() string {
	var b strings.Builder

	writeFolded(&b, "BEGIN:VCARD")
	writeFolded(&b, "VERSION:3.0")

	// Name
	if v.FirstName != "" || v.LastName != "" {
		writeFolded(&b, "N:"+escapeVCard(v.LastName)+";"+escapeVCard(v.FirstName)+";;;")
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
		writeFolded(&b, "FN:"+escapeVCard(fn))
	}

	if v.Organization != "" {
		writeFolded(&b, "ORG:"+escapeVCard(v.Organization))
	}

	if v.Title != "" {
		writeFolded(&b, "TITLE:"+escapeVCard(v.Title))
	}

	if v.Phone != "" {
		writeFolded(&b, "TEL;TYPE=HOME:"+v.Phone)
	}
	if v.PhoneWork != "" {
		writeFolded(&b, "TEL;TYPE=WORK:"+v.PhoneWork)
	}
	if v.PhoneMobile != "" {
		writeFolded(&b, "TEL;TYPE=CELL:"+v.PhoneMobile)
	}

	// Email. RFC 2426 §3.3.2: EMAIL default TYPE is INTERNET; emit it
	// explicitly alongside HOME/WORK so strict parsers don't fall back to
	// a non-SMTP mailbox type (e.g. X.400).
	if v.Email != "" {
		writeFolded(&b, "EMAIL;TYPE=INTERNET,HOME:"+v.Email)
	}
	if v.EmailWork != "" {
		writeFolded(&b, "EMAIL;TYPE=INTERNET,WORK:"+v.EmailWork)
	}

	// Website. URL passes through escapeVCard so that an embedded ';' (e.g.
	// matrix/session parameters like ';jsessionid=...') doesn't terminate
	// the property early and break vCard structure. Conformant parsers
	// unescape URI values on read.
	if v.Website != "" {
		writeFolded(&b, "URL:"+escapeVCard(v.Website))
	}

	if v.Address != "" {
		writeFolded(&b, "ADR:;;"+escapeVCard(v.Address)+";;;;")
	}

	if v.Note != "" {
		writeFolded(&b, "NOTE:"+escapeVCard(v.Note))
	}

	writeFolded(&b, "END:VCARD")

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
