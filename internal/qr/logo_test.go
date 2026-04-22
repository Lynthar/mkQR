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

func TestSavePNGWithLogoHaloFollowsBackground(t *testing.T) {
	// QR with a distinctive cyan background; halo should match so the logo
	// halo doesn't look like a white sticker on a colored code.
	bg := color.NRGBA{R: 0, G: 200, B: 200, A: 255}
	gen := NewGenerator(Options{Level: LevelH, Size: 256, BackgroundColor: bg})
	qr, err := gen.Generate("halo color regression")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "mkqr-halo-*")
	if err != nil {
		t.Fatalf("mkdir temp: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logoPath := filepath.Join(tmpDir, "logo.png")
	logoImg := image.NewRGBA(image.Rect(0, 0, 20, 20))
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	for y := 0; y < 20; y++ {
		for x := 0; x < 20; x++ {
			logoImg.SetRGBA(x, y, red)
		}
	}
	f, err := os.Create(logoPath)
	if err != nil {
		t.Fatalf("create logo: %v", err)
	}
	if err := png.Encode(f, logoImg); err != nil {
		t.Fatalf("encode logo: %v", err)
	}
	f.Close()

	outPath := filepath.Join(tmpDir, "out.png")
	if err := SavePNGWithLogo(qr, outPath, 256, DefaultLogoOptions(logoPath)); err != nil {
		t.Fatalf("SavePNGWithLogo: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	decoded, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode output: %v", err)
	}

	// logoSize = 0.20 * 256 = 51 → logo rect center ± ~25; padSize = 61 → pad ± ~30.
	// Sample a pixel in the halo ring (inside pad, outside logo): x=155, y=128.
	r, g, b, _ := decoded.At(155, 128).RGBA()
	rb, gb, bb := r>>8, g>>8, b>>8
	// Before the fix the halo was hardcoded white; after, it should match bg (cyan).
	if rb > 80 {
		t.Errorf("halo pixel R=%d too high for cyan bg — is halo still white?", rb)
	}
	if gb < 120 || bb < 120 {
		t.Errorf("halo pixel G=%d B=%d — should be cyan-dominated", gb, bb)
	}
}

func TestBoxResizeEmpty(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 0, 0))
	dst := boxResize(src, 5, 5)
	if dst.Bounds().Dx() != 5 || dst.Bounds().Dy() != 5 {
		t.Errorf("Empty src should still produce sized dst")
	}
}
