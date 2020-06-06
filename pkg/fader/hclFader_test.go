package fader

/*
import (
	"fmt"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	cRed            = color.RGBA64{R: 0xffff, G: 0x0000, B: 0x0000, A: 0xffff}
	cGreen          = color.RGBA64{R: 0x0000, G: 0xffff, B: 0x0000, A: 0xffff}
	redGreen0Point1 = color.RGBA64{R: 0xfedd, G: 0x4021, B: 0x0000, A: 0xffff}
	redGreen0Point9 = color.RGBA64{R: 0x6a81, G: 0xf02a, B: 0x0000, A: 0xffff}
)

var (
	greenPalette        = &HCLFader{palette: []color.Color{cGreen}, granularity: 1, wrapping: false}
	redGreenPalette     = &HCLFader{palette: []color.Color{cRed, cGreen}, granularity: 1, wrapping: false}
	redGreenWrapPalette = &HCLFader{palette: []color.Color{cRed, cGreen}, granularity: 1, wrapping: true}
)

func TestNewHCLFader(t *testing.T) {
	cases := []struct {
		desc        string
		colors      []color.Color
		granularity int
		wrapping    bool
		expected    Fader
		stepMap     map[float64]color.Color
	}{
		{
			desc:     "empty_fader",
			expected: &HCLFader{palette: nil, granularity: 1, wrapping: false},
			stepMap:  map[float64]color.Color{1.0: color.Black},
		},
		{
			desc:     "green_color_fader",
			colors:   []color.Color{cGreen},
			expected: greenPalette,
			stepMap: map[float64]color.Color{
				0.1:  cGreen,
				0.0:  cGreen,
				-0.1: cGreen,
				1.0:  cGreen,
				1.1:  cGreen,
			},
		},
		{
			desc:     "red_green_color_fader",
			colors:   []color.Color{cRed, cGreen},
			expected: redGreenPalette,
			stepMap: map[float64]color.Color{
				0.1:  redGreen0Point1,
				0.9:  redGreen0Point9,
				0.0:  cRed,
				-0.1: cRed,
				1.0:  cGreen,
				1.1:  cGreen,
			},
		},
		// wrapping
		{
			desc:     "wrapping_red_green_color_fader",
			colors:   []color.Color{cRed, cGreen},
			wrapping: true,
			expected: redGreenWrapPalette,
			stepMap: map[float64]color.Color{
				0.1:  redGreen0Point1,
				0.9:  redGreen0Point9,
				0.0:  cRed,
				-0.1: redGreen0Point1,
				1.0:  cGreen,
				1.1:  redGreen0Point9,
				2.1:  redGreen0Point1,
				-5.1: redGreen0Point9,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			fader := NewHCLFader("", c.colors, c.granularity, c.wrapping)
			assert.Equal(t, c.expected, fader)
			for step, expectedColor := range c.stepMap {
				t.Run(fmt.Sprintf("%s_%f", c.desc, step), func(t *testing.T) {
					fr, fg, fb, fa := fader.Fade(step).RGBA()
					er, eg, eb, ea := expectedColor.RGBA()
					assert.Equal(t, er, fr, "R")
					assert.Equal(t, eg, fg, "G")
					assert.Equal(t, eb, fb, "B")
					assert.Equal(t, ea, fa, "A")
				})
			}
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
			colors:      []color.Color{cGreen},
			granularity: 2,
			expected:    []float64{0.0},
			expectedLen: 1,
		},
		{
			desc:        "increment_red_color_fader_10",
			colors:      []color.Color{cRed},
			granularity: 10,
			expected:    []float64{0.0},
			expectedLen: 1,
		},

		{
			desc:        "increment_green_color_fader_10",
			colors:      []color.Color{cRed, cGreen},
			granularity: 10,
			expected:    []float64{0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0},
			expectedLen: 11,
		},

		{
			desc:        "increment_wrapping_green_color_fader_10",
			wrapping:    true,
			colors:      []color.Color{cRed, cGreen},
			granularity: 10,
			expected:    []float64{0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0, 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8, 1.9},
			expectedLen: 20,
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			fader := NewHCLFader("", c.colors, c.granularity, c.wrapping)
			increments := fader.GetIncrements()
			assert.Equal(t, c.expectedLen, len(increments))
			assert.Equal(t, c.expected, increments)
		})
	}
}
*/
