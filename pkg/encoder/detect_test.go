package encoder

import "testing"

func TestDetect(t *testing.T) {
	tests := []struct {
		input    string
		expected ContentType
	}{
		// URLs
		{"https://example.com", TypeURL},
		{"http://example.com", TypeURL},
		{"HTTPS://EXAMPLE.COM", TypeURL},
		{"example.com", TypeURL},
		{"github.com/user/repo", TypeURL},

		// Proxy protocols
		{"vmess://base64content", TypeProxy},
		{"vless://uuid@host:port", TypeProxy},
		{"ss://base64", TypeProxy},
		{"ssr://base64", TypeProxy},
		{"trojan://password@host:port", TypeProxy},
		{"hysteria://host:port", TypeProxy},
		{"hysteria2://host:port", TypeProxy},

		// WiFi
		{"WIFI:T:WPA;S:MyNetwork;P:password;;", TypeWiFi},
		{"wifi:T:WPA;S:test;;", TypeWiFi},

		// OTP
		{"otpauth://totp/GitHub:user?secret=ABC", TypeOTP},
		{"otpauth://hotp/Service:user?secret=ABC&counter=0", TypeOTP},

		// Email
		{"mailto:test@example.com", TypeEmail},
		{"test@example.com", TypeEmail},
		{"user.name+tag@example.co.uk", TypeEmail},

		// Phone
		{"tel:+1234567890", TypePhone},

		// SMS
		{"sms:+1234567890", TypeSMS},
		{"smsto:+1234567890", TypeSMS},

		// Geo
		{"geo:40.7128,-74.0060", TypeGeo},

		// vCard
		{"BEGIN:VCARD\nVERSION:3.0\nFN:John Doe\nEND:VCARD", TypeVCard},
		{"begin:vcard\nversion:3.0", TypeVCard},

		// Calendar event
		{"BEGIN:VEVENT\nSUMMARY:Meeting\nEND:VEVENT", TypeEvent},

		// Plain text
		{"Hello World", TypeText},
		{"Some random text", TypeText},
		{"12345", TypeText},
	}

	for _, tt := range tests {
		t.Run(tt.input[:min(len(tt.input), 30)], func(t *testing.T) {
			result := Detect(tt.input)
			if result != tt.expected {
				t.Errorf("Detect(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDetectAndDescribe(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    ContentType
		expectedDescNot string // description should NOT be empty
	}{
		{"https://example.com", TypeURL, ""},
		{"test@example.com", TypeEmail, ""},
		{"Hello World", TypeText, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			contentType, description := DetectAndDescribe(tt.input)
			if contentType != tt.expectedType {
				t.Errorf("DetectAndDescribe(%q) type = %q, want %q", tt.input, contentType, tt.expectedType)
			}
			if description == "" {
				t.Errorf("DetectAndDescribe(%q) description is empty", tt.input)
			}
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
