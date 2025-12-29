package encoder

import (
	"fmt"
	"strings"
)

// WiFiEncryption represents WiFi encryption type
type WiFiEncryption string

const (
	WPA    WiFiEncryption = "WPA"
	WEP    WiFiEncryption = "WEP"
	NoPass WiFiEncryption = "nopass"
)

// ParseWiFiEncryption parses encryption type from string
func ParseWiFiEncryption(s string) (WiFiEncryption, error) {
	switch strings.ToUpper(s) {
	case "WPA", "WPA2", "WPA3", "WPA2/WPA3":
		return WPA, nil
	case "WEP":
		return WEP, nil
	case "NOPASS", "NONE", "OPEN", "":
		return NoPass, nil
	default:
		return "", fmt.Errorf("unknown encryption type: %s (use WPA, WEP, or nopass)", s)
	}
}

// WiFi encodes WiFi network information
type WiFi struct {
	SSID       string
	Password   string
	Encryption WiFiEncryption
	Hidden     bool
}

// Encode returns the WiFi QR code format string
// Format: WIFI:T:<encryption>;S:<ssid>;P:<password>;H:<hidden>;;
func (w *WiFi) Encode() string {
	// Escape special characters in SSID and password
	ssid := escapeWiFiString(w.SSID)
	password := escapeWiFiString(w.Password)

	encryption := w.Encryption
	if encryption == "" {
		if w.Password == "" {
			encryption = NoPass
		} else {
			encryption = WPA
		}
	}

	hidden := ""
	if w.Hidden {
		hidden = "H:true;"
	}

	return fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;%s;", encryption, ssid, password, hidden)
}

// escapeWiFiString escapes special characters for WiFi QR format
func escapeWiFiString(s string) string {
	// Characters that need escaping: \ ; , " :
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`;`, `\;`,
		`,`, `\,`,
		`"`, `\"`,
		`:`, `\:`,
	)
	return replacer.Replace(s)
}
