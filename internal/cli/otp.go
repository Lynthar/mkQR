package cli

import (
	"fmt"

	"github.com/Lynthar/mkQR/internal/encoder"
	"github.com/spf13/cobra"
)

var (
	otpSecret    string
	otpIssuer    string
	otpAccount   string
	otpAlgorithm string
	otpDigits    int
	otpPeriod    int
	otpCounter   int
	otpTypeHOTP  bool
)

var otpCmd = &cobra.Command{
	Use:   "otp",
	Short: "Generate QR code for OTP (2FA authenticator)",
	Long: `Generate a QR code for setting up two-factor authentication.

This creates an otpauth:// URL compatible with Google Authenticator,
Authy, Microsoft Authenticator, and other TOTP/HOTP apps.

Examples:
  mkqr otp -s "JBSWY3DPEHPK3PXP" -i "GitHub" -a "user@example.com"
  mkqr otp --secret "ABCD1234" --issuer "AWS" --account "myaccount"
  mkqr otp -s "SECRET" -i "Service" -a "user" --digits 8 --period 60`,
	RunE: runOTP,
}

func init() {
	otpCmd.Flags().StringVarP(&otpSecret, "secret", "s", "", "Secret key (base32 encoded) [required]")
	otpCmd.Flags().StringVarP(&otpIssuer, "issuer", "i", "", "Service/issuer name [required]")
	otpCmd.Flags().StringVarP(&otpAccount, "account", "a", "", "Account name/email [required]")
	otpCmd.Flags().StringVar(&otpAlgorithm, "algorithm", "SHA1", "Hash algorithm (SHA1/SHA256/SHA512)")
	otpCmd.Flags().IntVar(&otpDigits, "digits", 6, "Number of digits (6 or 8)")
	otpCmd.Flags().IntVar(&otpPeriod, "period", 30, "Time period in seconds (TOTP)")
	otpCmd.Flags().IntVar(&otpCounter, "counter", 0, "Initial counter value (HOTP)")
	otpCmd.Flags().BoolVar(&otpTypeHOTP, "hotp", false, "Use HOTP (counter-based) instead of TOTP")

	otpCmd.MarkFlagRequired("secret")
	otpCmd.MarkFlagRequired("issuer")
	otpCmd.MarkFlagRequired("account")

	rootCmd.AddCommand(otpCmd)
}

func runOTP(cmd *cobra.Command, args []string) error {
	// Validate secret is valid base32
	if err := encoder.ValidateSecret(otpSecret); err != nil {
		return fmt.Errorf("invalid secret: %w", err)
	}

	// Validate digits
	if otpDigits != 6 && otpDigits != 8 {
		return fmt.Errorf("digits must be 6 or 8, got %d", otpDigits)
	}

	// Validate period
	if otpPeriod <= 0 {
		return fmt.Errorf("period must be a positive number, got %d", otpPeriod)
	}

	// Validate algorithm
	validAlgorithms := map[string]bool{"SHA1": true, "SHA256": true, "SHA512": true}
	if !validAlgorithms[otpAlgorithm] {
		return fmt.Errorf("algorithm must be SHA1, SHA256, or SHA512, got %s", otpAlgorithm)
	}

	otpType := encoder.TOTP
	if otpTypeHOTP {
		otpType = encoder.HOTP
	}

	otp := &encoder.OTP{
		Type:      otpType,
		Secret:    otpSecret,
		Issuer:    otpIssuer,
		Account:   otpAccount,
		Algorithm: otpAlgorithm,
		Digits:    otpDigits,
		Period:    otpPeriod,
		Counter:   otpCounter,
	}

	content := otp.Encode()

	if !quiet {
		fmt.Fprintf(cmd.ErrOrStderr(), "OTP: %s (%s)\n", otpIssuer, otpAccount)
	}

	return generateQR(content)
}
