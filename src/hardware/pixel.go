package hardware

import (
	"log"
)

// Red Green Blue an numBytePixel as constants
const (
	R = iota
	G
	B
	NumBytePixel
)

// Pixel hardware struct
type Pixel []uint8

// NewPixelFromInts creates a new rgb byte struct from single ints
func NewPixelFromInts(r, g, b uint32) Pixel {
	byteSlice := make([]byte, 0, NumBytePixel)
	return append(byteSlice, uint8(r), uint8(g), uint8(b))
}

// NewPixelFromInt creates a new rgb byte struct from integer
func NewPixelFromInt(c int) Pixel {
	byteSlice := make([]byte, 0, NumBytePixel)
	return append(byteSlice, uint8(c>>16), uint8(c>>8), uint8(c))
}

// NewPixelFromSlice creates a new rgb byte struct from uint8 slice
func NewPixelFromSlice(s []uint8, pos int) Pixel {
	pixPos := pos * NumBytePixel
	if len(s) < pixPos+NumBytePixel-1 {
		log.Fatalf("no correct byteslice %d with offset %d", len(s), pixPos+NumBytePixel-1)
	}
	return s[pixPos : pixPos+NumBytePixel]
}

// ToSlice converts to an slice color value
func (p Pixel) ToSlice() []uint8 {
	return ([]uint8)(p)
}

// ToInt converts to an int color value
func (p Pixel) ToInt() int {
	return int(p[R])<<16 | int(p[G])<<8 | int(p[B])
}

// Equals checks wheather the color of orhter and this pixel is identical
func (p Pixel) Equals(other Pixel) bool {
	return p[R] == other[R] && p[G] == other[G] && p[B] == other[B]
}
