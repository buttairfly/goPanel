package device

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

const (
	// NumBytePerColor is the number of bytes per pixel
	NumBytePerColor = 3
	// NumByteToRepresentHex is the number of bytes to represent one byte as hex number
	NumByteToRepresentHex = 2
)

var ledDevice LedDevice

// LedDevice interface for all
type LedDevice interface {
	pipepart.PixelPiperSink
	Open() error
	Write(command string) (int, error)
	Close() error
	GetDeviceType() Type
	GetCurrentFrame() hardware.Frame
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
	ledDevice = pixelDevice
	return pixelDevice, nil
}

// GetLedDevice gets the current base LedDevice
func GetLedDevice() LedDevice {
	return ledDevice
}
