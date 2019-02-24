package config

import (
	"fmt"
	"image"
	"os"
	"testing"

	"github.com/buttairfly/goPanel/src/testhelper"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTileConfigSnakeMapFile(t *testing.T) {
	const testFolder = "testdata/"
	cases := []struct {
		desc         string
		generator    TileConfigSnakeGenerator
		expectedFile string
		actualFile   string
		numPixel     int
		testPixel    map[string]int
		err          error
	}{
		{
			desc: "snake_horizontal_c0_0-0_10-10",
			generator: TileConfigSnakeGenerator{
				startPoint:      image.Point{X: 0, Y: 0},
				endPoint:        image.Point{X: 10, Y: 10},
				direction:       horizontal,
				connectionOrder: 0,
			},
			numPixel: 100,
			testPixel: map[string]int{
				" 0": 0, " 9": 9,
				"90": 99, "99": 90,
			},
			expectedFile: ".config.json",
			actualFile:   "actual.config.json",
		},
		{
			desc: "snake_horizontal_c0_20-0_10-10",
			generator: TileConfigSnakeGenerator{
				startPoint:      image.Point{X: 20, Y: 0},
				endPoint:        image.Point{X: 10, Y: 10},
				direction:       horizontal,
				connectionOrder: 0,
			},
			numPixel: 100,
			testPixel: map[string]int{
				" 0": 9, " 9": 0,
				"90": 90, "99": 99,
			},
			expectedFile: ".config.json",
			actualFile:   "actual.config.json",
		},
		{
			desc: "snake_vertical_c0_20-0_10-10",
			generator: TileConfigSnakeGenerator{
				startPoint:      image.Point{X: 20, Y: 0},
				endPoint:        image.Point{X: 10, Y: 10},
				direction:       vertical,
				connectionOrder: 0,
			},
			testPixel: map[string]int{
				" 0": 99, " 9": 0,
				"90": 90, "99": 9,
			},
			numPixel:     100,
			expectedFile: ".config.json",
			actualFile:   "actual.config.json",
		},
		{
			desc: "snake_vertical_c1_10-0_0-10",
			generator: TileConfigSnakeGenerator{
				startPoint:      image.Point{X: 10, Y: 0},
				endPoint:        image.Point{X: 0, Y: 10},
				direction:       vertical,
				connectionOrder: 1,
			},
			testPixel: map[string]int{
				" 0": 99, " 9": 0,
				"90": 90, "99": 9,
			},
			numPixel:     100,
			expectedFile: ".config.json",
			actualFile:   "actual.config.json",
		},
		{
			desc: "snake_vertical_c0_0-0_3-3",
			generator: TileConfigSnakeGenerator{
				startPoint:      image.Point{X: 0, Y: 0},
				endPoint:        image.Point{X: 3, Y: 3},
				direction:       vertical,
				connectionOrder: 0,
			},
			testPixel: map[string]int{
				" 0": 0, " 1": 5, " 2": 6,
				" 3": 1, " 4": 4, " 5": 7,
				" 6": 2, " 7": 3, " 8": 8,
			},
			numPixel:     9,
			expectedFile: ".config.json",
			actualFile:   "actual.config.json",
		},
		{
			desc: "snake_vertical_c0_0-0_4-4",
			generator: TileConfigSnakeGenerator{
				startPoint:      image.Point{X: 0, Y: 0},
				endPoint:        image.Point{X: 4, Y: 4},
				direction:       vertical,
				connectionOrder: 0,
			},
			testPixel: map[string]int{
				" 0": 0, " 1": 7, " 2": 8, " 3": 15,
				" 4": 1, " 5": 6, " 6": 9, " 7": 14,
				" 8": 2, " 9": 5, "10": 10, "11": 13,
				"12": 3, "13": 4, "14": 11, "15": 12,
			},
			numPixel:     16,
			expectedFile: ".config.json",
			actualFile:   "actual.config.json",
		},
		{
			desc: "snake_vertical_c0_3-3_0-0",
			generator: TileConfigSnakeGenerator{
				startPoint:      image.Point{X: 3, Y: 3},
				endPoint:        image.Point{X: 0, Y: 0},
				direction:       vertical,
				connectionOrder: 0,
			},
			testPixel: map[string]int{
				" 0": 8, " 1": 3, " 2": 2,
				" 3": 7, " 4": 4, " 5": 1,
				" 6": 6, " 7": 5, " 8": 0,
			},
			numPixel:     9,
			expectedFile: ".config.json",
			actualFile:   "actual.config.json",
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			expectedFile := fmt.Sprintf("%stile.%s%s", testFolder, c.desc, c.expectedFile)
			actualFile := fmt.Sprintf("%stile.%s_%s", testFolder, c.desc, c.actualFile)
			genConfig, err := NewTileConfigSnakeMapFile(c.generator)
			require.NoError(t, err)
			require.NotNil(t, genConfig)

			assert.Equal(t, c.numPixel, genConfig.NumHardwarePixel(), "error not enough pixel NumHardwarePixel")
			assert.Equal(t, c.numPixel, len(genConfig.GetLedStripeMap()), "error not enough pixel GetLedStripeMap")
			assert.Equal(t, c.generator.connectionOrder, genConfig.GetConnectionOrder(), "error not correct connection order")
			assert.Equal(t, image.Rectangle{Min: c.generator.startPoint, Max: c.generator.endPoint}.Canon(), genConfig.GetBounds(), "error not correct bounds")

			if testhelper.RecordCall() {
				t.Logf("Write Config to file %v", expectedFile)
				require.NoError(t, genConfig.WriteToFile(expectedFile))
			}

			readConfig, err2 := NewTileConfigFromPath(expectedFile)
			require.NoError(t, err2)

			t.Log(cmp.Diff(readConfig, genConfig))
			assert.True(t, cmp.Equal(readConfig, genConfig), "error read and generated tile config are not equal")
			assert.Equal(t, c.err, genConfig.WriteToFile(actualFile), "error occurred in file write")
			defer os.Remove(actualFile)
			testhelper.Diff(t, expectedFile, actualFile)

			for k, v := range c.testPixel {
				assert.Equal(t, v, genConfig.GetLedStripeMap()[k], "testPixel not equal at pos: %s", k)
			}
		})
	}
}
