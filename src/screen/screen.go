package screen

import (
	"image"

	"github.com/buttairfly/goPanel/src/device"
	"github.com/buttairfly/goPanel/src/screen/raw"
)

type screen struct {
	device        device.SpiDevice
	width, height int
	modules       []Module
	image         raw.Image
}

type Module interface {
	Serialize(image raw.Image) []byte
}

func (s *screen) Bounds() image.Rectangle {
	return image.Rect(0, 0, s.width, s.height)
}
