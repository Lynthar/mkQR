package qr

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/skip2/go-qrcode"
)

// OutputFormat represents the output format
type OutputFormat string

const (
	FormatTerminal OutputFormat = "terminal"
	FormatPNG      OutputFormat = "png"
	FormatBase64   OutputFormat = "base64"
)

// DetectFormat detects output format from filename
func DetectFormat(filename string) OutputFormat {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".png":
		return FormatPNG
	default:
		return FormatPNG // Default to PNG for files
	}
}

// SavePNG saves the QR code as a PNG file
func SavePNG(qr *qrcode.QRCode, filename string, size int) error {
	// Ensure directory exists
	dir := filepath.Dir(filename)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	if err := qr.WriteFile(size, filename); err != nil {
		return fmt.Errorf("failed to write PNG file: %w", err)
	}
	return nil
}

// ToBase64 returns the QR code as a base64-encoded PNG string
func ToBase64(qr *qrcode.QRCode, size int) (string, error) {
	png, err := qr.PNG(size)
	if err != nil {
		return "", fmt.Errorf("failed to generate PNG: %w", err)
	}
	return base64.StdEncoding.EncodeToString(png), nil
}
