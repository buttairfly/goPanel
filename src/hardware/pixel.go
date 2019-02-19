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

// NewPixelFromSlice creates a new rgb byte struct from byte slice
func NewPixelFromSlice(s []byte, pos int) Pixel {
	if len(s) != NumBytePixel+pos {
		log.Panic("no correct byteslice", len(s), "with offset", NumBytePixel+pos)
	}
	return s[pos : pos+NumBytePixel]
}

func (p Pixel) toSlice() []uint8 {
	return ([]uint8)(p)
}

func (p Pixel) toInt() int {
	return int(p[R])<<16 | int(p[G])<<8 | int(p[B])
}

// TODO: move to higher frame
func (p Pixel) brighten(scale uint8) {
	//TODO: handle overflow
	p[R] *= scale
	p[G] *= scale
	p[B] *= scale
}

func (p Pixel) dim(scale uint8) {
	if scale != 0 {
		p[R] /= scale
		p[G] /= scale
		p[B] /= scale
	}
}

func (p Pixel) equals(other Pixel) bool {
	return p[R] == other[R] && p[G] == other[G] && p[B] == other[B]
}
