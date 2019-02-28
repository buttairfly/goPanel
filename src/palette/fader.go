package palette

import "image/color"

// Fader interface to fade colors in a palette
type Fader interface {
	Convert(c color.Color) color.Color
	Index(c color.Color) int
}

type fader struct {
	palette color.Palette
	weights []int
}

// NewFader creates a Fader from a palette
func NewFader(palette color.Palette) Fader {
	weights := make([]int, len(palette))
	for i := range weights {
		weights[i] = 1
	}
	return &fader{
		palette: palette,
		weights: weights,
	}
}

func (f *fader) Convert(c color.Color) color.Color {
	return f.palette.Convert(c)
}

func (f *fader) Index(c color.Color) int {
	return f.palette.Index(c)
}

func (f *fader) AddColor(c color.Color, pos int) {
	numColors := len(f.palette)
	if pos > numColors {
		pos = numColors
	}
	/*if f.palette[pos].Equals(c) {

	}*/
}
