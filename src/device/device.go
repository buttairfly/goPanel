package device

import (
	"fmt"
	"sync"

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
	GetType() Type
}

// Type is a LedDevice type
type Type string

const (
	// Print debug print device
	Print = Type("print")
	// Serial high level serial tty device
	Serial = Type("serial")
)

// NewLedDevice creates a new Led device
func NewLedDevice(t Type, length int) (LedDevice, error) {
	var pixelDevice LedDevice
	switch t {
	case Print:
		pixelDevice = NewPrintDevice(length)
	case Serial:
		pixelDevice = NewSerialDevice(length)
	default:
		return nil, fmt.Errorf("unkown led device type: %v", t)
	}
	if err := pixelDevice.Open(); err != nil {
		return nil, err
	}
	return pixelDevice, nil
}
