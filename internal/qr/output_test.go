package qr

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		filename string
		expected OutputFormat
	}{
		{"output.png", FormatPNG},
		{"output.PNG", FormatPNG},
		{"path/to/file.png", FormatPNG},
		{"output.jpg", FormatPNG},  // defaults to PNG
		{"output", FormatPNG},       // defaults to PNG
		{"output.jpeg", FormatPNG},  // defaults to PNG
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := DetectFormat(tt.filename)
			if result != tt.expected {
				t.Errorf("DetectFormat(%q) = %q, want %q", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestSavePNG(t *testing.T) {
	gen := NewGenerator(DefaultOptions())
	qr, err := gen.Generate("Test content")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "mkqr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test saving to file
	filename := filepath.Join(tmpDir, "test.png")
	err = SavePNG(qr, filename, 256)
	if err != nil {
		t.Fatalf("SavePNG() error: %v", err)
	}

	// Verify file exists
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("File not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("File is empty")
	}

	// Test saving to nested directory
	nestedFile := filepath.Join(tmpDir, "nested", "dir", "test.png")
	err = SavePNG(qr, nestedFile, 256)
	if err != nil {
		t.Fatalf("SavePNG() with nested dir error: %v", err)
	}

	// Verify nested file exists
	_, err = os.Stat(nestedFile)
	if err != nil {
		t.Fatalf("Nested file not created: %v", err)
	}
}

func TestToBase64(t *testing.T) {
	gen := NewGenerator(DefaultOptions())
	qr, err := gen.Generate("Test content")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	base64Str, err := ToBase64(qr, 256)
	if err != nil {
		t.Fatalf("ToBase64() error: %v", err)
	}

	if len(base64Str) == 0 {
		t.Error("ToBase64() returned empty string")
	}

	// Base64 should only contain valid characters
	validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	for _, c := range base64Str {
		found := false
		for _, vc := range validChars {
			if c == vc {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ToBase64() contains invalid character: %c", c)
		}
	}
}
