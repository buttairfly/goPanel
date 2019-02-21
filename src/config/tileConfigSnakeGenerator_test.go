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
		err          error
	}{
		{
			desc: "snake_horizontal_c0_0-0_9-9",
			generator: TileConfigSnakeGenerator{
				startPoint:      image.Point{X: 0, Y: 0},
				endPoint:        image.Point{X: 9, Y: 9},
				direction:       horizontal,
				connectionOrder: 0,
			},
			numPixel:     100,
			expectedFile: ".config.json",
			actualFile:   "actual.config.json",
		},
		{
			desc: "snake_horizontal_c0_19-0_10-9",
			generator: TileConfigSnakeGenerator{
				startPoint:      image.Point{X: 19, Y: 0},
				endPoint:        image.Point{X: 10, Y: 9},
				direction:       horizontal,
				connectionOrder: 0,
			},
			numPixel:     100,
			expectedFile: ".config.json",
			actualFile:   "actual.config.json",
		},
		{
			desc: "snake_vertical___c0_19-0_10-9",
			generator: TileConfigSnakeGenerator{
				startPoint:      image.Point{X: 19, Y: 0},
				endPoint:        image.Point{X: 10, Y: 9},
				direction:       vertical,
				connectionOrder: 0,
			},
			numPixel:     100,
			expectedFile: ".config.json",
			actualFile:   "actual.config.json",
		},
		{
			desc: "snake_vertical___c1_9-0_0-9",
			generator: TileConfigSnakeGenerator{
				startPoint:      image.Point{X: 9, Y: 0},
				endPoint:        image.Point{X: 0, Y: 9},
				direction:       vertical,
				connectionOrder: 1,
			},
			numPixel:     100,
			expectedFile: ".config.json",
			actualFile:   "actual.config.json",
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			expectedFile := fmt.Sprintf("%s%s%s", testFolder, c.desc, c.expectedFile)
			actualFile := fmt.Sprintf("%s%s_%s", testFolder, c.desc, c.actualFile)
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
			assert.True(t, cmp.Equal(readConfig, genConfig), "error read and generated tile config are not equal")
			assert.Equal(t, c.err, genConfig.WriteToFile(actualFile), "error occurred in file write")
			defer os.Remove(actualFile)
			testhelper.Diff(t, expectedFile, actualFile)
		})
	}
}
