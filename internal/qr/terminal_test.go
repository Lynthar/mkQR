package qr

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderTerminal(t *testing.T) {
	gen := NewGenerator(DefaultOptions())
	qr, err := gen.Generate("Test")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	tests := []struct {
		name   string
		config TerminalConfig
	}{
		{"normal", TerminalConfig{Invert: false, Small: false}},
		{"inverted", TerminalConfig{Invert: true, Small: false}},
		{"small", TerminalConfig{Invert: false, Small: true}},
		{"small inverted", TerminalConfig{Invert: true, Small: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			RenderTerminal(&buf, qr, tt.config)

			output := buf.String()
			if len(output) == 0 {
				t.Error("RenderTerminal() produced empty output")
			}

			// Should have multiple lines
			lines := strings.Split(output, "\n")
			if len(lines) < 5 {
				t.Errorf("RenderTerminal() produced too few lines: %d", len(lines))
			}
		})
	}
}

func TestRenderNormal(t *testing.T) {
	// Create a simple 2x2 bitmap for testing
	bitmap := [][]bool{
		{true, false},
		{false, true},
	}

	var buf bytes.Buffer
	renderNormal(&buf, bitmap, false)
	output := buf.String()

	// Should contain block characters
	if !strings.Contains(output, "██") && !strings.Contains(output, "  ") {
		t.Error("renderNormal() output doesn't contain expected block characters")
	}

	// Test inverted
	buf.Reset()
	renderNormal(&buf, bitmap, true)
	outputInverted := buf.String()

	// Inverted should be different
	if output == outputInverted {
		t.Error("renderNormal() inverted output is same as normal")
	}
}

func TestRenderSmall(t *testing.T) {
	// Create a simple 4x2 bitmap for testing
	bitmap := [][]bool{
		{true, false},
		{false, true},
		{true, true},
		{false, false},
	}

	var buf bytes.Buffer
	renderSmall(&buf, bitmap, false)
	output := buf.String()

	if len(output) == 0 {
		t.Error("renderSmall() produced empty output")
	}

	// Small mode should use half-block characters
	hasHalfBlock := strings.Contains(output, "▀") ||
		strings.Contains(output, "▄") ||
		strings.Contains(output, "█") ||
		strings.Contains(output, " ")
	if !hasHalfBlock {
		t.Error("renderSmall() output doesn't contain expected half-block characters")
	}

	// Test inverted
	buf.Reset()
	renderSmall(&buf, bitmap, true)
	outputInverted := buf.String()

	// Inverted should be different
	if output == outputInverted {
		t.Error("renderSmall() inverted output is same as normal")
	}
}

func TestRenderSmallOddRows(t *testing.T) {
	// Test with odd number of rows
	bitmap := [][]bool{
		{true, false},
		{false, true},
		{true, false},
	}

	var buf bytes.Buffer
	// Should not panic with odd rows
	renderSmall(&buf, bitmap, false)

	if len(buf.String()) == 0 {
		t.Error("renderSmall() with odd rows produced empty output")
	}
}
