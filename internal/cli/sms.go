package cli

import (
	"fmt"

	"github.com/Lynthar/mkQR/internal/encoder"
	"github.com/spf13/cobra"
)

var smsBody string

var smsCmd = &cobra.Command{
	Use:   "sms <number>",
	Short: "Generate QR code for SMS",
	Long: `Generate a QR code that opens the SMS app when scanned.

Examples:
  mkqr sms +1234567890
  mkqr sms +1234567890 -b "Hello!"
  mkqr sms 10086 --body "CXYE"`,
	Args: cobra.ExactArgs(1),
	RunE: runSMS,
}

func init() {
	smsCmd.Flags().StringVarP(&smsBody, "body", "b", "", "Message body")

	rootCmd.AddCommand(smsCmd)
}

func runSMS(cmd *cobra.Command, args []string) error {
	sms := &encoder.SMS{
		Number: args[0],
		Body:   smsBody,
	}

	content := sms.Encode()

	if !quiet {
		fmt.Fprintf(cmd.ErrOrStderr(), "SMS: %s\n", args[0])
	}

	return generateQR(content)
}
