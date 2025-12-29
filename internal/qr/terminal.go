package qr

import (
	"fmt"
	"io"
	"strings"

	"github.com/skip2/go-qrcode"
)

// TerminalConfig configures terminal output
type TerminalConfig struct {
	Invert bool // Invert colors for dark terminals
	Small  bool // Use half-block characters for compact output
}

// RenderTerminal renders a QR code to the terminal
func RenderTerminal(w io.Writer, qr *qrcode.QRCode, cfg TerminalConfig) {
	bitmap := qr.Bitmap()

	if cfg.Small {
		renderSmall(w, bitmap, cfg.Invert)
	} else {
		renderNormal(w, bitmap, cfg.Invert)
	}
}

// renderNormal renders using full block characters (2x2 pixels per character)
func renderNormal(w io.Writer, bitmap [][]bool, invert bool) {
	white := "██"
	black := "  "
	if invert {
		white, black = black, white
	}

	// Add quiet zone
	width := len(bitmap[0])
	quietLine := strings.Repeat(white, width+4)

	// Top quiet zone
	fmt.Fprintln(w, quietLine)
	fmt.Fprintln(w, quietLine)

	for _, row := range bitmap {
		fmt.Fprint(w, white+white) // Left quiet zone
		for _, cell := range row {
			if cell {
				fmt.Fprint(w, black)
			} else {
				fmt.Fprint(w, white)
			}
		}
		fmt.Fprintln(w, white+white) // Right quiet zone
	}

	// Bottom quiet zone
	fmt.Fprintln(w, quietLine)
	fmt.Fprintln(w, quietLine)
}

// renderSmall renders using half-block characters (2 rows per line)
func renderSmall(w io.Writer, bitmap [][]bool, invert bool) {
	// Unicode half blocks: ▀ (upper), ▄ (lower), █ (full), " " (empty)
	const (
		upperHalf = "▀"
		lowerHalf = "▄"
		fullBlock = "█"
		emptyBlock = " "
	)

	height := len(bitmap)
	width := len(bitmap[0])

	// Quiet zone width
	qz := 2

	// Top quiet zone (represented as full blocks in inverted sense)
	quietChar := fullBlock
	if invert {
		quietChar = emptyBlock
	}
	quietLine := strings.Repeat(quietChar, width+qz*2)
	fmt.Fprintln(w, quietLine)

	// Process two rows at a time
	for y := 0; y < height; y += 2 {
		// Left quiet zone
		fmt.Fprint(w, strings.Repeat(quietChar, qz))

		for x := 0; x < width; x++ {
			upper := bitmap[y][x]
			lower := false
			if y+1 < height {
				lower = bitmap[y+1][x]
			}

			// In bitmap: true = black (QR module), false = white (background)
			// For terminal: we want white background, black modules
			// Without invert: white=█, black=" "
			// With invert: white=" ", black="█"

			var char string
			if invert {
				// Inverted: black modules should be visible (█), white should be empty
				switch {
				case upper && lower:
					char = fullBlock
				case upper && !lower:
					char = upperHalf
				case !upper && lower:
					char = lowerHalf
				default:
					char = emptyBlock
				}
			} else {
				// Normal: white background (█), black modules (space)
				switch {
				case !upper && !lower:
					char = fullBlock
				case !upper && lower:
					char = upperHalf
				case upper && !lower:
					char = lowerHalf
				default:
					char = emptyBlock
				}
			}
			fmt.Fprint(w, char)
		}

		// Right quiet zone
		fmt.Fprintln(w, strings.Repeat(quietChar, qz))
	}

	// Bottom quiet zone
	fmt.Fprintln(w, quietLine)
}
