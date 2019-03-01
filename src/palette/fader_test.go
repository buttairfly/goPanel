package palette

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFader(t *testing.T) {
	red := color.RGBA64{R: 0xffff, G: 0x00, B: 0x00, A: 0xffff}
	green := color.RGBA64{R: 0x00, G: 0xffff, B: 0x00, A: 0xffff}
	green_red_0_1 := color.RGBA64{R: 0xfedd, G: 0x4021, B: 0x00, A: 0xffff}
	green_red_0_9 := color.RGBA64{R: 0x6a81, G: 0xf02a, B: 0x00, A: 0xffff}
	cases := []struct {
		desc      string
		colors    []color.Color
		expected  Fader
		step      float64
		fadeColor color.Color
	}{
		{
			desc:      "empty_fader",
			expected:  fader(nil),
			step:      1.0,
			fadeColor: color.Black,
		},
		{
			desc:      "green_color_fader_0.1",
			colors:    []color.Color{green},
			expected:  fader([]color.Color{green}),
			step:      0.1,
			fadeColor: green,
		},
		{
			desc:      "green_color_fader_0.0",
			colors:    []color.Color{green},
			expected:  fader([]color.Color{green}),
			step:      0.0,
			fadeColor: green,
		},
		{
			desc:      "green_color_fader_-0.1",
			colors:    []color.Color{green},
			expected:  fader([]color.Color{green}),
			step:      -0.1,
			fadeColor: green,
		},
		{
			desc:      "green_color_fader_1.0",
			colors:    []color.Color{green},
			expected:  fader([]color.Color{green}),
			step:      1.0,
			fadeColor: green,
		},
		{
			desc:      "green_color_fader_1.1",
			colors:    []color.Color{green},
			expected:  fader([]color.Color{green}),
			step:      1.1,
			fadeColor: green,
		},
		{
			desc:      "red_green_color_fader_0.1",
			colors:    []color.Color{red, green},
			expected:  fader([]color.Color{red, green}),
			step:      0.1,
			fadeColor: green_red_0_1,
		},
		{
			desc:      "red_green_color_fader_0.9",
			colors:    []color.Color{red, green},
			expected:  fader([]color.Color{red, green}),
			step:      0.9,
			fadeColor: green_red_0_9,
		},
		{
			desc:      "red_green_color_fader_0.0",
			colors:    []color.Color{red, green},
			expected:  fader([]color.Color{red, green}),
			step:      0.0,
			fadeColor: red,
		},
		{
			desc:      "red_green_color_fader_-0.1",
			colors:    []color.Color{red, green},
			expected:  fader([]color.Color{red, green}),
			step:      -0.1,
			fadeColor: red,
		},
		{
			desc:      "red_green_color_fader_1.0",
			colors:    []color.Color{red, green},
			expected:  fader([]color.Color{red, green}),
			step:      1.0,
			fadeColor: green,
		},
		{
			desc:      "red_green_color_fader_1.1",
			colors:    []color.Color{red, green},
			expected:  fader([]color.Color{red, green}),
			step:      1.1,
			fadeColor: green,
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			fader := NewFader(c.colors)
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
