package device

import (
	"log"
	"sync"
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
	Print  = Type("print")
	WS2801 = Type("ws2801")
)

func NewSpiDevice() (SpiDevice, error) {
	var pixelDevice SpiDevice
	pixelDevice = NewWs2801Device(200)

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
