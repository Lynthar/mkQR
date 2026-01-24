package qr

import (
	"testing"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected ErrorCorrectionLevel
		hasError bool
	}{
		{"L", LevelL, false},
		{"l", LevelL, false},
		{"M", LevelM, false},
		{"m", LevelM, false},
		{"Q", LevelQ, false},
		{"q", LevelQ, false},
		{"H", LevelH, false},
		{"h", LevelH, false},
		{"X", LevelM, true},
		{"invalid", LevelM, true},
		{"", LevelM, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseLevel(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("ParseLevel(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseLevel(%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	if opts.Level != LevelM {
		t.Errorf("DefaultOptions().Level = %v, want %v", opts.Level, LevelM)
	}
	if opts.Size != 256 {
		t.Errorf("DefaultOptions().Size = %v, want 256", opts.Size)
	}
	if opts.ForegroundColor == nil {
		t.Error("DefaultOptions().ForegroundColor is nil")
	}
	if opts.BackgroundColor == nil {
		t.Error("DefaultOptions().BackgroundColor is nil")
	}
}

func TestGeneratorGenerate(t *testing.T) {
	gen := NewGenerator(DefaultOptions())

	// Test basic generation
	qr, err := gen.Generate("Hello World")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	if qr == nil {
		t.Fatal("Generate() returned nil QR code")
	}

	// Test empty content (should fail)
	_, err = gen.Generate("")
	if err == nil {
		t.Error("Generate() with empty content should fail")
	}
}

func TestGeneratorGeneratePNG(t *testing.T) {
	gen := NewGenerator(DefaultOptions())

	png, err := gen.GeneratePNG("Test content")
	if err != nil {
		t.Fatalf("GeneratePNG() error: %v", err)
	}
	if len(png) == 0 {
		t.Fatal("GeneratePNG() returned empty bytes")
	}

	// Check PNG magic bytes
	if len(png) < 8 {
		t.Fatal("GeneratePNG() returned too short data")
	}
	// PNG magic: 137 80 78 71 13 10 26 10
	pngMagic := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	for i, b := range pngMagic {
		if png[i] != b {
			t.Errorf("GeneratePNG() invalid PNG header at byte %d: got %x, want %x", i, png[i], b)
		}
	}
}

func TestErrorCorrectionLevelToQRCode(t *testing.T) {
	tests := []struct {
		level    ErrorCorrectionLevel
		expected int // qrcode.RecoveryLevel values
	}{
		{LevelL, 0}, // qrcode.Low
		{LevelM, 1}, // qrcode.Medium
		{LevelQ, 2}, // qrcode.High
		{LevelH, 3}, // qrcode.Highest
	}

	for _, tt := range tests {
		result := tt.level.toQRCode()
		if int(result) != tt.expected {
			t.Errorf("ErrorCorrectionLevel(%v).toQRCode() = %v, want %v", tt.level, result, tt.expected)
		}
	}
}
