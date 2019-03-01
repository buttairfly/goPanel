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

type fader struct {
	palette  color.Palette
	wrapping bool
}

const epsilon float64 = 0.000001

// NewFader creates a Fader from a palette
func NewFader(palette color.Palette, wrapping bool) Fader {
	return &fader{palette: palette, wrapping: wrapping}
}

func (f fader) Convert(c color.Color) color.Color {
	return color.Palette(f.palette).Convert(c)
}

func (f fader) Index(c color.Color) int {
	return color.Palette(f.palette).Index(c)
}

func (f fader) AddColor(c color.Color, pos int) {
	numColors := len(f.palette)
	if pos > numColors {
		pos = numColors
	}
	f.palette = append(f.palette, c)
	if pos != numColors {
		copy(f.palette[pos+1:], f.palette[pos:])
		f.palette[pos] = c
	}
}

func (f fader) AddLastColor(c color.Color) {
	f.AddColor(c, len(f.palette))
}

func (f fader) Fade(step float64) color.Color {
	paletteLen := len(f.palette)
	if paletteLen == 0 {
		return color.Black
	}
	if paletteLen == 1 {
		return f.palette[0]
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
		return f.palette[baseStep]
	}

	c1, _ := colorful.MakeColor(f.palette[baseStep])
	c2, _ := colorful.MakeColor(f.palette[baseStep+1])
	return c1.BlendHcl(c2, step-math.Trunc(step)).Clamped()
}

func (f fader) GetIncrements(granularity int) []float64 {
	lenPalette := len(f.palette)
	if lenPalette < 2 {
		return []float64{0.0}
	}
	numSteps := granularity * lenPalette
	if f.wrapping {
		numSteps = granularity * (lenPalette + 1)
	}
	increments := make([]float64, numSteps)
	position := 0.0
	for i := range increments {
		increments[i] = position
		position += 1.0 / float64(granularity)
	}
	return increments
}
