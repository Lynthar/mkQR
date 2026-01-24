package encoder

import (
	"strings"
	"testing"
)

func TestOTPEncode(t *testing.T) {
	tests := []struct {
		name     string
		otp      OTP
		contains []string
	}{
		{
			name: "basic TOTP",
			otp: OTP{
				Secret:  "JBSWY3DPEHPK3PXP",
				Issuer:  "GitHub",
				Account: "user@example.com",
			},
			contains: []string{
				"otpauth://totp/",
				"GitHub:user%40example.com",
				"secret=JBSWY3DPEHPK3PXP",
				"issuer=GitHub",
			},
		},
		{
			name: "HOTP with counter",
			otp: OTP{
				Type:    HOTP,
				Secret:  "ABCD1234",
				Issuer:  "Service",
				Account: "user",
				Counter: 10,
			},
			contains: []string{
				"otpauth://hotp/",
				"counter=10",
			},
		},
		{
			name: "custom algorithm and digits",
			otp: OTP{
				Secret:    "SECRET",
				Issuer:    "App",
				Account:   "user",
				Algorithm: "SHA256",
				Digits:    8,
			},
			contains: []string{
				"algorithm=SHA256",
				"digits=8",
			},
		},
		{
			name: "custom period",
			otp: OTP{
				Secret:  "SECRET",
				Issuer:  "App",
				Account: "user",
				Period:  60,
			},
			contains: []string{
				"period=60",
			},
		},
		{
			name: "default values not included",
			otp: OTP{
				Secret:    "SECRET",
				Issuer:    "App",
				Account:   "user",
				Algorithm: "SHA1",
				Digits:    6,
				Period:    30,
			},
			contains: []string{
				"secret=SECRET",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.otp.Encode()
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Encode() missing %q in result: %s", expected, result)
				}
			}
		})
	}
}

func TestOTPEncodeDefaultsNotIncluded(t *testing.T) {
	otp := OTP{
		Secret:    "SECRET",
		Issuer:    "App",
		Account:   "user",
		Algorithm: "SHA1",
		Digits:    6,
		Period:    30,
	}
	result := otp.Encode()

	// These defaults should NOT appear in the URL
	if strings.Contains(result, "algorithm=") {
		t.Errorf("Default algorithm=SHA1 should not be in URL: %s", result)
	}
	if strings.Contains(result, "digits=") {
		t.Errorf("Default digits=6 should not be in URL: %s", result)
	}
	if strings.Contains(result, "period=") {
		t.Errorf("Default period=30 should not be in URL: %s", result)
	}
}

func TestValidateSecret(t *testing.T) {
	tests := []struct {
		secret   string
		hasError bool
	}{
		{"JBSWY3DPEHPK3PXP", false},
		{"ABCD1234", false},
		{"abcd1234", false},
		{"ABCD 1234", false},       // spaces allowed (removed)
		{"ABCD-1234-EFGH", false},  // hyphens allowed (removed)
		{"", true},                 // empty
		{"ABCD1890", true},         // invalid chars (8, 9, 0)
		{"ABCD!@#$", true},         // special chars
		{"12345678", true},         // all invalid digits
	}

	for _, tt := range tests {
		t.Run(tt.secret, func(t *testing.T) {
			err := ValidateSecret(tt.secret)
			if tt.hasError && err == nil {
				t.Errorf("ValidateSecret(%q) expected error, got nil", tt.secret)
			}
			if !tt.hasError && err != nil {
				t.Errorf("ValidateSecret(%q) unexpected error: %v", tt.secret, err)
			}
		})
	}
}
