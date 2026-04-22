package qr

import (
	"image/color"
	"testing"
)

func TestParseColor(t *testing.T) {
	tests := []struct {
		in   string
		want color.NRGBA
	}{
		// Named
		{"black", color.NRGBA{0, 0, 0, 255}},
		{"WHITE", color.NRGBA{255, 255, 255, 255}},
		{"red", color.NRGBA{255, 0, 0, 255}},
		{"transparent", color.NRGBA{0, 0, 0, 0}},
		// 3-digit hex (each expanded to double)
		{"#000", color.NRGBA{0, 0, 0, 255}},
		{"#fff", color.NRGBA{255, 255, 255, 255}},
		{"#f00", color.NRGBA{255, 0, 0, 255}},
		{"#abc", color.NRGBA{0xaa, 0xbb, 0xcc, 255}},
		// 6-digit hex
		{"#ff0000", color.NRGBA{255, 0, 0, 255}},
		{"#00ff00", color.NRGBA{0, 255, 0, 255}},
		{"#123456", color.NRGBA{0x12, 0x34, 0x56, 255}},
		// 8-digit hex with alpha
		{"#ff000080", color.NRGBA{255, 0, 0, 0x80}},
		{"#00000000", color.NRGBA{0, 0, 0, 0}},
		// Leading '#' optional
		{"ff0000", color.NRGBA{255, 0, 0, 255}},
		{"f00", color.NRGBA{255, 0, 0, 255}},
		// Uppercase
		{"#FF0000", color.NRGBA{255, 0, 0, 255}},
		// Whitespace tolerated
		{"  #fff  ", color.NRGBA{255, 255, 255, 255}},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := ParseColor(tt.in)
			if err != nil {
				t.Fatalf("ParseColor(%q) error: %v", tt.in, err)
			}
			gotN, ok := got.(color.NRGBA)
			if !ok {
				t.Fatalf("ParseColor(%q) returned %T, want color.NRGBA", tt.in, got)
			}
			if gotN != tt.want {
				t.Errorf("ParseColor(%q) = %+v, want %+v", tt.in, gotN, tt.want)
			}
		})
	}
}

func TestParseColorErrors(t *testing.T) {
	bad := []string{
		"",
		"notacolor",
		"#",
		"#ff",
		"#ggg",
		"#1234567",   // 7 digits, unsupported length
		"#123456789", // 9 digits
		"rgb(0,0,0)",
	}
	for _, s := range bad {
		t.Run(s, func(t *testing.T) {
			_, err := ParseColor(s)
			if err == nil {
				t.Errorf("ParseColor(%q) expected error, got nil", s)
			}
		})
	}
}

func TestColorToSVG(t *testing.T) {
	tests := []struct {
		name       string
		in         color.Color
		wantFill   string
		wantOpNone bool // true means fill-opacity should be empty
	}{
		{"nil defaults to black", nil, "#000000", true},
		{"opaque black", color.NRGBA{0, 0, 0, 255}, "#000000", true},
		{"opaque white", color.NRGBA{255, 255, 255, 255}, "#ffffff", true},
		{"opaque red", color.NRGBA{255, 0, 0, 255}, "#ff0000", true},
		{"transparent", color.NRGBA{0, 0, 0, 0}, "none", true},
		{"semi-transparent", color.NRGBA{255, 0, 0, 128}, "#ff0000", false},
		{"from premultiplied RGBA opaque", color.RGBA{255, 0, 0, 255}, "#ff0000", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fill, op := colorToSVG(tt.in)
			if fill != tt.wantFill {
				t.Errorf("fill = %q, want %q", fill, tt.wantFill)
			}
			if (op == "") != tt.wantOpNone {
				t.Errorf("fill-opacity presence mismatch: got %q, wantEmpty=%v", op, tt.wantOpNone)
			}
		})
	}
}
