package qr

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"  // register GIF decoder
	_ "image/jpeg" // register JPEG decoder
	"image/png"
	"os"
	"path/filepath"

	"github.com/skip2/go-qrcode"
)

// LogoOptions configures logo embedding into a QR code.
type LogoOptions struct {
	// LogoPath is the filesystem path to the logo image (PNG, JPEG, or GIF).
	LogoPath string
	// SizeRatio is logo edge length as a fraction of QR size (default 0.20).
	// Values above 0.35 are rejected: the QR typically becomes unscannable.
	SizeRatio float64
	// PadRatio is the white padding around the logo as a fraction of logo size (default 0.10).
	PadRatio float64
}

// DefaultLogoOptions returns sensible defaults for a given logo path.
func DefaultLogoOptions(logoPath string) LogoOptions {
	return LogoOptions{
		LogoPath:  logoPath,
		SizeRatio: 0.20,
		PadRatio:  0.10,
	}
}

// SavePNGWithLogo writes the QR code as a PNG with a logo composited at center.
// The caller should generate the QR code with error correction level H so the
// masked modules can still be recovered during scanning.
func SavePNGWithLogo(qr *qrcode.QRCode, filename string, size int, opts LogoOptions) error {
	if opts.LogoPath == "" {
		return fmt.Errorf("logo path is required")
	}
	if opts.SizeRatio <= 0 {
		opts.SizeRatio = 0.20
	}
	if opts.SizeRatio > 0.35 {
		return fmt.Errorf("logo size ratio %.2f exceeds safe maximum 0.35 (would make QR unscannable)", opts.SizeRatio)
	}
	if opts.PadRatio < 0 {
		opts.PadRatio = 0
	}

	logo, err := loadImage(opts.LogoPath)
	if err != nil {
		return err
	}

	qrImg := qr.Image(size)
	logoSize := int(float64(size) * opts.SizeRatio)
	if logoSize < 8 {
		logoSize = 8
	}
	padSize := logoSize + int(float64(logoSize)*opts.PadRatio*2)

	out := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(out, out.Bounds(), qrImg, image.Point{}, draw.Src)

	padRect := centeredRect(size, padSize)
	draw.Draw(out, padRect, &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	resized := boxResize(logo, logoSize, logoSize)
	logoRect := centeredRect(size, logoSize)
	draw.Draw(out, logoRect, resized, image.Point{}, draw.Over)

	dir := filepath.Dir(filename)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()
	if err := png.Encode(f, out); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}
	return nil
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open logo: %w", err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode logo (supported: PNG, JPEG, GIF): %w", err)
	}
	return img, nil
}

func centeredRect(canvas, inner int) image.Rectangle {
	offset := (canvas - inner) / 2
	return image.Rect(offset, offset, offset+inner, offset+inner)
}

// boxResize downsizes src to dstW x dstH using area averaging (box filter).
// Quality is good for any downscale; upscaling degrades to nearest-neighbor sampling.
// Zero external dependencies.
func boxResize(src image.Image, dstW, dstH int) *image.RGBA {
	sb := src.Bounds()
	sw, sh := sb.Dx(), sb.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	if sw == 0 || sh == 0 {
		return dst
	}
	for dy := 0; dy < dstH; dy++ {
		y0 := sb.Min.Y + dy*sh/dstH
		y1 := sb.Min.Y + (dy+1)*sh/dstH
		if y1 <= y0 {
			y1 = y0 + 1
		}
		for dx := 0; dx < dstW; dx++ {
			x0 := sb.Min.X + dx*sw/dstW
			x1 := sb.Min.X + (dx+1)*sw/dstW
			if x1 <= x0 {
				x1 = x0 + 1
			}
			var r, g, b, a uint64
			var n uint64
			for sy := y0; sy < y1; sy++ {
				for sx := x0; sx < x1; sx++ {
					cr, cg, cb, ca := src.At(sx, sy).RGBA()
					r += uint64(cr)
					g += uint64(cg)
					b += uint64(cb)
					a += uint64(ca)
					n++
				}
			}
			if n == 0 {
				n = 1
			}
			dst.SetRGBA(dx, dy, color.RGBA{
				R: uint8((r / n) >> 8),
				G: uint8((g / n) >> 8),
				B: uint8((b / n) >> 8),
				A: uint8((a / n) >> 8),
			})
		}
	}
	return dst
}
