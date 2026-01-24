package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Lynthar/mkQR/internal/encoder"
	"github.com/Lynthar/mkQR/internal/qr"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	outputFile  string
	outputSize  int
	errorLevel  string
	quiet       bool
	invert      bool
	small       bool
	showVersion bool

	// Version info (set at build time)
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "mkqr [content]",
	Short: "Generate QR codes from various data types",
	Long: `mkQR - A fast, flexible QR code generator for the command line.

Supports WiFi, URLs, contacts, OTP, and many more formats.
Can output to terminal, PNG, or base64.

Examples:
  mkqr "Hello World"                    # Generate QR for text
  mkqr "https://github.com"             # Auto-detect URL
  mkqr wifi -s "MyNetwork" -p "pass"    # WiFi network
  mkqr "vmess://..." -o proxy.png       # Save proxy QR to file
  echo "text" | mkqr                    # Read from stdin`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRoot,
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "Output file (PNG format)")
	rootCmd.PersistentFlags().IntVar(&outputSize, "size", 256, "QR code size in pixels")
	rootCmd.PersistentFlags().StringVarP(&errorLevel, "level", "l", "M", "Error correction level (L/M/Q/H)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Suppress non-essential output")
	rootCmd.PersistentFlags().BoolVar(&invert, "invert", false, "Invert colors (for dark terminals)")
	rootCmd.PersistentFlags().BoolVar(&small, "small", false, "Use compact display mode")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "Show version information")
}

func runRoot(cmd *cobra.Command, args []string) error {
	// Handle version flag
	if showVersion {
		fmt.Printf("mkqr version %s\n", Version)
		fmt.Printf("  commit: %s\n", GitCommit)
		fmt.Printf("  built:  %s\n", BuildDate)
		return nil
	}

	var content string

	if len(args) > 0 {
		content = args[0]
	} else {
		// Check if stdin has data
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// Reading from pipe
			reader := bufio.NewReader(os.Stdin)
			data, err := io.ReadAll(reader)
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			content = strings.TrimSpace(string(data))
		} else {
			return cmd.Help()
		}
	}

	if content == "" {
		return fmt.Errorf("no content provided")
	}

	// Auto-detect content type
	contentType, description := encoder.DetectAndDescribe(content)
	if !quiet {
		fmt.Fprintf(os.Stderr, "Detected: %s\n", description)
	}

	// For detected URLs without protocol, add https://
	if contentType == encoder.TypeURL && !strings.HasPrefix(strings.ToLower(content), "http") {
		content = "https://" + content
	}

	return generateQR(content)
}

// generateQR is the common QR generation logic
func generateQR(content string) error {
	// Validate size
	if outputSize <= 0 {
		return fmt.Errorf("size must be a positive number, got %d", outputSize)
	}

	level, err := qr.ParseLevel(errorLevel)
	if err != nil {
		return err
	}

	opts := qr.Options{
		Level: level,
		Size:  outputSize,
	}

	gen := qr.NewGenerator(opts)
	qrCode, err := gen.Generate(content)
	if err != nil {
		return err
	}

	// Output to file or terminal
	if outputFile != "" {
		if err := qr.SavePNG(qrCode, outputFile, outputSize); err != nil {
			return err
		}
		if !quiet {
			fmt.Fprintf(os.Stderr, "Saved to: %s\n", outputFile)
		}
	} else {
		// Render to terminal
		cfg := qr.TerminalConfig{
			Invert: invert,
			Small:  small,
		}
		qr.RenderTerminal(os.Stdout, qrCode, cfg)
	}

	return nil
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// GetRootCmd returns the root command (for adding subcommands)
func GetRootCmd() *cobra.Command {
	return rootCmd
}
