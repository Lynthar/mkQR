package encoder

import (
	"strings"
	"testing"
)

func TestVCardEncode(t *testing.T) {
	tests := []struct {
		name     string
		vcard    VCard
		contains []string
	}{
		{
			name: "full name",
			vcard: VCard{
				FirstName: "John",
				LastName:  "Doe",
			},
			contains: []string{
				"BEGIN:VCARD",
				"VERSION:3.0",
				"N:Doe;John;;;",
				"FN:John Doe",
				"END:VCARD",
			},
		},
		{
			name: "first name only",
			vcard: VCard{
				FirstName: "John",
			},
			contains: []string{
				"N:;John;;;",
				"FN:John",
			},
		},
		{
			name: "last name only",
			vcard: VCard{
				LastName: "Doe",
			},
			contains: []string{
				"N:Doe;;;;",
				"FN:Doe",
			},
		},
		{
			name: "with organization and title",
			vcard: VCard{
				FirstName:    "Jane",
				LastName:     "Smith",
				Organization: "Acme Inc",
				Title:        "Engineer",
			},
			contains: []string{
				"ORG:Acme Inc",
				"TITLE:Engineer",
			},
		},
		{
			name: "with phone numbers",
			vcard: VCard{
				FirstName:   "Test",
				Phone:       "+1234567890",
				PhoneMobile: "+0987654321",
			},
			contains: []string{
				"TEL;TYPE=HOME:+1234567890",
				"TEL;TYPE=CELL:+0987654321",
			},
		},
		{
			name: "with email",
			vcard: VCard{
				FirstName: "Test",
				Email:     "test@example.com",
			},
			contains: []string{
				"EMAIL;TYPE=HOME:test@example.com",
			},
		},
		{
			name: "with special characters",
			vcard: VCard{
				FirstName: "John,Jr",
				LastName:  "O'Brien;Sr",
				Note:      "Line1\nLine2",
			},
			contains: []string{
				`N:O'Brien\;Sr;John\,Jr;;;`,
				`NOTE:Line1\nLine2`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.vcard.Encode()
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Encode() missing %q in result:\n%s", expected, result)
				}
			}
		})
	}
}

func TestVCardNoExtraSpaceInFN(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		wantFN    string
	}{
		{"both names", "John", "Doe", "FN:John Doe"},
		{"first only", "John", "", "FN:John"},
		{"last only", "", "Doe", "FN:Doe"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vcard := VCard{
				FirstName: tt.firstName,
				LastName:  tt.lastName,
			}
			result := vcard.Encode()
			if !strings.Contains(result, tt.wantFN) {
				t.Errorf("Expected %q in result, got:\n%s", tt.wantFN, result)
			}
			// Make sure there's no double space or trailing space in FN
			if strings.Contains(result, "FN: ") || strings.Contains(result, "FN:  ") {
				t.Errorf("FN has leading space in result:\n%s", result)
			}
		})
	}
}

func TestEscapeVCard(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{`with\backslash`, `with\\backslash`},
		{"with,comma", `with\,comma`},
		{"with;semicolon", `with\;semicolon`},
		{"with\nnewline", `with\nnewline`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapeVCard(tt.input)
			if result != tt.expected {
				t.Errorf("escapeVCard(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
