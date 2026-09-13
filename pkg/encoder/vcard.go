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
		writeFolded(&b, "N:"+escapeText(v.LastName)+";"+escapeText(v.FirstName)+";;;")
		writeFolded(&b, "FN:"+escapeText(v.FormattedName()))
	}

	if v.Organization != "" {
		writeFolded(&b, "ORG:"+escapeText(v.Organization))
	}

	if v.Title != "" {
		writeFolded(&b, "TITLE:"+escapeText(v.Title))
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

	// Website. URL passes through escapeText so that an embedded ';' (e.g.
	// matrix/session parameters like ';jsessionid=...') doesn't terminate
	// the property early and break vCard structure. Conformant parsers
	// unescape URI values on read.
	if v.Website != "" {
		writeFolded(&b, "URL:"+escapeText(v.Website))
	}

	if v.Address != "" {
		writeFolded(&b, "ADR:;;"+escapeText(v.Address)+";;;;")
	}

	if v.Note != "" {
		writeFolded(&b, "NOTE:"+escapeText(v.Note))
	}

	writeFolded(&b, "END:VCARD")

	return b.String()
}

// FormattedName is the FN value Encode emits: "First Last", or whichever of
// the two is set; empty when neither is.
func (v *VCard) FormattedName() string {
	switch {
	case v.FirstName != "" && v.LastName != "":
		return v.FirstName + " " + v.LastName
	case v.FirstName != "":
		return v.FirstName
	default:
		return v.LastName
	}
}
