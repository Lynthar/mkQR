package encoder

import (
	"fmt"
	"net/url"
	"strings"
)

// OTPType represents the type of OTP
type OTPType string

const (
	TOTP OTPType = "totp" // Time-based OTP
	HOTP OTPType = "hotp" // Counter-based OTP
)

// OTP encodes one-time password configuration for authenticator apps
type OTP struct {
	Type      OTPType
	Secret    string
	Issuer    string
	Account   string
	Algorithm string // SHA1, SHA256, SHA512 (default: SHA1)
	Digits    int    // 6 or 8 (default: 6)
	Period    int    // For TOTP: seconds (default: 30)
	Counter   int    // For HOTP: initial counter value
}

// Encode returns the otpauth:// URL
// Format: otpauth://TYPE/LABEL?PARAMETERS
func (o *OTP) Encode() string {
	otpType := o.Type
	if otpType == "" {
		otpType = TOTP
	}

	// Build label: "Issuer:Account" or just "Account"
	var label string
	if o.Issuer != "" {
		label = url.PathEscape(o.Issuer) + ":" + url.PathEscape(o.Account)
	} else {
		label = url.PathEscape(o.Account)
	}

	// Build parameters
	var params []string
	params = append(params, "secret="+strings.ToUpper(o.Secret))

	if o.Issuer != "" {
		params = append(params, "issuer="+url.QueryEscape(o.Issuer))
	}

	if o.Algorithm != "" && o.Algorithm != "SHA1" {
		params = append(params, "algorithm="+o.Algorithm)
	}

	if o.Digits != 0 && o.Digits != 6 {
		params = append(params, fmt.Sprintf("digits=%d", o.Digits))
	}

	if otpType == TOTP && o.Period != 0 && o.Period != 30 {
		params = append(params, fmt.Sprintf("period=%d", o.Period))
	}

	if otpType == HOTP {
		params = append(params, fmt.Sprintf("counter=%d", o.Counter))
	}

	return fmt.Sprintf("otpauth://%s/%s?%s", otpType, label, strings.Join(params, "&"))
}
