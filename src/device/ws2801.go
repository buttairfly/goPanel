package device

import (
	"fmt"

	"sync"

	"log"

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

func (ws *ws2801) Close() error {
	return ws.device.Close()
}

func (ws *ws2801) Write(data []byte) (int, error) {
	if len(data) != ws.numLed*WS2801NumColor {
		return 0, fmt.Errorf(
			"could not write %v bytes of data, %v is needed",
			len(data), WS2801NumColor*ws.numLed)
	}
	_, err := ws.device.Send(data)
	return len(data), err
}

func (ws *ws2801) SetInput(input <-chan []byte) {
	ws.input = input
}

func (ws *ws2801) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer ws.Close()
	for buffer := range ws.input {
		_, err := ws.Write(buffer)
		if err != nil {
			log.Panic(err)
		}
	}
}

func (ws *ws2801) GetType() Type {
	return WS2801
}
