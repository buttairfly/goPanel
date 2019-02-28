package palette

import (
	"image/color"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

// Fader interface to fade colors in a palette
type Fader interface {
	Convert(c color.Color) color.Color
	Index(c color.Color) int
	AddColor(c color.Color, pos int)
	AddLastColor(c color.Color)
	NumStepsAndIncrement(granularity int) (int, float64)
	Fade(step float64) color.Color
}

type fader color.Palette

// NewFader creates a Fader from a palette
func NewFader(palette color.Palette) Fader {
	return fader(palette)
}

func (f fader) Convert(c color.Color) color.Color {
	return color.Palette(f).Convert(c)
}

func (f fader) Index(c color.Color) int {
	return color.Palette(f).Index(c)
}

func (f fader) AddColor(c color.Color, pos int) {
	numColors := len(f)
	if pos > numColors {
		pos = numColors
	}
	f = append(f, c)
	if pos != numColors {
		copy(f[pos+1:], f[pos:])
		f[pos] = c
	}
}

func (f fader) AddLastColor(c color.Color) {
	f.AddColor(c, len(f))
}

func (f fader) Fade(step float64) color.Color {
	paletteLen := len(f)
	if paletteLen == 0 {
		return color.Black
	}
	if paletteLen == 1 {
		return f[0]
	}
	if step < 0.0 {
		step = 0.0
	}
	if step > float64(paletteLen-2) {
		step = float64(paletteLen - 2)
	}
	c1, _ := colorful.MakeColor(f[int(step)])
	c2, _ := colorful.MakeColor(f[int(step)+1])

	return c1.BlendHcl(c2, step-math.Trunc(step))
}

func (f fader) NumStepsAndIncrement(granularity int) (int, float64) {
	paletteLength := len(f)
	if paletteLength < 2 {
		return 1, 0.0
	}
	return len(f) * granularity, 1.0 / float64(granularity)
}
