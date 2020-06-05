package fader

import "image/color"

type Fader interface {
	Convert(c color.Color) color.Color
	Index(c color.Color) int
	AddColor(c color.Color, pos int)
	Fade(step float64) color.Color
	GetIncrements() []float64
}
