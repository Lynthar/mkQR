package cli

import (
	"fmt"

	"github.com/Lynthar/mkQR/internal/encoder"
	"github.com/spf13/cobra"
)

var (
	emailTo      string
	emailCC      string
	emailBCC     string
	emailSubject string
	emailBody    string
)

var emailCmd = &cobra.Command{
	Use:   "email <address>",
	Short: "Generate QR code for email",
	Long: `Generate a QR code that opens an email composer when scanned.

Examples:
  mkqr email hello@example.com
  mkqr email hello@example.com -s "Subject" -b "Message body"
  mkqr email support@example.com --subject "Help Request" --cc "manager@example.com"`,
	Args: cobra.ExactArgs(1),
	RunE: runEmail,
}

func init() {
	emailCmd.Flags().StringVarP(&emailSubject, "subject", "s", "", "Email subject")
	emailCmd.Flags().StringVarP(&emailBody, "body", "b", "", "Email body")
	emailCmd.Flags().StringVar(&emailCC, "cc", "", "CC recipients")
	emailCmd.Flags().StringVar(&emailBCC, "bcc", "", "BCC recipients")

	rootCmd.AddCommand(emailCmd)
}

func runEmail(cmd *cobra.Command, args []string) error {
	email := &encoder.Email{
		To:      args[0],
		CC:      emailCC,
		BCC:     emailBCC,
		Subject: emailSubject,
		Body:    emailBody,
	}

	content := email.Encode()

	if !quiet {
		fmt.Fprintf(cmd.ErrOrStderr(), "Email: %s\n", args[0])
	}

	return generateQR(content)
}
