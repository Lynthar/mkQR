package cli

import (
	"fmt"

	"github.com/Lynthar/mkQR/internal/encoder"
	"github.com/spf13/cobra"
)

var phoneCmd = &cobra.Command{
	Use:   "phone <number>",
	Short: "Generate QR code for a phone number",
	Long: `Generate a QR code that initiates a phone call when scanned.

Examples:
  mkqr phone +1234567890
  mkqr phone "+86 138 1234 5678"
  mkqr phone 4001234567`,
	Args: cobra.ExactArgs(1),
	RunE: runPhone,
}

func init() {
	rootCmd.AddCommand(phoneCmd)
}

func runPhone(cmd *cobra.Command, args []string) error {
	phone := &encoder.Phone{
		Number: args[0],
	}

	content := phone.Encode()

	if !quiet {
		fmt.Fprintf(cmd.ErrOrStderr(), "Phone: %s\n", args[0])
	}

	return generateQR(content)
}
