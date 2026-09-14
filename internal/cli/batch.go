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
	batchCmd.Flags().StringVar(&batchPrefix, "prefix", "qr_", "Filename prefix (avoid path separators or '..'; not sanitized)")

	rootCmd.AddCommand(batchCmd)
}

func runBatch(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	gen, err := buildGenerator(cmd.ErrOrStderr())
	if err != nil {
		return err
	}

	// -o names a single file; batch names one file per line under --output-dir.
	if outputFile != "" && !quiet {
		cmd.PrintErrln("Note: batch writes PNG files into --output-dir; -o is ignored")
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
		defer func() { _ = file.Close() }()
		scanner = bufio.NewScanner(file)
	}
	// Raise the per-line cap from the default 64KB — a single vmess:// or
	// subscription line can easily exceed that when it carries a large
	// base64 payload, and hitting the cap aborts the whole batch.
	scanner.Buffer(make([]byte, 64*1024), 10*1024*1024)

	// --prefix may carry path components, so files can land outside
	// --output-dir; the summary must name where they actually went.
	nameFor := func(line int) string {
		return filepath.Join(batchOutputDir, fmt.Sprintf("%s%04d.png", batchPrefix, line))
	}
	writtenDir := filepath.Dir(nameFor(1))

	count := 0
	failed := 0
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
			cmd.PrintErrf("Error on line %d: %v\n", lineNum, err)
			failed++
			continue
		}

		noteSizeRaised(cmd.ErrOrStderr(), qrCode)

		// PNG only (SVG batch output isn't wired up). The number is the input
		// line, so every name points back at the line that produced it.
		filename := nameFor(lineNum)
		if logoPath != "" {
			if err := qr.SavePNGWithLogo(qrCode, filename, outputSize, qr.DefaultLogoOptions(logoPath)); err != nil {
				cmd.PrintErrf("Error saving line %d: %v\n", lineNum, err)
				failed++
				continue
			}
		} else {
			if err := qr.SavePNG(qrCode, filename, outputSize); err != nil {
				cmd.PrintErrf("Error saving line %d: %v\n", lineNum, err)
				failed++
				continue
			}
		}

		if !quiet {
			cmd.PrintErrf("[%d] %s -> %s\n", lineNum, previewOf(line, 40), filename)
		}

		count++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	if !quiet {
		cmd.PrintErrf("\nGenerated %d QR codes in %s\n", count, writtenDir)
	}

	if failed > 0 {
		return fmt.Errorf("%d of %d lines failed", failed, count+failed)
	}

	return nil
}
