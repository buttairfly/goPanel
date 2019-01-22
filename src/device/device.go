package device

import (
	"fmt"
	"sync"
)

const (
	NumBytePerColor = 3
)

type SpiDevice interface {
	Open() error
	Run(wg *sync.WaitGroup)
	Write(data []byte) (int, error)
	Close() error
	SetInput(<-chan []byte)
	GetType() Type
}

type Type string

const (
	Print  = Type("print")  // Print debug print device
	WS2801 = Type("ws2801") // direct spi serial device
	Serial = Type("serial") // high level serial tty device
)

func NewSpiDevice(t Type, length int) (SpiDevice, error) {
	var pixelDevice SpiDevice
	switch t {
	case Print:
		pixelDevice = NewPrintDevice(length)
		if err := pixelDevice.Open(); err != nil {
			return nil, err
		}
		return pixelDevice, nil
	case WS2801:
		pixelDevice = NewWs2801Device(length)
		if err := pixelDevice.Open(); err != nil {
			return nil, err
		}
	case Serial:
		pixelDevice = NewSerialDevice(length)
		if err := pixelDevice.Open(); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unkown spi device type: %v", t)
	}
	return pixelDevice, nil
}
