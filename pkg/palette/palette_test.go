package palette

import (
	"fmt"
	"math"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/stretchr/testify/assert"
)

var (
	pBlack            = colorful.Color{R: 0, G: 0, B: 0}
	pRed              = colorful.Color{R: 1.0, G: 0, B: 0}
	pGreen            = colorful.Color{R: 0, G: 1.0, B: 0}
	cgRedGreen        = colorful.Color{R: 0.8415372297953739, G: 0.6529456392653107, B: 0}
	emptyPalette      = []paletteColor(nil)
	green0Palette     = []paletteColor{{pos: 0.0, color: pGreen}}
	green1Palette     = []paletteColor{{pos: 1.0, color: pGreen}}
	green1red0Palette = []paletteColor{{pos: 1.0, color: pGreen}, {pos: 0.0, color: pRed}}
	red0green1Palette = []paletteColor{{pos: 0.0, color: pRed}, {pos: 1.0, color: pGreen}}
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
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			palette := NewPalette()
			for _, paletteColor := range c.paletteColors {
				palette.AddAt(paletteColor.color, paletteColor.pos)
			}
			assert.Equal(t, &(c.expected), palette)
			for step, expectedColor := range c.stepMap {
				t.Run(fmt.Sprintf("%f", step), func(t *testing.T) {
					assert.Equal(t, expectedColor.Hex(), palette.Blend(step).Hex(),
						"Hex representation, colorful representation %v",
						palette.Blend(step),
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
