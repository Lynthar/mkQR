package cli

import (
	"fmt"

	"github.com/Lynthar/mkQR/internal/encoder"
	"github.com/spf13/cobra"
)

var (
	vcardFirstName string
	vcardLastName  string
	vcardOrg       string
	vcardTitle     string
	vcardPhone     string
	vcardMobile    string
	vcardEmail     string
	vcardWebsite   string
	vcardAddress   string
	vcardNote      string
)

var vcardCmd = &cobra.Command{
	Use:   "vcard",
	Short: "Generate QR code for a contact (vCard)",
	Long: `Generate a QR code containing contact information in vCard format.

This can be scanned to add a contact to your phone's address book.

Examples:
  mkqr vcard -f "John" --last "Doe" -p "+1234567890"
  mkqr vcard --first "Jane" --last "Smith" --email "jane@example.com" --org "Acme Inc"
  mkqr vcard -f "张" --last "三" -m "+8613812345678" --org "公司名"`,
	RunE: runVCard,
}

func init() {
	vcardCmd.Flags().StringVarP(&vcardFirstName, "first", "f", "", "First name")
	vcardCmd.Flags().StringVar(&vcardLastName, "last", "", "Last name")
	vcardCmd.Flags().StringVarP(&vcardOrg, "org", "O", "", "Organization")
	vcardCmd.Flags().StringVarP(&vcardTitle, "title", "t", "", "Job title")
	vcardCmd.Flags().StringVarP(&vcardPhone, "phone", "p", "", "Phone number")
	vcardCmd.Flags().StringVarP(&vcardMobile, "mobile", "m", "", "Mobile phone")
	vcardCmd.Flags().StringVarP(&vcardEmail, "email", "e", "", "Email address")
	vcardCmd.Flags().StringVarP(&vcardWebsite, "website", "w", "", "Website URL")
	vcardCmd.Flags().StringVarP(&vcardAddress, "address", "a", "", "Address")
	vcardCmd.Flags().StringVarP(&vcardNote, "note", "n", "", "Note")

	rootCmd.AddCommand(vcardCmd)
}

func runVCard(cmd *cobra.Command, args []string) error {
	// Validate at least some info is provided
	if vcardFirstName == "" && vcardLastName == "" && vcardPhone == "" && vcardEmail == "" {
		return fmt.Errorf("at least one of --first, --last, --phone, or --email is required")
	}

	vcard := &encoder.VCard{
		FirstName:    vcardFirstName,
		LastName:     vcardLastName,
		Organization: vcardOrg,
		Title:        vcardTitle,
		Phone:        vcardPhone,
		PhoneMobile:  vcardMobile,
		Email:        vcardEmail,
		Website:      vcardWebsite,
		Address:      vcardAddress,
		Note:         vcardNote,
	}

	content := vcard.Encode()

	if !quiet {
		name := vcardFirstName
		if vcardLastName != "" {
			name += " " + vcardLastName
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "Contact: %s\n", name)
	}

	return generateQR(content)
}
