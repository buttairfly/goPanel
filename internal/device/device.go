package device

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

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
	Run(cancelCtx context.Context, wg *sync.WaitGroup)
	Write(command string) (int, error)
	Close() error
	SetInput(<-chan hardware.Frame)
	GetType() Type
}

// NewLedDevice creates a new Led device
func NewLedDevice(ledDeviceConfig *LedDeviceConfig, length int, logger *zap.Logger) (LedDevice, error) {
	var pixelDevice LedDevice
	switch ledDeviceConfig.Type {
	case Print:
		pixelDevice = NewPrintDevice(length, ledDeviceConfig.PrintConfig, logger)
	case Serial:
		pixelDevice = NewSerialDevice(length, ledDeviceConfig.SerialConfig, logger)
	default:
		return nil, fmt.Errorf("unkown led device type: %v", ledDeviceConfig.Type)
	}
	if err := pixelDevice.Open(); err != nil {
		return nil, err
	}
	return pixelDevice, nil
}
