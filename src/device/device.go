package device

import (
	"fmt"
	"sync"

	"github.com/buttairfly/goPanel/src/config"
	"github.com/buttairfly/goPanel/src/hardware"
)

const (
	// NumBytePerColor is the number of bytes per pixel
	NumBytePerColor = 3
	// NumByteToRepresentHex is the number of bytes to represent one byte as hex number
	NumByteToRepresentHex = 2
)

// LedDevice interface for all
type LedDevice interface {
	Open() error
	Run(wg *sync.WaitGroup)
	Write(data []byte) (int, error)
	Close() error
	SetInput(<-chan hardware.Frame)
	GetType() config.Type
}

// NewLedDevice creates a new Led device
func NewLedDevice(t config.Type, length int) (LedDevice, error) {
	var pixelDevice LedDevice
	switch t {
	case config.Print:
		pixelDevice = NewPrintDevice(length)
	case config.Serial:
		pixelDevice = NewSerialDevice(length)
	default:
		return nil, fmt.Errorf("unkown led device type: %v", t)
	}
	if err := pixelDevice.Open(); err != nil {
		return nil, err
	}
	return pixelDevice, nil
}
