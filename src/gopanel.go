package main

import (
	"log"
	"sync"

	"github.com/buttairfly/goPanel/src/device"
)

func main() {
	const (
		panelLed   = 200
		bufferSize = panelLed * device.NumBytePerColor
	)
	pixelDevice, err := device.NewLedDevice(device.Serial, panelLed)
	//pixelDevice, err := device.NewLedDevice(device.WS2801, panelLed)
	//pixelDevice, err := device.NewLedDevice(device.Print, panelLed)
	if err != nil {
		log.Fatal(err)
	}
	defer pixelDevice.Close()

	inputChan := make(chan []byte)
	pixelDevice.SetInput(inputChan)
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go pixelDevice.Run(wg)
	for {
		for c := 0; c < 0x100; c++ {
			data := make([]byte, bufferSize, bufferSize)
			for i := range data {
				data[i] = byte(i + c)
			}
			inputChan <- data
			//pixelDevice.Write(data)
			//time.Sleep(time.Second)
		}
	}
}
