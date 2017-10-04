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

func (pd *printDevice) Close() error {
	return nil
}

func (pd *printDevice) Write(data []byte) (int, error) {
	log.Printf("%+x", data)
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
