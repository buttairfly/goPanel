package module

import (
	"image"

	"github.com/buttairfly/goPanel/src/device"
	"github.com/buttairfly/goPanel/src/screen/raw"
)

type module struct {
	deviceName device.Name
	width      int
	height     int
	numPix     int
	origin     image.Point
	pixLUT     map[image.Point]int
	pixCor     map[ColorPoint]float64
	colLUT     map[raw.RGB8Color]int
}

type ColorPoint struct {
	point    image.Point
	rgbOrder raw.RGB8Color
}

func NewModule(configPath string) *module {

}

func (m *module) Byte(img raw.Image) []byte {
	bytes := make([]byte, m.numPix)
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			if ledPos, ok := m.pixLUT[image.Point{x, y}]; ok {
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
	corr := m.pixCor[ColorPoint{image.Point{x, y}, c}]
	return byte(float64(img.Canvas[imgX][imgY].GetColor(c)) * corr)
}
