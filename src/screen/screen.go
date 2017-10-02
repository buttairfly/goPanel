package screen

import (
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
	Byte(image raw.Image) []byte
}
