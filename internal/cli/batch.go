package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Lynthar/mkQR/internal/qr"
	"github.com/Lynthar/mkQR/pkg/encoder"
	"github.com/spf13/cobra"
)

var (
	batchOutputDir string
	batchPrefix    string
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
	// Raise the per-line cap from the default 64KB — a single vmess:// or
	// subscription line can easily exceed that when it carries a large
	// base64 payload, and hitting the cap aborts the whole batch.
	scanner.Buffer(make([]byte, 64*1024), 10*1024*1024)

	gen, err := buildGenerator(cmd.ErrOrStderr())
	if err != nil {
		return err
	}

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
		if contentType == encoder.TypeURL {
			content = ensureHTTPScheme(line)
		}

		// Generate QR code
		qrCode, err := gen.Generate(content)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error on line %d: %v\n", lineNum, err)
			continue
		}

		// Save to file (PNG only — SVG batch output isn't wired up yet).
		filename := filepath.Join(batchOutputDir, fmt.Sprintf("%s%04d.png", batchPrefix, count+1))
		if logoPath != "" {
			if err := qr.SavePNGWithLogo(qrCode, filename, outputSize, qr.DefaultLogoOptions(logoPath)); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error saving line %d: %v\n", lineNum, err)
				continue
			}
		} else {
			if err := qr.SavePNG(qrCode, filename, outputSize); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error saving line %d: %v\n", lineNum, err)
				continue
			}
		}

		if !quiet {
			// Truncate by runes, not bytes, so multi-byte characters (e.g.
			// CJK) aren't sliced mid-character into garbled output.
			preview := line
			if runes := []rune(preview); len(runes) > 40 {
				preview = string(runes[:40]) + "..."
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
