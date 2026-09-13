package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Pins the notice to the payload: the type in "WiFi: <ssid> (<type>)" must be
// the T: field of the QR actually written, whether given with -e or inferred
// from -p. The PNG is compared against one generated from the literal payload.
func TestWifiNoticeMatchesPayloadType(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		typ     string
		payload string
	}{
		{"password infers WPA", []string{"-p", "pw"}, "WPA", "WIFI:T:WPA;S:Net;P:pw;;"},
		{"no password infers nopass", nil, "nopass", "WIFI:T:nopass;S:Net;P:;;"},
		{"explicit WEP", []string{"-p", "pw", "-e", "wep"}, "WEP", "WIFI:T:WEP;S:Net;P:pw;;"},
		{"explicit nopass keeps password", []string{"-p", "pw", "-e", "nopass"}, "nopass", "WIFI:T:nopass;S:Net;P:pw;;"},
	}

	origPassword, origEncryption, origOut := wifiPassword, wifiEncryption, outputFile
	t.Cleanup(func() { wifiPassword, wifiEncryption, outputFile = origPassword, origEncryption, origOut })
	silenceStderr(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Flags bind to package globals that survive Execute; omitted
			// flags must start from their defaults, not the previous case.
			wifiPassword, wifiEncryption = "", ""
			out := filepath.Join(t.TempDir(), "got.png")

			var buf bytes.Buffer
			cmd := GetRootCmd()
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs(append([]string{"wifi", "-s", "Net", "-o", out}, tt.args...))
			if err := cmd.Execute(); err != nil {
				t.Fatalf("wifi %v: %v", tt.args, err)
			}

			if want := "WiFi: Net (" + tt.typ + ")"; !strings.Contains(buf.String(), want) {
				t.Errorf("notice lacks %q, output was:\n%s", want, buf.String())
			}

			got, err := os.ReadFile(out)
			if err != nil {
				t.Fatalf("read output: %v", err)
			}
			gen, err := buildGenerator(io.Discard)
			if err != nil {
				t.Fatalf("buildGenerator: %v", err)
			}
			want, err := gen.GeneratePNG(tt.payload)
			if err != nil {
				t.Fatalf("GeneratePNG: %v", err)
			}
			if !bytes.Equal(got, want) {
				t.Errorf("QR written for %v does not encode %q", tt.args, tt.payload)
			}
		})
	}
}
