package encoder

import (
	"strings"
	"testing"
)

func TestWiFiEncode(t *testing.T) {
	tests := []struct {
		name     string
		wifi     WiFi
		expected string
	}{
		{
			name: "basic WPA network",
			wifi: WiFi{
				SSID:       "MyNetwork",
				Password:   "password123",
				Encryption: WPA,
			},
			expected: "WIFI:T:WPA;S:MyNetwork;P:password123;;",
		},
		{
			name: "open network",
			wifi: WiFi{
				SSID:       "FreeWiFi",
				Encryption: NoPass,
			},
			expected: "WIFI:T:nopass;S:FreeWiFi;P:;;",
		},
		{
			name: "hidden network",
			wifi: WiFi{
				SSID:       "HiddenNet",
				Password:   "secret",
				Encryption: WPA,
				Hidden:     true,
			},
			expected: "WIFI:T:WPA;S:HiddenNet;P:secret;H:true;;",
		},
		{
			name: "SSID with special characters",
			wifi: WiFi{
				SSID:       "My;Network:Test",
				Password:   "pass;word",
				Encryption: WPA,
			},
			expected: `WIFI:T:WPA;S:My\;Network\:Test;P:pass\;word;;`,
		},
		{
			name: "auto-detect WPA when password provided",
			wifi: WiFi{
				SSID:     "AutoNet",
				Password: "mypass",
			},
			expected: "WIFI:T:WPA;S:AutoNet;P:mypass;;",
		},
		{
			name: "auto-detect nopass when no password",
			wifi: WiFi{
				SSID: "OpenNet",
			},
			expected: "WIFI:T:nopass;S:OpenNet;P:;;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.wifi.Encode()
			if result != tt.expected {
				t.Errorf("Encode() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestParseWiFiEncryption(t *testing.T) {
	tests := []struct {
		input    string
		expected WiFiEncryption
		hasError bool
	}{
		{"WPA", WPA, false},
		{"wpa", WPA, false},
		{"WPA2", WPA, false},
		{"WPA3", WPA, false},
		{"WEP", WEP, false},
		{"wep", WEP, false},
		{"nopass", NoPass, false},
		{"NOPASS", NoPass, false},
		{"none", NoPass, false},
		{"open", NoPass, false},
		{"", NoPass, false},
		{"invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseWiFiEncryption(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("ParseWiFiEncryption(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseWiFiEncryption(%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ParseWiFiEncryption(%q) = %q, want %q", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestEscapeWiFiString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{`with\backslash`, `with\\backslash`},
		{"with;semicolon", `with\;semicolon`},
		{"with,comma", `with\,comma`},
		{`with"quote`, `with\"quote`},
		{"with:colon", `with\:colon`},
		{`all\;,":chars`, `all\\\;\,\"\:chars`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapeWiFiString(tt.input)
			if result != tt.expected {
				t.Errorf("escapeWiFiString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
