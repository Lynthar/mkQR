package qr

import (
	"fmt"
	"image/color"

	"github.com/skip2/go-qrcode"
)

// ErrorCorrectionLevel represents QR code error correction level
type ErrorCorrectionLevel int

const (
	LevelL ErrorCorrectionLevel = iota // 7% recovery
	LevelM                             // 15% recovery
	LevelQ                             // 25% recovery
	LevelH                             // 30% recovery
)

// ParseLevel parses a string into ErrorCorrectionLevel
func ParseLevel(s string) (ErrorCorrectionLevel, error) {
	switch s {
	case "L", "l":
		return LevelL, nil
	case "M", "m":
		return LevelM, nil
	case "Q", "q":
		return LevelQ, nil
	case "H", "h":
		return LevelH, nil
	default:
		return LevelM, fmt.Errorf("invalid error correction level: %s (use L, M, Q, or H)", s)
	}
}

func (l ErrorCorrectionLevel) toQRCode() qrcode.RecoveryLevel {
	switch l {
	case LevelL:
		return qrcode.Low
	case LevelM:
		return qrcode.Medium
	case LevelQ:
		return qrcode.High
	case LevelH:
		return qrcode.Highest
	default:
		return qrcode.Medium
	}
}

// Options configures QR code generation
type Options struct {
	Level           ErrorCorrectionLevel
	Size            int
	ForegroundColor color.Color
	BackgroundColor color.Color
}

// DefaultOptions returns default QR generation options
func DefaultOptions() Options {
	return Options{
		Level:           LevelM,
		Size:            256,
		ForegroundColor: color.Black,
		BackgroundColor: color.White,
	}
}

// Generator creates QR codes
type Generator struct {
	opts Options
}

// NewGenerator creates a new QR code generator
func NewGenerator(opts Options) *Generator {
	return &Generator{opts: opts}
}

// Generate creates a QR code from content
func (g *Generator) Generate(content string) (*qrcode.QRCode, error) {
	qr, err := qrcode.New(content, g.opts.Level.toQRCode())
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Set colors with defaults if not specified
	if g.opts.ForegroundColor != nil {
		qr.ForegroundColor = g.opts.ForegroundColor
	} else {
		qr.ForegroundColor = color.Black
	}
	if g.opts.BackgroundColor != nil {
		qr.BackgroundColor = g.opts.BackgroundColor
	} else {
		qr.BackgroundColor = color.White
	}

	return qr, nil
}

// GeneratePNG generates a QR code and returns it as PNG bytes
func (g *Generator) GeneratePNG(content string) ([]byte, error) {
	qr, err := g.Generate(content)
	if err != nil {
		return nil, err
	}
	return qr.PNG(g.opts.Size)
}
