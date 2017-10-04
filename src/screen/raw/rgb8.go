package raw

import "log"

type RGB8 struct {
	R, G, B byte
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
		return 0
	}
}

type RGB8Color string

const (
	R = RGB8Color("r")
	G = RGB8Color("g")
	B = RGB8Color("b")
)

var RGB8Space = map[RGB8Color]struct{}{
	R: {},
	G: {},
	B: {},
}
