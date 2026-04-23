package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var textCmd = &cobra.Command{
	Use:   "text <content>",
	Short: "Generate QR code for plain text",
	Long: `Generate a QR code for plain text content.

Use this when you want to ensure the content is treated as plain text,
without any automatic URL or format detection.

Examples:
  mkqr text "Hello World"
  mkqr text "This is a multi-line\ntext message"
  mkqr text "Some text that looks like a URL: example.com"`,
	Args: cobra.ExactArgs(1),
	RunE: runText,
}

func init() {
	rootCmd.AddCommand(textCmd)
}

func runText(cmd *cobra.Command, args []string) error {
	content := args[0]

	if !quiet {
		// Truncate by runes, not bytes, so multi-byte characters (e.g. CJK)
		// aren't sliced mid-character into garbled output.
		preview := content
		if runes := []rune(preview); len(runes) > 50 {
			preview = string(runes[:50]) + "..."
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "Text: %s\n", preview)
	}

	return generateQR(content)
}
