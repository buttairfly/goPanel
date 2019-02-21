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
	if len(s) != NumBytePixel+pos {
		log.Panic("no correct byteslice", len(s), "with offset", NumBytePixel+pos)
	}
	return s[pos : pos+NumBytePixel]
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

// TODO: move to higher frame

// Brighten will make the pixel more light by scale
func (p Pixel) Brighten(scale uint8) {
	//TODO: handle overflow
	p[R] *= scale
	p[G] *= scale
	p[B] *= scale
}

// Dim dimms the pixel by scale
func (p Pixel) Dim(scale uint8) {
	if scale != 0 {
		p[R] /= scale
		p[G] /= scale
		p[B] /= scale
	}
}
