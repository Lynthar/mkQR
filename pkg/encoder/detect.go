package encoder

import (
	"regexp"
	"strings"
)

// Package-level regex compiled once at init time. Previously these were
// compiled inside Detect(), which meant every batch line re-did the work.
//
// urlPattern matches bare URLs (no scheme): one-or-more domain labels (1-63
// chars, alnum + inner hyphens) followed by a ≥2-letter TLD, plus optional
// :port, /path, ?query, or #fragment. The alphabetic-TLD requirement keeps
// IPv4 literals ("127.0.0.1") and version strings ("1.2.3") from auto-
// detecting as URLs; users who want those as URLs should prefix http://
// explicitly. Explicit schemes (http://, https://, proxy://, etc.) are
// caught by the prefix checks in Detect() before this pattern is consulted.
var (
	urlPattern   = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}(:\d{1,5})?([/?#].*)?$`)
	emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// Detect analyzes input and returns its content type
func Detect(input string) ContentType {
	input = strings.TrimSpace(input)

	// Check for specific protocols/prefixes
	lowerInput := strings.ToLower(input)

	// Proxy protocols
	proxyPrefixes := []string{
		"vmess://", "vless://", "ss://", "ssr://",
		"trojan://", "hysteria://", "hysteria2://",
		"tuic://", "naive+https://", "nekoray://",
	}
	for _, prefix := range proxyPrefixes {
		if strings.HasPrefix(lowerInput, prefix) {
			return TypeProxy
		}
	}

	// URL
	if strings.HasPrefix(lowerInput, "http://") || strings.HasPrefix(lowerInput, "https://") {
		return TypeURL
	}

	// WiFi
	if strings.HasPrefix(strings.ToUpper(input), "WIFI:") {
		return TypeWiFi
	}

	// OTP
	if strings.HasPrefix(lowerInput, "otpauth://") {
		return TypeOTP
	}

	// Email (mailto:)
	if strings.HasPrefix(lowerInput, "mailto:") {
		return TypeEmail
	}

	// Phone (tel:)
	if strings.HasPrefix(lowerInput, "tel:") {
		return TypePhone
	}

	// SMS
	if strings.HasPrefix(lowerInput, "sms:") || strings.HasPrefix(lowerInput, "smsto:") {
		return TypeSMS
	}

	// Geo
	if strings.HasPrefix(lowerInput, "geo:") {
		return TypeGeo
	}

	// vCard
	if strings.HasPrefix(strings.ToUpper(input), "BEGIN:VCARD") {
		return TypeVCard
	}

	// Calendar event
	if strings.HasPrefix(strings.ToUpper(input), "BEGIN:VEVENT") {
		return TypeEvent
	}

	// Check if it looks like a URL without protocol
	if urlPattern.MatchString(input) {
		return TypeURL
	}

	// Check if it looks like an email address
	if emailPattern.MatchString(input) {
		return TypeEmail
	}

	// Default to plain text
	return TypeText
}

// DetectAndDescribe returns the detected type and a human-readable description
func DetectAndDescribe(input string) (ContentType, string) {
	contentType := Detect(input)

	descriptions := map[ContentType]string{
		TypeText:    "Plain text",
		TypeURL:     "URL",
		TypeWiFi:    "WiFi network",
		TypeVCard:   "Contact card (vCard)",
		TypeEmail:   "Email",
		TypePhone:   "Phone number",
		TypeSMS:     "SMS message",
		TypeOTP:     "One-time password (2FA)",
		TypeGeo:     "Geographic location",
		TypeEvent:   "Calendar event",
		TypeProxy:   "Proxy configuration",
		TypeUnknown: "Unknown",
	}

	return contentType, descriptions[contentType]
}
