package raw

import (
	"log"

	"github.com/buttairfly/goPanel/src/helper"
)

type RGB8 struct {
	R, G, B byte
}

const (
	rgb8min = 0
	rgb8max = 0xff
	rgbaMin = 0
	rgbaMax = 0xffff
)

func (c *RGB8) RGBA() (r, g, b, a uint32) {
	return c.to16(R), c.to16(G), c.to16(B), rgbaMax
}

func (c *RGB8) to16(color RGB8Color) uint32 {
	return uint32(helper.IntMap(
		int(c.GetColor(color)),
		rgb8min, rgb8max, rgbaMin, rgbaMax))
}

func (c *RGB8) GetColor(color RGB8Color) byte {
	switch color {
	case R:
		return c.R
	case G:
		return c.G
	case B:
		return c.B
	default:
		log.Panicf("Unknown rgb8 color: %v", color)
	}
	return 0
}

func (c *RGB8) SetColor(val byte, color RGB8Color) {
	switch color {
	case R:
		c.R = val
	case G:
		c.G = val
	case B:
		c.B = val
	default:
		log.Panicf("Unknown rgb8 color: %v", color)
	}
}

type RGB8Color int

const (
	R RGB8Color = iota
	G
	B
)

var RGB8Space = map[RGB8Color]struct{}{
	R: {},
	G: {},
	B: {},
}

func (c RGB8Color) String() string {
	switch c {
	case R:
		return "r"
	case G:
		return "g"
	case B:
		return "b"
	default:
		return "unknown color"
	}
}

func RGB8ColorFromString(s string) RGB8Color {
	switch s {
	case "r":
		return R
	case "g":
		return G
	case "b":
		return B
	default:
		log.Panicf("Unknown rgb8 color string: %v", s)
	}
	return -1
}
