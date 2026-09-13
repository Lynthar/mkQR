package encoder

import "strings"

// textEscaper applies the TEXT-value escaping shared by iCalendar (RFC 5545
// §3.3.11) and vCard 3.0 (RFC 2426 §2.4.2): backslash, semicolon, comma, LF.
var textEscaper = strings.NewReplacer(`\`, `\\`, `;`, `\;`, `,`, `\,`, "\n", `\n`)

func escapeText(s string) string { return textEscaper.Replace(s) }
