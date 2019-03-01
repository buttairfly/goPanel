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
	GetIncrements() []float64
	Fade(step float64) color.Color
}

type fader struct {
	palette     color.Palette
	granularity int
	wrapping    bool
}

// NewFader creates a Fader from a palette
func NewFader(palette color.Palette, granularity int, wrapping bool) Fader {
	if granularity < 1 {
		granularity = 1
	}
	return &fader{palette: palette, granularity: granularity, wrapping: wrapping}
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
	if math.Abs(fadeValue) < 0.1/float64(f.granularity) {
		return f.palette[baseStep]
	}

	c1, _ := colorful.MakeColor(f.palette[baseStep])
	c2, _ := colorful.MakeColor(f.palette[baseStep+1])
	return c1.BlendHcl(c2, step-math.Trunc(step)).Clamped()
}

func (f fader) GetIncrements() []float64 {
	lenPalette := len(f.palette)
	if lenPalette < 2 {
		return []float64{0.0}
	}
	numSteps := f.granularity*(lenPalette-1) + 1
	if f.wrapping {
		numSteps = f.granularity * lenPalette
	}
	increments := make([]float64, numSteps)
	for num := range increments {
		increments[num] = float64(num) / float64(f.granularity)
	}
	return increments
}
