package qr

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

// namedColors is a small set of CSS-style color names intentionally kept
// minimal — the hex path covers everything else.
var namedColors = map[string]color.NRGBA{
	"black":       {R: 0, G: 0, B: 0, A: 255},
	"white":       {R: 255, G: 255, B: 255, A: 255},
	"red":         {R: 255, G: 0, B: 0, A: 255},
	"green":       {R: 0, G: 128, B: 0, A: 255},
	"blue":        {R: 0, G: 0, B: 255, A: 255},
	"yellow":      {R: 255, G: 255, B: 0, A: 255},
	"cyan":        {R: 0, G: 255, B: 255, A: 255},
	"magenta":     {R: 255, G: 0, B: 255, A: 255},
	"transparent": {R: 0, G: 0, B: 0, A: 0},
}

// ParseColor converts a user-supplied color string into a color.Color.
//
// Accepted forms (case-insensitive, leading '#' optional):
//
//	#rgb          shorthand hex (each digit doubled: f → ff)
//	#rrggbb       standard hex
//	#rrggbbaa     hex with straight alpha
//	<name>        black, white, red, green, blue, yellow, cyan, magenta, transparent
//
// The returned Color is a color.NRGBA (straight alpha), which preserves the
// user's intent under downstream compositing regardless of the consumer.
func ParseColor(s string) (color.Color, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return nil, fmt.Errorf("empty color")
	}
	if c, ok := namedColors[s]; ok {
		return c, nil
	}
	hex := strings.TrimPrefix(s, "#")
	switch len(hex) {
	case 3:
		r, err1 := parseHexByte(hex[0:1] + hex[0:1])
		g, err2 := parseHexByte(hex[1:2] + hex[1:2])
		b, err3 := parseHexByte(hex[2:3] + hex[2:3])
		if err := firstErr(err1, err2, err3); err != nil {
			return nil, fmt.Errorf("invalid hex color %q: %w", s, err)
		}
		return color.NRGBA{R: r, G: g, B: b, A: 255}, nil
	case 6:
		r, err1 := parseHexByte(hex[0:2])
		g, err2 := parseHexByte(hex[2:4])
		b, err3 := parseHexByte(hex[4:6])
		if err := firstErr(err1, err2, err3); err != nil {
			return nil, fmt.Errorf("invalid hex color %q: %w", s, err)
		}
		return color.NRGBA{R: r, G: g, B: b, A: 255}, nil
	case 8:
		r, err1 := parseHexByte(hex[0:2])
		g, err2 := parseHexByte(hex[2:4])
		b, err3 := parseHexByte(hex[4:6])
		a, err4 := parseHexByte(hex[6:8])
		if err := firstErr(err1, err2, err3, err4); err != nil {
			return nil, fmt.Errorf("invalid hex color %q: %w", s, err)
		}
		return color.NRGBA{R: r, G: g, B: b, A: a}, nil
	}
	return nil, fmt.Errorf("invalid color %q (use a name like 'red', or hex like '#f00', '#ff0000', '#ff0000ff')", s)
}

func parseHexByte(s string) (uint8, error) {
	v, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		return 0, err
	}
	return uint8(v), nil
}

func firstErr(errs ...error) error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
}

// colorToSVG renders a color.Color into an SVG fill string plus an optional
// fill-opacity value. Fully transparent input returns ("none", "").
func colorToSVG(c color.Color) (fill, fillOpacity string) {
	if c == nil {
		return "#000000", ""
	}
	n := color.NRGBAModel.Convert(c).(color.NRGBA)
	if n.A == 0 {
		return "none", ""
	}
	fill = fmt.Sprintf("#%02x%02x%02x", n.R, n.G, n.B)
	if n.A < 255 {
		fillOpacity = fmt.Sprintf("%.3f", float64(n.A)/255.0)
	}
	return fill, fillOpacity
}
