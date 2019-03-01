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
	GetIncrements(granularity int) []float64
	Fade(step float64) color.Color
}

type fader color.Palette

const epsilon float64 = 0.000001

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
	if step > float64(paletteLen-1) {
		step = float64(paletteLen - 1)
	}
	fadeValue := step - math.Trunc(step)
	baseStep := int(math.Trunc(step))
	if math.Abs(fadeValue) < epsilon {
		return f[baseStep]
	}

	c1, _ := colorful.MakeColor(f[baseStep])
	c2, _ := colorful.MakeColor(f[baseStep+1])
	return c1.BlendHcl(c2, step-math.Trunc(step)).Clamped()
}

func (f fader) GetIncrements(granularity int) []float64 {
	increments := make([]float64, granularity*len(f))
	position := 0.0
	for i := range increments {
		increments[i] = position
		position += 1.0 / float64(granularity)
	}
	return increments
}
