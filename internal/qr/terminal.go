package qr

import (
	"fmt"
	"io"
	"strings"

	"github.com/skip2/go-qrcode"
)

// TerminalConfig configures terminal output
type TerminalConfig struct {
	Invert bool // Invert colors for light terminals (default rendering targets dark terminals)
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

// renderNormal renders using full block characters (2 columns per QR module).
//
// Variable names describe the QR role, not the glyph's appearance. On dark
// terminals "██" renders in the foreground color (light) and "  " renders in
// the background color (dark); so by default modulePx=spaces makes QR modules
// look dark and bgPx=blocks makes QR background look light — the scannable
// "dark modules on light" appearance. --invert swaps the pair for light
// terminals.
func renderNormal(w io.Writer, bitmap [][]bool, invert bool) {
	modulePx := "  "
	bgPx := "██"
	if invert {
		modulePx, bgPx = bgPx, modulePx
	}

	width := len(bitmap[0])
	quietLine := strings.Repeat(bgPx, width+4)

	// Top quiet zone
	fmt.Fprintln(w, quietLine)
	fmt.Fprintln(w, quietLine)

	for _, row := range bitmap {
		fmt.Fprint(w, bgPx+bgPx) // Left quiet zone
		for _, cell := range row {
			if cell {
				fmt.Fprint(w, modulePx)
			} else {
				fmt.Fprint(w, bgPx)
			}
		}
		fmt.Fprintln(w, bgPx+bgPx) // Right quiet zone
	}

	// Bottom quiet zone
	fmt.Fprintln(w, quietLine)
	fmt.Fprintln(w, quietLine)
}

// renderSmall renders using half-block characters (2 QR rows per terminal line).
//
// Same dark-terminal-first convention as renderNormal: bitmap[y][x] == true
// is a QR module and should render dark, which on a dark terminal means the
// "empty" (space-colored) half of the glyph. --invert flips that mapping for
// light terminals.
func renderSmall(w io.Writer, bitmap [][]bool, invert bool) {
	const (
		upperHalf  = "▀"
		lowerHalf  = "▄"
		fullBlock  = "█"
		emptyBlock = " "
	)

	height := len(bitmap)
	width := len(bitmap[0])

	qz := 2 // quiet zone width in columns

	// Quiet zones use the glyph that renders as QR background: fullBlock on
	// dark terminals (default), emptyBlock on light (invert=true).
	quietChar := fullBlock
	if invert {
		quietChar = emptyBlock
	}
	quietLine := strings.Repeat(quietChar, width+qz*2)
	fmt.Fprintln(w, quietLine)

	for y := 0; y < height; y += 2 {
		fmt.Fprint(w, strings.Repeat(quietChar, qz))

		for x := 0; x < width; x++ {
			upper := bitmap[y][x]
			lower := false
			if y+1 < height {
				lower = bitmap[y+1][x]
			}

			var char string
			if invert {
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

		fmt.Fprintln(w, strings.Repeat(quietChar, qz))
	}

	fmt.Fprintln(w, quietLine)
}
