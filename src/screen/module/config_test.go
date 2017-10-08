package module

import (
	"image"
	"testing"

	"io/ioutil"

	"os"

	"errors"

	"github.com/buttairfly/goPanel/src/device"
	"github.com/buttairfly/goPanel/src/screen/raw"
	"github.com/buttairfly/goPanel/src/testhelper"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewModulesFromConfig(t *testing.T) {
	cases := []struct {
		desc       string
		config     config
		resultFile string
		actualFile string
		fileName   string
		err        error
	}{
		{
			desc: "20x10_print_ws2801",
			config: config{[]moduleConfig{
				{
					DeviceType: device.Print,
					Height:     10,
					Width:      20,
					LineOrder:  LineOrderXY,
					Mirror:     MirrorNo,
					Rotation:   Rotate0,
					Origin:     image.Point{X: 0, Y: 0},
					ColLUT:     map[raw.RGB8Color]int{raw.R: 0, raw.G: 1, raw.B: 2},
				},
				{
					DeviceType: device.WS2801,
					Height:     10,
					Width:      10,
					LineOrder:  LineOrderSnake,
					Mirror:     MirrorNo,
					Rotation:   Rotate270,
					Origin:     image.Point{X: 0, Y: 0},
					ColLUT:     map[raw.RGB8Color]int{raw.R: 0, raw.G: 1, raw.B: 2},
				},
				{
					DeviceType: device.WS2801,
					Height:     10,
					Width:      10,
					LineOrder:  LineOrderSnake,
					Mirror:     MirrorNo,
					Rotation:   Rotate270,
					Origin:     image.Point{X: 10, Y: 0},
					ColLUT:     map[raw.RGB8Color]int{raw.R: 0, raw.G: 1, raw.B: 2},
				},
			}},
			resultFile: "expected.module",
			actualFile: "actual.module",
			fileName:   "20x10_print_ws2801.json",
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			if testhelper.RecordCall() {
				t.Logf("Write config to file %v", c.fileName)
				require.NoError(t, c.config.WriteToFile(c.fileName))
			}
			modules, err := NewModulesFromConfig(c.fileName)
			data := []byte(modulesToStr(modules))
			if testhelper.RecordCall() {
				t.Logf("Write result to file %v", c.resultFile)
				require.NoError(t, ioutil.WriteFile(c.resultFile, data, 0644))
			}
			assert.Equal(t, c.err, err, "error occurred")
			assert.NoError(t, err, "error reading file")
			assert.NoError(t, ioutil.WriteFile(c.actualFile, data, 0644))
			defer os.Remove(c.actualFile)
			testhelper.Diff(t, c.resultFile, c.actualFile)
		})
	}
}

func modulesToStr(modules []module) string {
	spew.Config.SortKeys = true
	return spew.Sdump(modules)
}

func TestConfigRotate(t *testing.T) {
	cases := []struct {
		desc            string
		width, height   int
		origin          image.Point
		rotation        rotation
		input, expected image.Point
		err             error
	}{
		{
			desc:     "invalid",
			width:    2,
			height:   3,
			rotation: rotation("invalid"),
			input:    image.Point{0, 0},
			err:      errors.New("no correct Rotation=(invalid) set"),
		},
		{
			desc:     "1x1_90Grad_0_0",
			width:    1,
			height:   1,
			rotation: Rotate90,
			input:    image.Point{0, 0},
			expected: image.Point{0, 0},
		},
		{
			desc:     "2x2_90Grad_0_0",
			width:    2,
			height:   2,
			rotation: Rotate90,
			input:    image.Point{0, 0},
			expected: image.Point{1, 0},
		},
		{
			desc:     "2x2_90Grad_1_0",
			width:    2,
			height:   2,
			rotation: Rotate90,
			input:    image.Point{1, 0},
			expected: image.Point{1, 1},
		},
		{
			desc:     "2x3_0Grad_0_0",
			width:    2,
			height:   3,
			rotation: Rotate0,
			input:    image.Point{0, 0},
			expected: image.Point{0, 0},
		},
		{
			desc:     "2x3_90Grad_0_0",
			width:    2,
			height:   3,
			rotation: Rotate90,
			input:    image.Point{0, 0},
			expected: image.Point{2, 0},
		},
		{
			desc:     "2x3_180Grad_0_0",
			width:    2,
			height:   3,
			rotation: Rotate180,
			input:    image.Point{0, 0},
			expected: image.Point{1, 2},
		},
		{
			desc:     "2x3_270Grad_0_0_with_offset",
			width:    2,
			height:   3,
			rotation: Rotate270,
			input:    image.Point{0, 0},
			expected: image.Point{0, 1},
		},
		{
			desc:     "2x3_90Grad_0_0",
			width:    2,
			height:   3,
			origin:   image.Point{10, 0},
			rotation: Rotate90,
			input:    image.Point{0, 0},
			expected: image.Point{2, 0},
		},
		{
			desc:     "10x10_90Grad_0_0",
			width:    10,
			height:   10,
			origin:   image.Point{10, 0},
			rotation: Rotate90,
			input:    image.Point{0, 0},
			expected: image.Point{9, 0},
		},
		{
			desc:     "10x10_90Grad_9_0",
			width:    10,
			height:   10,
			origin:   image.Point{10, 0},
			rotation: Rotate90,
			input:    image.Point{9, 0},
			expected: image.Point{9, 9},
		},
		{
			desc:     "10x9_90Grad_8_1",
			width:    10,
			height:   9,
			origin:   image.Point{10, 0},
			rotation: Rotate90,
			input:    image.Point{8, 1},
			expected: image.Point{7, 8},
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			mc := moduleConfig{
				Width:    c.width,
				Height:   c.height,
				Origin:   c.origin,
				Rotation: c.rotation,
			}
			newP, err := mc.translateRotation(c.input)
			assert.Equal(t, c.err, err, "error equals")
			assert.Equal(t, c.expected, newP, "point is expected")
		})
	}
}
