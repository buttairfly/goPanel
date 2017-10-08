package module

import (
	"image"

	"fmt"

	"github.com/buttairfly/goPanel/src/device"
	"github.com/buttairfly/goPanel/src/screen/raw"
)

type module struct {
	deviceName device.Type
	width      int
	height     int
	numPix     int
	origin     image.Point
	pixLUT     map[image.Point]int
	pixCor     map[ColorPoint]float64
	colLUT     map[raw.RGB8Color]int
}

type ColorPoint struct {
	image.Point
	rgbType raw.RGB8Color
}

func (cp ColorPoint) String() string {
	return fmt.Sprintf("%v %v", cp.Point, cp.rgbType)
}

func (m *module) Serialize(img raw.Image) []byte {
	bytes := make([]byte, m.numPix)
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			if ledPos, ok := m.pixLUT[image.Point{X: x, Y: y}]; ok {
				for c := range m.colLUT {
					bytes[ledPos+m.colLUT[c]] = m.getValue(img, x, y, c)
				}
			}
		}
	}
	return bytes
}

func (m *module) getValue(img raw.Image, x, y int, c raw.RGB8Color) byte {
	imgX := m.origin.X + x
	imgY := m.origin.Y + y
	corr := m.pixCor[ColorPoint{image.Point{X: x, Y: y}, c}]
	return byte(float64(img.Canvas[imgX][imgY].GetColor(c)) * corr)
}
