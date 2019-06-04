package device

import (
	"fmt"
	"sync"

	"github.com/buttairfly/goPanel/internal/config"
	"github.com/buttairfly/goPanel/internal/hardware"
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
func NewLedDevice(deviceConfig *config.DeviceConfig, length int) (LedDevice, error) {
	var pixelDevice LedDevice
	switch deviceConfig.Type {
	case config.Print:
		pixelDevice = NewPrintDevice(length)
	case config.Serial:
		pixelDevice = NewSerialDevice(length, deviceConfig.SerialConfig)
	default:
		return nil, fmt.Errorf("unkown led device type: %v", deviceConfig.Type)
	}
	if err := pixelDevice.Open(); err != nil {
		return nil, err
	}
	return pixelDevice, nil
}
