package encoder_test

import (
	"fmt"

	"github.com/Lynthar/mkQR/pkg/encoder"
)

func ExampleWiFi_Encode() {
	s := (&encoder.WiFi{SSID: "Home", Password: "secret", Encryption: encoder.WPA}).Encode()
	fmt.Println(s)
	// Output: WIFI:T:WPA;S:Home;P:secret;;
}

// An unset Encryption is inferred from the password: nopass without one, WPA
// with one. ParseWiFiEncryption("") returns that same unset value.
func ExampleWiFi_Encode_inferredEncryption() {
	fmt.Println((&encoder.WiFi{SSID: "Cafe"}).Encode())
	fmt.Println((&encoder.WiFi{SSID: "Home", Password: "secret"}).Encode())
	// Output:
	// WIFI:T:nopass;S:Cafe;P:;;
	// WIFI:T:WPA;S:Home;P:secret;;
}

// Encode never checks the secret, so a library caller has to.
func ExampleOTP_Encode() {
	otp := &encoder.OTP{
		Type:    encoder.TOTP,
		Secret:  "JBSWY3DPEHPK3PXP",
		Issuer:  "Example",
		Account: "alice@example.com",
	}
	if err := encoder.ValidateSecret(otp.Secret); err != nil {
		fmt.Println("invalid secret:", err)
		return
	}
	fmt.Println(otp.Encode())
	// Output: otpauth://totp/Example:alice%40example.com?secret=JBSWY3DPEHPK3PXP&issuer=Example
}

func ExampleValidateSecret() {
	fmt.Println(encoder.ValidateSecret("JBSW Y3DP EHPK 3PXP"))
	fmt.Println(encoder.ValidateSecret("not base32!"))
	// Output:
	// <nil>
	// secret must be a valid base32 string (A-Z, 2-7)
}

func ExampleDetectAndDescribe() {
	for _, in := range []string{"https://example.com", "vmess://eyJ2IjogIjIifQ", "hello"} {
		contentType, description := encoder.DetectAndDescribe(in)
		fmt.Printf("%s -> %s\n", contentType, description)
	}
	// Output:
	// url -> URL
	// proxy -> Proxy configuration
	// text -> Plain text
}
