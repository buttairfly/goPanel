package fader

/*
import (
	"image/color"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

// HCLFader struct to fade colors in a palette
type HCLFader struct {
	name        string
	palette     color.Palette
	granularity int
	wrapping    bool
	//luminance   float64
}

// NewHCLFader creates a HCLFader from a palette
func NewHCLFader(name string, palette color.Palette, granularity int, wrapping bool) Fader {
	if granularity < 1 {
		granularity = 1
	}
	return &HCLFader{name: name, palette: palette, granularity: granularity, wrapping: wrapping}
}

// Convert calls the Palette color convert function
func (f *HCLFader) Convert(c color.Color) color.Color {
	return color.Palette(f.palette).Convert(c)
}

// Index gets the index of the palette
func (f *HCLFader) Index(c color.Color) int {
	return color.Palette(f.palette).Index(c)
}

// AddColor adds a new color to existing fader at pos
func (f *HCLFader) AddColor(c color.Color, pos int) {
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

// AddLastColor appends a color to the end of the fader
func (f *HCLFader) AddLastColor(c color.Color) {
	f.AddColor(c, len(f.palette))
}

// Fade fades between the colors with equal distance and dependend on step [0.0,1.0]
func (f *HCLFader) Fade(step float64) color.Color {
	paletteLen := len(f.palette)
	if paletteLen == 0 {
		return color.Black
	}
	if paletteLen == 1 {
		return f.palette[0]
	}

	numColorsToFade := float64(paletteLen)
	if !f.wrapping {
		numColorsToFade -= 1.0
	}

	if step < 0.0 {
		if f.wrapping {
			for step < 0.0 { // this may happen more than once
				step += numColorsToFade // it is a subtract since step is negative
			}
		} else {
			step = 0.0
		}
	}
	if step > numColorsToFade {
		if f.wrapping {
			for step > numColorsToFade { // this may happen more then once
				step -= numColorsToFade
			}
		} else {
			step = numColorsToFade
		}
	}
	fadeValue := step - math.Trunc(step)
	baseStep := int(math.Trunc(step))
	if math.Abs(fadeValue) < 0.1/float64(f.granularity) {
		return f.palette[baseStep]
	}
	var c2 colorful.Color

	c1, _ := colorful.MakeColor(f.palette[baseStep])
	if f.wrapping && baseStep+1 == paletteLen {
		c2, _ = colorful.MakeColor(f.palette[0])
	} else {
		c2, _ = colorful.MakeColor(f.palette[baseStep+1])
	}
	return c1.BlendHcl(c2, step-math.Trunc(step)).Clamped()
}

// GetIncrements returns the steps array for the Fader
func (f *HCLFader) GetIncrements() []float64 {
	lenPalette := len(f.palette)
	if lenPalette < 2 {
		return []float64{0.0}
	}
	var numSteps int
	if f.wrapping {
		numSteps = f.granularity * lenPalette
	} else {
		numSteps = f.granularity*(lenPalette-1) + 1
	}
	increments := make([]float64, numSteps)
	for num := range increments {
		increments[num] = float64(num) / float64(f.granularity)
	}
	return increments
}
*/
