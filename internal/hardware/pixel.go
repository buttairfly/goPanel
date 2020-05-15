package hardware

import (
	"fmt"
	"image/color"

	"go.uber.org/zap"
)

// Red Green Blue an numBytePixel as constants
const (
	R = iota
	G
	B
	NumBytePixel
)

// Pixel interface implements color interface
type Pixel interface {
	color.Color
	Equals(color.Color) bool
	Slice() []uint8
	Int() int
	Hex() string
}

// Pixel hardware struct
type pixel []uint8

var _ Pixel = (*pixel)(nil)

// NewPixelFromColor creates a new rgb byte struct from color.Color
func NewPixelFromColor(c color.Color) Pixel {
	r, g, b, _ := c.RGBA()
	slice := make(pixel, 0, NumBytePixel)
	p := append(slice, uint8(r>>8), uint8(g>>8), uint8(b>>8))
	return &p
}

// NewPixelFromInts creates a new rgb byte struct from single ints
func NewPixelFromInts(r, g, b uint32) Pixel {
	slice := make(pixel, 0, NumBytePixel)
	p := append(slice, uint8(r), uint8(g), uint8(b))
	return &p
}

// NewPixelFromInt creates a new rgb byte struct from integer
func NewPixelFromInt(c int) Pixel {
	slice := make(pixel, 0, NumBytePixel)
	p := append(slice, uint8(c>>16), uint8(c>>8), uint8(c))
	return &p
}

// NewPixelFromSlice creates a new rgb byte struct from uint8 slice
func NewPixelFromSlice(s []uint8, pos int, logger *zap.Logger) Pixel {
	pixPos := pos * NumBytePixel
	pixPosEnd := pixPos + NumBytePixel
	if len(s) < pixPosEnd-1 {
		logger.Sugar().Fatalf("no correct byteslice %d with offset %d", len(s), pixPosEnd-1)
	}
	p := pixel(s[pixPos:pixPosEnd])
	return &p
}

// NewPixelFromHex parses a "html" hex color-string, either in the 3 "#f0c" or 6 "#ff1034" digits form.
func NewPixelFromHex(hex string) (Pixel, error) {
	format := "#%02x%02x%02x"
	if len(hex) == 4 {
		format = "#%1x%1x%1x"
	}

	var r, g, b uint8
	n, err := fmt.Sscanf(hex, format, &r, &g, &b)
	if err != nil {
		return nil, err
	}
	if n != 3 {
		return nil, fmt.Errorf("color: %v is not a hex-color", hex)
	}

	slice := make(pixel, 0, NumBytePixel)
	p := append(slice, uint8(r), uint8(g), uint8(b))
	return &p, nil
}

// RGBA implements color.Color interface
func (p *pixel) RGBA() (r, g, b, a uint32) {
	r = uint32((*p)[R])
	r |= r << 8
	g = uint32((*p)[G])
	g |= g << 8
	b = uint32((*p)[B])
	b |= b << 8
	a = uint32(0xFFFF)
	return
}

// Slice converts to an slice color value
func (p *pixel) Slice() []uint8 {
	return ([]uint8)(*p)
}

// Int converts to an int color value
func (p *pixel) Int() int {
	return int((*p)[R])<<16 | int((*p)[G])<<8 | int((*p)[B])
}

// Hex returns the hex "html" representation of the color, as in #ff0080
func (p *pixel) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x", (*p)[R], (*p)[G], (*p)[B])
}

// Equals checks wheather the color of another color is identical
func (p *pixel) Equals(c color.Color) bool {
	cr, cg, cb, ca := c.RGBA()
	pr, pg, pb, pa := p.RGBA()
	return cr == pr && cg == pg && cb == pb && ca == pa
}
