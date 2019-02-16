package device

import (
	"fmt"
	"log"
	"sync"

	"github.com/luismesas/goPi/spi"
)

const (
	ws2801Mode  = uint8(0)
	ws2801Bpw   = uint8(8)
	ws2801Speed = uint32(1000000)
	ws2801Delay = uint16(1000)
)

type ws2801Device struct {
	device *spi.SPIDevice
	numLed int
	input  <-chan []byte
}

// NewWs2801Device creates a new direct spi WS2801 device
func NewWs2801Device(numLed int) LedDevice {
	ws := new(ws2801Device)
	ws.device = spi.NewSPIDevice(spi.DEFAULT_BUS, spi.DEFAULT_CHIP)
	ws.numLed = numLed
	return ws
}

func (ws *ws2801Device) Open() error {
	err := ws.device.Open()
	if err != nil {
		return err
	}

	err = ws.device.SetMode(ws2801Mode)
	if err != nil {
		return err
	}

	err = ws.device.SetBitsPerWord(ws2801Bpw)
	if err != nil {
		return err
	}

	err = ws.device.SetSpeed(ws2801Speed)
	if err != nil {
		return err
	}

	return nil
}

func (ws *ws2801Device) Close() error {
	return ws.device.Close()
}

func (ws *ws2801Device) Write(data []byte) (int, error) {
	if len(data) != ws.numLed*NumBytePerColor {
		return 0, fmt.Errorf(
			"could not write %v bytes of data, %v is needed",
			len(data), ws.numLed*NumBytePerColor)
	}
	array := [3]byte{0, 0, 0}
	for i := 0; i < 3; i++ {
		array[i] = data[i]
	}
	_, err := ws.device.Send(array)
	return len(data), err
}

func (ws *ws2801Device) SetInput(input <-chan []byte) {
	ws.input = input
}

func (ws *ws2801Device) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer ws.Close()
	for buffer := range ws.input {
		_, err := ws.Write(buffer)
		if err != nil {
			log.Panic(err)
		}
	}
}

func (ws *ws2801Device) GetType() Type {
	return WS2801
}
