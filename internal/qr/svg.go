package qr

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skip2/go-qrcode"
)

// SaveSVG writes the QR code as an SVG file.
// size is the rendered pixel size; the SVG is square (size x size).
func SaveSVG(qr *qrcode.QRCode, filename string, size int) error {
	dir := filepath.Dir(filename)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	data := ToSVG(qr, size)
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write SVG file: %w", err)
	}
	return nil
}

// ToSVG renders the QR code to an SVG byte slice.
// The bitmap from skip2/go-qrcode already includes the standard quiet zone.
// ForegroundColor / BackgroundColor on qr are honored; a fully-transparent
// background omits the background rect so the QR sits on page transparency.
func ToSVG(qr *qrcode.QRCode, size int) []byte {
	bitmap := qr.Bitmap()
	n := len(bitmap)

	fgFill, fgOp := colorToSVG(qr.ForegroundColor)
	bgFill, bgOp := colorToSVG(qr.BackgroundColor)

	var b bytes.Buffer
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	fmt.Fprintf(&b,
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" shape-rendering="crispEdges">`+"\n",
		size, size, n, n)

	// Emit background rect unless the background is fully transparent.
	if bgFill != "none" {
		fmt.Fprintf(&b, `<rect width="%d" height="%d" fill="%s"`, n, n, bgFill)
		if bgOp != "" {
			fmt.Fprintf(&b, ` fill-opacity="%s"`, bgOp)
		}
		b.WriteString(`/>` + "\n")
	}

	// Coalesce horizontal runs of dark modules into compact path commands.
	fmt.Fprintf(&b, `<path fill="%s"`, fgFill)
	if fgOp != "" {
		fmt.Fprintf(&b, ` fill-opacity="%s"`, fgOp)
	}
	b.WriteString(` d="`)
	for y := 0; y < n; y++ {
		x := 0
		for x < n {
			if bitmap[y][x] {
				start := x
				for x < n && bitmap[y][x] {
					x++
				}
				runLen := x - start
				fmt.Fprintf(&b, "M%d %dh%dv1h-%dz", start, y, runLen, runLen)
			} else {
				x++
			}
		}
	}
	b.WriteString(`"/>` + "\n")
	b.WriteString(`</svg>` + "\n")
	return b.Bytes()
}
