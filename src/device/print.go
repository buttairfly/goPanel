package device

import (
	"log"
	"sync"
)

type printDevice struct {
	input <-chan []byte
}

func NewPrintDevice() *printDevice {
	pd := new(printDevice)
	return pd
}

func (pd *printDevice) Open() error {
	log.Print("Open print device")
	return nil
}

func (pd *printDevice) Release() error {
	return nil
}

func (pd *printDevice) Write(data []byte) error {
	log.Printf("%+x", data)
	return nil
}

func (pd *printDevice) SetInput(input <-chan []byte) {
	pd.input = input
}

func (pd *printDevice) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer pd.Release()
	for buffer := range pd.input {
		pd.Write(buffer)
	}
}

func (pd *printDevice) GetName() Name {
	return PrintDevice
}
