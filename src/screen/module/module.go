package module

import (
	"image"

	"github.com/buttairfly/goPanel/src/device"
	"github.com/buttairfly/goPanel/src/screen/raw"
)

type module struct {
	deviceType device.Type
	width      int
	height     int
	numPix     int
	origin     image.Point
	pixLUT     map[image.Point]int
	pixCor     map[raw.ColorPoint]float64
	colLUT     map[raw.RGB8Color]int
}

func (m *module) Serialize(img image.Image) []byte {
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

func (m *module) GetDevice() device.Type {
	return m.deviceType
}

func (m *module) Bounds() image.Rectangle {
	return image.Rect(m.origin.X, m.origin.Y, m.width, m.height)
}

func (m *module) getValue(img image.Image, x, y int, c raw.RGB8Color) byte {
	imgX := m.origin.X + x
	imgY := m.origin.Y + y
	corr := m.pixCor[raw.ColorPoint{Point: image.Point{X: x, Y: y}, C: c}]
	col := raw.RGB8Model.Convert(img.At(imgX, imgY)).(raw.RGB8)
	return byte(float64(col.GetColor(c)) * corr)
}

func (m *module) GetNumPix() int {
	return m.numPix
}
