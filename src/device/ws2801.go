package device

import (
	"fmt"

	"github.com/luismesas/goPi/spi"
)

const (
	ws2801_MODE  = uint8(0)
	ws2801_BPW   = uint8(8)
	ws2801_SPEED = uint32(1000000)
	ws2801_DELAY = uint16(550)

	WS2801_NUM_COLORS = 3
)

type ws2801 struct {
	device *spi.SPIDevice
	numLed int
}

func NewWs2801Device(devPath string, numLed int) *ws2801 {
	ws := new(ws2801)
	ws.device = spi.NewSPIDevice(spi.DEFAULT_BUS, spi.DEFAULT_CHIP, ws2801_DELAY)
	ws.numLed = numLed
	return ws
}

func (ws *ws2801) Open() error {
	err := ws.device.Open()
	if err != nil {
		return err
	}

	err = ws.device.SetMode(ws2801_MODE)
	if err != nil {
		return err
	}

	err = ws.device.SetBitsPerWord(ws2801_BPW)
	if err != nil {
		return err
	}

	err = ws.device.SetSpeed(ws2801_SPEED)
	if err != nil {
		return err
	}

	return nil
}

func (ws *ws2801) Release() error {
	return ws.device.Close()
}

func (ws *ws2801) Write(data []byte) error {
	if len(data) != ws.numLed*WS2801_NUM_COLORS {
		return fmt.Errorf(
			"could not write %v bytes of data, %v is needed", len(data), WS2801_NUM_COLORS*ws.numLed)
	}
	_, err := ws.device.Send(data)
	return err
}
