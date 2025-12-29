package cli

import (
	"fmt"

	"github.com/Lynthar/mkQR/internal/encoder"
	"github.com/spf13/cobra"
)

var (
	wifiSSID       string
	wifiPassword   string
	wifiEncryption string
	wifiHidden     bool
)

var wifiCmd = &cobra.Command{
	Use:   "wifi",
	Short: "Generate QR code for WiFi network",
	Long: `Generate a QR code that allows devices to connect to a WiFi network.

Supported encryption types:
  WPA   - WPA/WPA2/WPA3 (default when password is provided)
  WEP   - WEP encryption
  nopass - Open network (no password)

Examples:
  mkqr wifi -s "MyNetwork" -p "password123"
  mkqr wifi --ssid "Home WiFi" --password "secret" --encryption WPA
  mkqr wifi -s "Guest" --encryption nopass
  mkqr wifi -s "Hidden Network" -p "pass" --hidden`,
	RunE: runWifi,
}

func init() {
	wifiCmd.Flags().StringVarP(&wifiSSID, "ssid", "s", "", "Network name (SSID) [required]")
	wifiCmd.Flags().StringVarP(&wifiPassword, "password", "p", "", "Network password")
	wifiCmd.Flags().StringVarP(&wifiEncryption, "encryption", "e", "", "Encryption type (WPA/WEP/nopass)")
	wifiCmd.Flags().BoolVarP(&wifiHidden, "hidden", "H", false, "Hidden network")

	wifiCmd.MarkFlagRequired("ssid")

	rootCmd.AddCommand(wifiCmd)
}

func runWifi(cmd *cobra.Command, args []string) error {
	encryption := encoder.WPA
	if wifiEncryption != "" {
		var err error
		encryption, err = encoder.ParseWiFiEncryption(wifiEncryption)
		if err != nil {
			return err
		}
	} else if wifiPassword == "" {
		encryption = encoder.NoPass
	}

	wifi := &encoder.WiFi{
		SSID:       wifiSSID,
		Password:   wifiPassword,
		Encryption: encryption,
		Hidden:     wifiHidden,
	}

	content := wifi.Encode()

	if !quiet {
		fmt.Fprintf(cmd.ErrOrStderr(), "WiFi: %s (%s)\n", wifiSSID, encryption)
	}

	return generateQR(content)
}
