package module

import (
	"image"
	"testing"

	"github.com/buttairfly/goPanel/src/device"
	"github.com/buttairfly/goPanel/src/screen/raw"
	"github.com/buttairfly/goPanel/src/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewModulesFromConfig(t *testing.T) {
	cases := []struct {
		desc            string
		config          config
		expectedModules []module
		fileName        string
		err             error
	}{
		{
			desc: "variants",
			config: config{[]moduleConfig{
				{
					DeviceName: device.Print,
					Height:     5,
					Width:      3,
					LineOrder:  LineOrderXY,
					Mirror:     MirrorNo,
					Rotation:   Rotate0,
					Origin:     image.Point{X: 0, Y: 0},
					ColLUT:     map[raw.RGB8Color]int{raw.R: 0, raw.G: 1, raw.B: 2},
				},
			}},
			expectedModules: []module{
				{
					deviceName: "print",
					height:     5,
					width:      3,
					pixLUT: map[image.Point]int{
						{0, 0}: 0,
						{1, 0}: 1,
						{2, 0}: 2,
						{0, 1}: 3,
						{1, 1}: 4,
						{2, 1}: 5,
						{0, 2}: 6,
						{1, 2}: 7,
						{2, 2}: 8,
						{0, 3}: 9,
						{1, 3}: 10,
						{2, 3}: 11,
						{0, 4}: 12,
						{1, 4}: 13,
						{2, 4}: 14,
					},
					colLUT: map[raw.RGB8Color]int{raw.R: 0, raw.G: 1, raw.B: 2},
					pixCor: map[ColorPoint]float64{
						{image.Point{0, 0}, raw.R}: 1.0,
						{image.Point{1, 0}, raw.R}: 1.0,
						{image.Point{2, 0}, raw.R}: 1.0,
						{image.Point{0, 1}, raw.R}: 1.0,
						{image.Point{1, 1}, raw.R}: 1.0,
						{image.Point{2, 1}, raw.R}: 1.0,
						{image.Point{0, 2}, raw.R}: 1.0,
						{image.Point{1, 2}, raw.R}: 1.0,
						{image.Point{2, 2}, raw.R}: 1.0,
						{image.Point{0, 3}, raw.R}: 1.0,
						{image.Point{1, 3}, raw.R}: 1.0,
						{image.Point{2, 3}, raw.R}: 1.0,
						{image.Point{0, 4}, raw.R}: 1.0,
						{image.Point{1, 4}, raw.R}: 1.0,
						{image.Point{2, 4}, raw.R}: 1.0,
						{image.Point{0, 0}, raw.G}: 1.0,
						{image.Point{1, 0}, raw.G}: 1.0,
						{image.Point{2, 0}, raw.G}: 1.0,
						{image.Point{0, 1}, raw.G}: 1.0,
						{image.Point{1, 1}, raw.G}: 1.0,
						{image.Point{2, 1}, raw.G}: 1.0,
						{image.Point{0, 2}, raw.G}: 1.0,
						{image.Point{1, 2}, raw.G}: 1.0,
						{image.Point{2, 2}, raw.G}: 1.0,
						{image.Point{0, 3}, raw.G}: 1.0,
						{image.Point{1, 3}, raw.G}: 1.0,
						{image.Point{2, 3}, raw.G}: 1.0,
						{image.Point{0, 4}, raw.G}: 1.0,
						{image.Point{1, 4}, raw.G}: 1.0,
						{image.Point{2, 4}, raw.G}: 1.0,
						{image.Point{0, 0}, raw.B}: 1.0,
						{image.Point{1, 0}, raw.B}: 1.0,
						{image.Point{2, 0}, raw.B}: 1.0,
						{image.Point{0, 1}, raw.B}: 1.0,
						{image.Point{1, 1}, raw.B}: 1.0,
						{image.Point{2, 1}, raw.B}: 1.0,
						{image.Point{0, 2}, raw.B}: 1.0,
						{image.Point{1, 2}, raw.B}: 1.0,
						{image.Point{2, 2}, raw.B}: 1.0,
						{image.Point{0, 3}, raw.B}: 1.0,
						{image.Point{1, 3}, raw.B}: 1.0,
						{image.Point{2, 3}, raw.B}: 1.0,
						{image.Point{0, 4}, raw.B}: 1.0,
						{image.Point{1, 4}, raw.B}: 1.0,
						{image.Point{2, 4}, raw.B}: 1.0,
					},
				},
			},
			fileName: "printDevice.json",
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			if testhelper.RecordCall() {
				t.Logf("Write config to file %v", c.fileName)
				require.NoError(t, c.config.WriteToFile(c.fileName))
			}
			modules, err := NewModulesFromConfig(c.fileName)
			assert.Equal(t, c.err, err, "error occoured")
			assert.Equal(t, c.expectedModules, modules, "modules are not equal")
		})
	}
}
