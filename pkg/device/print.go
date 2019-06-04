package device

import (
	"fmt"
	"log"
	"sync"

	"github.com/buttairfly/goPanel/internal/config"
	"github.com/buttairfly/goPanel/internal/hardware"
)

type printDevice struct {
	input  <-chan hardware.Frame
	numPix int
}

// NewPrintDevice creates a new printDevice
func NewPrintDevice(numPix int) LedDevice {
	pd := new(printDevice)
	pd.numPix = numPix
	return pd
}

func (pd *printDevice) Open() error {
	log.Print("Open print device")
	return nil
}

func (pd *printDevice) Close() error {
	return nil
}

func (pd *printDevice) Write(data []byte) (int, error) {
	log.Printf("%+x", data)
	if len(data) != pd.numPix {
		return 0, fmt.Errorf(
			"len write data %v does not equal numPix %v", len(data), pd.numPix)
	}
	return len(data), nil
}

func (pd *printDevice) SetInput(input <-chan hardware.Frame) {
	pd.input = input
}

func (pd *printDevice) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer pd.Close()
	for frame := range pd.input {
		_, err := pd.Write(([]byte)(frame.ToLedStripe().GetBuffer()))
		if err != nil {
			log.Panic(err)
		}
	}
}

func (pd *printDevice) GetType() config.Type {
	return config.Print
}
