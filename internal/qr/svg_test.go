package qr

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestToSVG(t *testing.T) {
	gen := NewGenerator(DefaultOptions())
	qr, err := gen.Generate("Test content")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	svg := string(ToSVG(qr, 256))

	if !strings.HasPrefix(svg, `<?xml`) {
		t.Error("SVG output missing XML declaration")
	}
	if !strings.Contains(svg, `<svg`) {
		t.Error("SVG output missing <svg> element")
	}
	if !strings.Contains(svg, `xmlns="http://www.w3.org/2000/svg"`) {
		t.Error("SVG output missing xmlns attribute")
	}
	if !strings.Contains(svg, `width="256"`) {
		t.Errorf("SVG output missing width=\"256\": %s", svg[:200])
	}
	if !strings.Contains(svg, `<path`) {
		t.Error("SVG output missing <path> element")
	}
	if !strings.Contains(svg, `</svg>`) {
		t.Error("SVG output missing closing </svg> tag")
	}
	if !strings.Contains(svg, `shape-rendering="crispEdges"`) {
		t.Error("SVG output missing crispEdges for scanner compatibility")
	}
}

func TestSaveSVG(t *testing.T) {
	gen := NewGenerator(DefaultOptions())
	qr, err := gen.Generate("Test content")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "mkqr-svg-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Basic save
	filename := filepath.Join(tmpDir, "test.svg")
	if err := SaveSVG(qr, filename, 256); err != nil {
		t.Fatalf("SaveSVG() error: %v", err)
	}
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("File not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("SVG file is empty")
	}

	// Save into a nested directory that does not yet exist
	nested := filepath.Join(tmpDir, "nested", "dir", "test.svg")
	if err := SaveSVG(qr, nested, 256); err != nil {
		t.Fatalf("SaveSVG() with nested dir error: %v", err)
	}
	if _, err := os.Stat(nested); err != nil {
		t.Fatalf("Nested SVG not created: %v", err)
	}
}

func TestDetectFormatSVG(t *testing.T) {
	cases := map[string]OutputFormat{
		"out.svg":          FormatSVG,
		"out.SVG":          FormatSVG,
		"path/to/file.svg": FormatSVG,
		"out.png":          FormatPNG,
		"out":              FormatPNG,
	}
	for filename, want := range cases {
		if got := DetectFormat(filename); got != want {
			t.Errorf("DetectFormat(%q) = %q, want %q", filename, got, want)
		}
	}
}
