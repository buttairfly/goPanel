package palette

import (
	"fmt"
	"math"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/stretchr/testify/assert"
)

var (
	pBlack              = colorful.Color{R: 0, G: 0, B: 0}
	pRed                = colorful.Color{R: 1.0, G: 0, B: 0}
	pGreen              = colorful.Color{R: 0, G: 1.0, B: 0}
	pBlue               = colorful.Color{R: 0, G: 0, B: 1.00}
	cgRedGreen          = colorful.Color{R: 0.8415372297953739, G: 0.6529456392653107, B: 0}
	emptyPalette        = []paletteColor(nil)
	green0Palette       = []paletteColor{{pos: 0.0, color: pGreen}}
	green1Palette       = []paletteColor{{pos: 1.0, color: pGreen}}
	red0green1Palette   = []paletteColor{{pos: 0.0, color: pRed}, {pos: 1.0, color: pGreen}}
	red0green0_5Palette = []paletteColor{{pos: 0.0, color: pRed}, {pos: 0.5, color: pGreen}}
	red0_5green1Palette = []paletteColor{{pos: 0.5, color: pRed}, {pos: 1.0, color: pGreen}}
	green1red0Palette   = []paletteColor{{pos: 1.0, color: pGreen}, {pos: 0.0, color: pRed}}

	pRGB0_1    = colorful.Color{R: 0.9451598952119228, G: 0.4773624195946065, B: 0}
	pRGB0_2    = colorful.Color{R: 0.7692833682910185, G: 0.731063820539255, B: 0}
	pRGB0_3    = colorful.Color{R: 0.41603785430985246, G: 0.938156725445569, B: 0}
	pRGB0_4    = colorful.Color{R: 0, G: 0.9120522553612874, B: 0.573324315612675}
	pRGB0_5    = colorful.Color{R: 0, G: 0.7297272229698231, B: 1}
	pRGB0_6    = colorful.Color{R: 0, G: 0.4758562396465739, B: 1}
	pRGB0_7    = colorful.Color{R: 0.515587672665168, G: 0, B: 0.912547403000855}
	pRGB0_8    = colorful.Color{R: 0.9209495225040252, G: 0, B: 0.6063025156213984}
	pRGB0_9    = colorful.Color{R: 1, G: 0, B: 0.31337268895733733}
	rgbPalette = []paletteColor{{pos: 0.0, color: pRed}, {pos: 1.0 / 3.0, color: pGreen}, {pos: 2.0 / 3.0, color: pBlue}, {pos: 1.0, color: pRed}}
)

var _ Palette = (*palette)(nil)

func TestPaletteAdd(t *testing.T) {
	cases := []struct {
		desc          string
		paletteColors []paletteColor
		expected      palette
		stepMap       map[float64]colorful.Color
	}{
		{
			desc:     "empty_palette",
			expected: emptyPalette,
			stepMap: map[float64]colorful.Color{
				0.0: pBlack,
				0.5: pBlack,
				1.0: pBlack,
			},
		},
		{
			desc:          "green0_palette",
			paletteColors: green0Palette,
			expected:      green0Palette,
			stepMap: map[float64]colorful.Color{
				0.0: pGreen,
				0.5: pGreen,
				1.0: pGreen,
			},
		},
		{
			desc:          "green1_palette",
			paletteColors: green1Palette,
			expected:      green1Palette,
			stepMap: map[float64]colorful.Color{
				0.0: pGreen,
				0.5: pGreen,
				1.0: pGreen,
			},
		},
		{
			desc:          "red0_green1_palette",
			paletteColors: green1red0Palette,
			expected:      red0green1Palette,
			stepMap: map[float64]colorful.Color{
				0.0: pRed,
				0.5: cgRedGreen,
				1.0: pGreen,
			},
		},
		{
			desc:          "red0_green0.5_palette",
			paletteColors: red0green0_5Palette,
			expected:      red0green0_5Palette,
			stepMap: map[float64]colorful.Color{
				0.0: pRed,
				0.5: pGreen,
				1.0: pGreen,
			},
		},
		{
			desc:          "red0.5_green1_palette",
			paletteColors: red0_5green1Palette,
			expected:      red0_5green1Palette,
			stepMap: map[float64]colorful.Color{
				0.0: pRed,
				0.5: pRed,
				1.0: pGreen,
			},
		},
		{
			desc:          "rgb_palette",
			paletteColors: rgbPalette,
			expected:      rgbPalette,
			stepMap: map[float64]colorful.Color{
				0.0:       pRed,
				1.0 / 3.0: pGreen,
				2.0 / 3.0: pBlue,
				1.0:       pRed,
				0.1:       pRGB0_1,
				0.2:       pRGB0_2,
				0.3:       pRGB0_3,
				0.4:       pRGB0_4,
				0.5:       pRGB0_5,
				0.6:       pRGB0_6,
				0.7:       pRGB0_7,
				0.8:       pRGB0_8,
				0.9:       pRGB0_9,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			palette := NewPalette()
			for _, paletteColor := range c.paletteColors {
				palette.AddAt(paletteColor.color, paletteColor.pos)
			}
			assert.Equal(t, &(c.expected), palette)
			for step, expectedColor := range c.stepMap {
				t.Run(fmt.Sprintf("%1.3f", step), func(t *testing.T) {
					calculatedColor := palette.Blend(step)
					assert.Equal(t, expectedColor.Hex(), calculatedColor.Hex(),
						"Hex representation, colorful representation %+v",
						calculatedColor,
					)
				})
			}
		})
	}
}

func TestGuaranteeBetween0And1(t *testing.T) {
	cases := []struct {
		desc   string
		posMap map[float64]float64
	}{
		{
			desc: "different_positions_within_0_and_1",
			posMap: map[float64]float64{
				0.0: 0.0,
				0.1: 0.1,
				0.5: 0.5,
				0.9: 0.9,
				1.0: 1.0,
			},
		},
		{
			desc: "different_positions_within_-1_and_2",
			posMap: map[float64]float64{
				1.1: 0.10000000000000009,
				1.5: 0.5,
				1.9: 0.8999999999999999,
				2.0: 0,

				-0.1: 0.9,
				-0.5: 0.5,
				-0.9: 0.09999999999999998,
				-1.0: 1.0,
			},
		},
		{
			desc: "different_positions_others",
			posMap: map[float64]float64{
				math.NaN(): 0,
				-11.0:      1,
				-10.0:      1,
				10:         0,
				11.0:       0,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			for pos, expectedPos := range c.posMap {
				t.Run(fmt.Sprintf("%f", pos), func(t *testing.T) {
					assert.Equal(t, expectedPos, guaranteeBetween0And1(pos))
				})
			}
		})
	}
}
