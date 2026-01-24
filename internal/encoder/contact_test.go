package encoder

import (
	"strings"
	"testing"
)

func TestEmailEncode(t *testing.T) {
	tests := []struct {
		name     string
		email    Email
		expected string
	}{
		{
			name:     "simple email",
			email:    Email{To: "test@example.com"},
			expected: "mailto:test@example.com",
		},
		{
			name: "with subject",
			email: Email{
				To:      "test@example.com",
				Subject: "Hello World",
			},
			expected: "mailto:test@example.com?subject=Hello+World",
		},
		{
			name: "with body",
			email: Email{
				To:   "test@example.com",
				Body: "Message body",
			},
			expected: "mailto:test@example.com?body=Message+body",
		},
		{
			name: "with CC and BCC",
			email: Email{
				To:  "test@example.com",
				CC:  "cc@example.com",
				BCC: "bcc@example.com",
			},
			expected: "mailto:test@example.com?cc=cc%40example.com&bcc=bcc%40example.com",
		},
		{
			name: "full email",
			email: Email{
				To:      "test@example.com",
				CC:      "cc@example.com",
				Subject: "Test",
				Body:    "Body",
			},
			expected: "mailto:test@example.com?cc=cc%40example.com&subject=Test&body=Body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.email.Encode()
			if result != tt.expected {
				t.Errorf("Encode() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestPhoneEncode(t *testing.T) {
	tests := []struct {
		name     string
		phone    Phone
		expected string
	}{
		{
			name:     "simple number",
			phone:    Phone{Number: "+1234567890"},
			expected: "tel:+1234567890",
		},
		{
			name:     "number with spaces",
			phone:    Phone{Number: "+1 234 567 890"},
			expected: "tel:+1234567890",
		},
		{
			name:     "local number",
			phone:    Phone{Number: "1234567890"},
			expected: "tel:1234567890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.phone.Encode()
			if result != tt.expected {
				t.Errorf("Encode() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSMSEncode(t *testing.T) {
	tests := []struct {
		name     string
		sms      SMS
		expected string
	}{
		{
			name:     "simple SMS",
			sms:      SMS{Number: "+1234567890"},
			expected: "sms:+1234567890",
		},
		{
			name: "SMS with body",
			sms: SMS{
				Number: "+1234567890",
				Body:   "Hello!",
			},
			expected: "sms:+1234567890?body=Hello%21",
		},
		{
			name: "number with spaces",
			sms: SMS{
				Number: "+1 234 567 890",
				Body:   "Test",
			},
			expected: "sms:+1234567890?body=Test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sms.Encode()
			if result != tt.expected {
				t.Errorf("Encode() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGeoEncode(t *testing.T) {
	tests := []struct {
		name     string
		geo      Geo
		contains []string
	}{
		{
			name: "simple coordinates",
			geo: Geo{
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			contains: []string{"geo:40.712800,-74.006000"},
		},
		{
			name: "with query",
			geo: Geo{
				Latitude:  39.9042,
				Longitude: 116.4074,
				Query:     "Beijing",
			},
			contains: []string{"geo:39.904200,116.407400?q=Beijing"},
		},
		{
			name: "query with spaces",
			geo: Geo{
				Latitude:  31.2304,
				Longitude: 121.4737,
				Query:     "Shanghai Tower",
			},
			contains: []string{"geo:31.230400,121.473700?q=Shanghai+Tower"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.geo.Encode()
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Encode() = %q, want to contain %q", result, expected)
				}
			}
		})
	}
}
