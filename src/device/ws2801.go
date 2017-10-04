package device

import (
	"fmt"

	"sync"

	"github.com/luismesas/goPi/spi"
)

const (
	ws2801Mode  = uint8(0)
	ws2801Bpw   = uint8(8)
	ws2801Speed = uint32(1000000)
	ws2801Delay = uint16(1000)

	WS2801NumColor = 3
)

type ws2801 struct {
	device *spi.SPIDevice
	numLed int
	input  <-chan []byte
}

func NewWs2801Device(numLed int) *ws2801 {
	ws := new(ws2801)
	ws.device = spi.NewSPIDevice(spi.DEFAULT_BUS, spi.DEFAULT_CHIP, ws2801Delay)
	ws.numLed = numLed
	return ws
}

func (ws *ws2801) Open() error {
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

func (ws *ws2801) Release() error {
	return ws.device.Close()
}

func (ws *ws2801) Write(data []byte) error {
	if len(data) != ws.numLed*WS2801NumColor {
		return fmt.Errorf(
			"could not write %v bytes of data, %v is needed",
			len(data), WS2801NumColor*ws.numLed)
	}
	_, err := ws.device.Send(data)
	return err
}

func (ws *ws2801) SetInput(input <-chan []byte) {
	ws.input = input
}

func (ws *ws2801) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer ws.Release()
	for buffer := range ws.input {
		ws.Write(buffer)
	}
}

func (ws *ws2801) GetName() Name {
	return WS2801
}
