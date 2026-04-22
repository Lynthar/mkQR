package qr

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestSavePNGWithLogo(t *testing.T) {
	gen := NewGenerator(Options{Level: LevelH, Size: 256})
	qr, err := gen.Generate("Test content for logo")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "mkqr-logo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a 32x32 solid red logo
	logoPath := filepath.Join(tmpDir, "logo.png")
	logoImg := image.NewRGBA(image.Rect(0, 0, 32, 32))
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			logoImg.SetRGBA(x, y, red)
		}
	}
	logoFile, err := os.Create(logoPath)
	if err != nil {
		t.Fatalf("Create logo: %v", err)
	}
	if err := png.Encode(logoFile, logoImg); err != nil {
		t.Fatalf("Encode logo: %v", err)
	}
	logoFile.Close()

	outPath := filepath.Join(tmpDir, "out.png")
	if err := SavePNGWithLogo(qr, outPath, 256, DefaultLogoOptions(logoPath)); err != nil {
		t.Fatalf("SavePNGWithLogo: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("Read output: %v", err)
	}
	if len(data) < 8 {
		t.Fatal("Output too small to be a PNG")
	}
	pngMagic := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	if !bytes.Equal(data[:8], pngMagic) {
		t.Error("Output is not a valid PNG")
	}

	decoded, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Decode output: %v", err)
	}

	// Center pixel should be dominated by red (the logo)
	center := decoded.At(128, 128)
	r, g, b, _ := center.RGBA()
	if r>>8 < 150 || g>>8 > 80 || b>>8 > 80 {
		t.Errorf("Center pixel not red-dominated: rgb=(%d,%d,%d)", r>>8, g>>8, b>>8)
	}
}

func TestSavePNGWithLogoRejectsOversized(t *testing.T) {
	gen := NewGenerator(Options{Level: LevelH, Size: 256})
	qr, err := gen.Generate("x")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	err = SavePNGWithLogo(qr, filepath.Join(os.TempDir(), "unused.png"), 256, LogoOptions{
		LogoPath:  "whatever.png",
		SizeRatio: 0.5,
	})
	if err == nil {
		t.Error("Expected error for size ratio > 0.35")
	}
}

func TestSavePNGWithLogoRejectsEmptyPath(t *testing.T) {
	gen := NewGenerator(DefaultOptions())
	qr, _ := gen.Generate("x")
	err := SavePNGWithLogo(qr, "out.png", 256, LogoOptions{})
	if err == nil {
		t.Error("Expected error for empty logo path")
	}
}

func TestSavePNGWithLogoInvalidFile(t *testing.T) {
	gen := NewGenerator(DefaultOptions())
	qr, _ := gen.Generate("x")
	err := SavePNGWithLogo(qr, "out.png", 256, DefaultLogoOptions("/nonexistent/path/logo.png"))
	if err == nil {
		t.Error("Expected error for non-existent logo file")
	}
}

func TestBoxResizeDownscale(t *testing.T) {
	// 100x100 gradient, resize to 10x10 — corners should differ
	src := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			src.SetRGBA(x, y, color.RGBA{R: uint8(x * 255 / 99), G: uint8(y * 255 / 99), B: 0, A: 255})
		}
	}
	dst := boxResize(src, 10, 10)
	if dst.Bounds().Dx() != 10 || dst.Bounds().Dy() != 10 {
		t.Fatalf("Wrong dst size: %v", dst.Bounds())
	}
	topLeft := dst.RGBAAt(0, 0)
	bottomRight := dst.RGBAAt(9, 9)
	if topLeft.R >= bottomRight.R {
		t.Errorf("Gradient not preserved: topLeft R=%d, bottomRight R=%d", topLeft.R, bottomRight.R)
	}
	if topLeft.G >= bottomRight.G {
		t.Errorf("Gradient not preserved: topLeft G=%d, bottomRight G=%d", topLeft.G, bottomRight.G)
	}
}

func TestBoxResizeEmpty(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 0, 0))
	dst := boxResize(src, 5, 5)
	if dst.Bounds().Dx() != 5 || dst.Bounds().Dy() != 5 {
		t.Errorf("Empty src should still produce sized dst")
	}
}
