package device

import (
	"fmt"
	"log"
	"sync"
)

type printDevice struct {
	input  <-chan []byte
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

func (pd *printDevice) SetInput(input <-chan []byte) {
	pd.input = input
}

func (pd *printDevice) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer pd.Close()
	for buffer := range pd.input {
		_, err := pd.Write(buffer)
		if err != nil {
			log.Panic(err)
		}
	}
}

func (pd *printDevice) GetType() Type {
	return Print
}
