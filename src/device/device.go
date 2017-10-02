package device

import (
	"log"
	"sync"
)

type SpiDevice interface {
	Open() error
	Run(wg *sync.WaitGroup)
	Write(data []byte) error
	Release() error
	SetInput(<-chan []byte)
	GetName() Name
}

type Name string

const (
	PrintDevice  = Name("printDevice")
	WS2801Device = Name("ws2801Device")
)

func NewSpiDevice() (SpiDevice, error) {
	var pixelDevice SpiDevice
	pixelDevice = NewWs2801Device("/dev/spidev0.0", 200)

	err := pixelDevice.Open()
	if err != nil {
		log.Print(err)
		pixelDevice = NewPrintDevice()
		if err := pixelDevice.Open(); err != nil {
			return nil, err
		}
	}
	return pixelDevice, err
}
