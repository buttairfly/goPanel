package palette

/*
var (
	emptyPalette        = &fader{palette: nil, granularity: 1, wrapping: false}
	greenPalette        = &fader{palette: []color.Color{green}, granularity: 1, wrapping: false}
	redGreenPalette     = &fader{palette: []color.Color{red, green}, granularity: 1, wrapping: false}
	redGreenWrapPalette = &fader{palette: []color.Color{red, green}, granularity: 1, wrapping: true}
	red                 = color.RGBA64{R: 0xffff, G: 0x0000, B: 0x0000, A: 0xffff}
	green               = color.RGBA64{R: 0x0000, G: 0xffff, B: 0x0000, A: 0xffff}
	redGreen0Point1     = color.RGBA64{R: 0xfedd, G: 0x4021, B: 0x0000, A: 0xffff}
	redGreen0Point9     = color.RGBA64{R: 0x6a81, G: 0xf02a, B: 0x0000, A: 0xffff}
)

func TestNewFader(t *testing.T) {
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
			expected: emptyPalette,
			stepMap:  map[float64]color.Color{1.0: color.Black},
		},
		{
			desc:     "green_color_fader",
			colors:   []color.Color{green},
			expected: greenPalette,
			stepMap: map[float64]color.Color{
				0.1:  green,
				0.0:  green,
				-0.1: green,
				1.0:  green,
				1.1:  green,
			},
		},
		{
			desc:     "red_green_color_fader",
			colors:   []color.Color{red, green},
			expected: redGreenPalette,
			stepMap: map[float64]color.Color{
				0.1:  redGreen0Point1,
				0.9:  redGreen0Point9,
				0.0:  red,
				-0.1: red,
				1.0:  green,
				1.1:  green,
			},
		},
		// wrapping
		{
			desc:     "wrapping_red_green_color_fader",
			colors:   []color.Color{red, green},
			wrapping: true,
			expected: redGreenWrapPalette,
			stepMap: map[float64]color.Color{
				0.1:  redGreen0Point1,
				0.9:  redGreen0Point9,
				0.0:  red,
				-0.1: redGreen0Point1,
				1.0:  green,
				1.1:  redGreen0Point9,
				2.1:  redGreen0Point1,
				-5.1: redGreen0Point9,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			fader := NewFader(c.colors, c.granularity, c.wrapping)
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

		{
			desc:        "increment_wrapping_green_color_fader_10",
			wrapping:    true,
			colors:      []color.Color{red, green},
			granularity: 10,
			expected:    []float64{0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0, 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8, 1.9},
			expectedLen: 20,
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
*/
