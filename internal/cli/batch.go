package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Lynthar/mkQR/internal/encoder"
	"github.com/Lynthar/mkQR/internal/qr"
	"github.com/spf13/cobra"
)

var (
	batchOutputDir string
	batchPrefix    string
	batchFormat    string
)

var batchCmd = &cobra.Command{
	Use:   "batch <file>",
	Short: "Generate QR codes from a file (one per line)",
	Long: `Generate multiple QR codes from a file containing one item per line.

Each line in the input file will generate a separate QR code.
Empty lines and lines starting with # are skipped.

Examples:
  mkqr batch urls.txt -O ./qrcodes/
  mkqr batch nodes.txt --output-dir ./out --prefix "node_"
  cat links.txt | mkqr batch - -O ./out/`,
	Args: cobra.ExactArgs(1),
	RunE: runBatch,
}

func init() {
	batchCmd.Flags().StringVarP(&batchOutputDir, "output-dir", "O", ".", "Output directory")
	batchCmd.Flags().StringVar(&batchPrefix, "prefix", "qr_", "Filename prefix")
	// Note: currently only PNG is supported, flag hidden to avoid confusion
	batchCmd.Flags().StringVar(&batchFormat, "format", "png", "Output format")
	batchCmd.Flags().MarkHidden("format")

	rootCmd.AddCommand(batchCmd)
}

func runBatch(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Validate size
	if outputSize <= 0 {
		return fmt.Errorf("size must be a positive number, got %d", outputSize)
	}

	// Create output directory
	if err := os.MkdirAll(batchOutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Open input file (or stdin if "-")
	var scanner *bufio.Scanner
	if inputFile == "-" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err := os.Open(inputFile)
		if err != nil {
			return fmt.Errorf("failed to open input file: %w", err)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	// Parse error correction level
	level, err := qr.ParseLevel(errorLevel)
	if err != nil {
		return err
	}

	opts := qr.Options{
		Level: level,
		Size:  outputSize,
	}
	gen := qr.NewGenerator(opts)

	count := 0
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Detect content type for logging
		contentType, _ := encoder.DetectAndDescribe(line)

		// Add https:// for URLs without protocol
		content := line
		if contentType == encoder.TypeURL && !strings.HasPrefix(strings.ToLower(line), "http") {
			content = "https://" + line
		}

		// Generate QR code
		qrCode, err := gen.Generate(content)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error on line %d: %v\n", lineNum, err)
			continue
		}

		// Save to file
		filename := filepath.Join(batchOutputDir, fmt.Sprintf("%s%04d.%s", batchPrefix, count+1, batchFormat))
		if err := qr.SavePNG(qrCode, filename, outputSize); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error saving line %d: %v\n", lineNum, err)
			continue
		}

		if !quiet {
			// Truncate long content for display
			preview := line
			if len(preview) > 40 {
				preview = preview[:40] + "..."
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "[%d] %s -> %s\n", count+1, preview, filename)
		}

		count++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	if !quiet {
		fmt.Fprintf(cmd.ErrOrStderr(), "\nGenerated %d QR codes in %s\n", count, batchOutputDir)
	}

	return nil
}
