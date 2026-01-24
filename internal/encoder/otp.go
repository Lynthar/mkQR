package encoder

import (
	"fmt"
	"net/url"
	"regexp"
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

// base32Pattern matches valid base32 characters (A-Z and 2-7, case insensitive)
var base32Pattern = regexp.MustCompile(`^[A-Za-z2-7]+=*$`)

// ValidateSecret checks if the secret is a valid base32 string
func ValidateSecret(secret string) error {
	// Remove spaces and hyphens (common in user-provided secrets)
	cleaned := strings.ReplaceAll(strings.ReplaceAll(secret, " ", ""), "-", "")
	if cleaned == "" {
		return fmt.Errorf("secret cannot be empty")
	}
	if !base32Pattern.MatchString(cleaned) {
		return fmt.Errorf("secret must be a valid base32 string (A-Z, 2-7)")
	}
	return nil
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
