package palette

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	emptyPalette    = &fader{palette: nil, granularity: 1, wrapping: false}
	greenPalette    = &fader{palette: []color.Color{green}, granularity: 1, wrapping: false}
	redGreenPalette = &fader{palette: []color.Color{red, green}, granularity: 1, wrapping: false}
	red             = color.RGBA64{R: 0xffff, G: 0x0000, B: 0x0000, A: 0xffff}
	green           = color.RGBA64{R: 0x0000, G: 0xffff, B: 0x0000, A: 0xffff}
	redGreen0Point1 = color.RGBA64{R: 0xfedd, G: 0x4021, B: 0x0000, A: 0xffff}
	redGreen0Point9 = color.RGBA64{R: 0x6a81, G: 0xf02a, B: 0x0000, A: 0xffff}
)

func TestNewFader(t *testing.T) {
	cases := []struct {
		desc        string
		colors      []color.Color
		granularity int
		wrapping    bool
		expected    Fader
		step        float64
		fadeColor   color.Color
	}{
		{
			desc:      "empty_fader",
			expected:  emptyPalette,
			step:      1.0,
			fadeColor: color.Black,
		},
		{
			desc:      "green_color_fader_0.1",
			colors:    []color.Color{green},
			expected:  greenPalette,
			step:      0.1,
			fadeColor: green,
		},
		{
			desc:      "green_color_fader_0.0",
			colors:    []color.Color{green},
			expected:  greenPalette,
			step:      0.0,
			fadeColor: green,
		},
		{
			desc:      "green_color_fader_-0.1",
			colors:    []color.Color{green},
			expected:  greenPalette,
			step:      -0.1,
			fadeColor: green,
		},
		{
			desc:      "green_color_fader_1.0",
			colors:    []color.Color{green},
			expected:  greenPalette,
			step:      1.0,
			fadeColor: green,
		},
		{
			desc:      "green_color_fader_1.1",
			colors:    []color.Color{green},
			expected:  greenPalette,
			step:      1.1,
			fadeColor: green,
		},
		{
			desc:      "red_green_color_fader_0.1",
			colors:    []color.Color{red, green},
			expected:  redGreenPalette,
			step:      0.1,
			fadeColor: redGreen0Point1,
		},
		{
			desc:      "red_green_color_fader_0.9",
			colors:    []color.Color{red, green},
			expected:  redGreenPalette,
			step:      0.9,
			fadeColor: redGreen0Point9,
		},
		{
			desc:      "red_green_color_fader_0.0",
			colors:    []color.Color{red, green},
			expected:  redGreenPalette,
			step:      0.0,
			fadeColor: red,
		},
		{
			desc:      "red_green_color_fader_-0.1",
			colors:    []color.Color{red, green},
			expected:  redGreenPalette,
			step:      -0.1,
			fadeColor: red,
		},
		{
			desc:      "red_green_color_fader_1.0",
			colors:    []color.Color{red, green},
			expected:  redGreenPalette,
			step:      1.0,
			fadeColor: green,
		},
		{
			desc:      "red_green_color_fader_1.1",
			colors:    []color.Color{red, green},
			expected:  redGreenPalette,
			step:      1.1,
			fadeColor: green,
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			fader := NewFader(c.colors, c.granularity, c.wrapping)
			assert.Equal(t, c.expected, fader)
			fr, fg, fb, fa := fader.Fade(c.step).RGBA()
			er, eg, eb, ea := c.fadeColor.RGBA()
			assert.Equal(t, er, fr, "R")
			assert.Equal(t, eg, fg, "G")
			assert.Equal(t, eb, fb, "B")
			assert.Equal(t, ea, fa, "A")
		})
	}
}

func TestFaderIncrements(t *testing.T) {
	cases := []struct {
		desc        string
		colors      []color.Color
		wrapping    bool
		granularity int
		expected    []float64
		expectedLen int
	}{
		{
			desc:        "increment_empty_fader_100",
			granularity: 100,
			expected:    []float64{0.0},
			expectedLen: 1,
		},
		{
			desc:        "increment_green_color_fader_2",
			colors:      []color.Color{green},
			granularity: 2,
			expected:    []float64{0.0},
			expectedLen: 1,
		},
		{
			desc:        "increment_red_color_fader_10",
			colors:      []color.Color{red},
			granularity: 10,
			expected:    []float64{0.0},
			expectedLen: 1,
		},

		{
			desc:        "increment_green_color_fader_10",
			colors:      []color.Color{red, green},
			granularity: 10,
			expected:    []float64{0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0},
			expectedLen: 11,
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			fader := NewFader(c.colors, c.granularity, c.wrapping)
			increments := fader.GetIncrements()
			assert.Equal(t, c.expectedLen, len(increments))
			assert.Equal(t, c.expected, increments)
		})
	}
}
