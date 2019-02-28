package hardware

import (
	"image/color"
	"log"

	"github.com/buttairfly/goPanel/src/palette.go"
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
	palette.Color
	ToSlice() []uint8
	ToInt() int
}

// Pixel hardware struct
type pixel []uint8

var _ Pixel = (*pixel)(nil)

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
func NewPixelFromSlice(s []uint8, pos int) Pixel {
	pixPos := pos * NumBytePixel
	if len(s) < pixPos+NumBytePixel-1 {
		log.Fatalf("no correct byteslice %d with offset %d", len(s), pixPos+NumBytePixel-1)
	}
	p := pixel(s[pixPos : pixPos+NumBytePixel])
	return &p
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

// ToSlice converts to an slice color value
func (p *pixel) ToSlice() []uint8 {
	return ([]uint8)(*p)
}

// ToInt converts to an int color value
func (p *pixel) ToInt() int {
	return int((*p)[R])<<16 | int((*p)[G])<<8 | int((*p)[B])
}

// Equals checks wheather the color of another color is identical
func (p *pixel) Equals(c color.Color) bool {
	cr, cg, cb, ca := c.RGBA()
	pr, pg, pb, pa := p.RGBA()
	return cr == pr && cg == pg && cb == pb && ca == pa
}
