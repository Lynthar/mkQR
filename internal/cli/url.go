package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var urlCmd = &cobra.Command{
	Use:   "url <url>",
	Short: "Generate QR code for a URL",
	Long: `Generate a QR code for a URL.

If the URL doesn't start with http:// or https://, https:// will be added automatically.

Examples:
  mkqr url https://github.com
  mkqr url github.com
  mkqr url "https://example.com/path?query=1"`,
	Args: cobra.ExactArgs(1),
	RunE: runURL,
}

func init() {
	rootCmd.AddCommand(urlCmd)
}

func runURL(cmd *cobra.Command, args []string) error {
	url := args[0]

	// Add https:// if no protocol specified
	if !strings.HasPrefix(strings.ToLower(url), "http://") &&
		!strings.HasPrefix(strings.ToLower(url), "https://") {
		url = "https://" + url
	}

	if !quiet {
		fmt.Fprintf(cmd.ErrOrStderr(), "URL: %s\n", url)
	}

	return generateQR(url)
}
