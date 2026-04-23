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
		// URL with port
		{"example.com:8080", TypeURL},
		{"example.com:8080/path", TypeURL},
		// URL with naked query / fragment
		{"example.com?q=1", TypeURL},
		{"example.com?q=hello&x=2", TypeURL},
		{"example.com#anchor", TypeURL},
		// Multi-label subdomains with hyphens
		{"my-site.example.com", TypeURL},
		{"a.my-site.com", TypeURL},
		{"deep.sub.domain.example.co.uk", TypeURL},

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
		// Intentionally NOT auto-detected as URL (user must add http:// to
		// force). Bare IPv4 / version strings are too easily confused with
		// plain content.
		{"127.0.0.1", TypeText},
		{"192.168.1.1", TypeText},
		{"1.2.3", TypeText},
		// Rejected URL-ish inputs
		{"example.c", TypeText},      // TLD must be ≥2 letters
		{"example.123", TypeText},    // TLD must be alphabetic
		{"example.com:abc", TypeText}, // port must be digits
		{".example.com", TypeText},   // leading dot (no label before it)
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
