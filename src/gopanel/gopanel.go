package main

import (
	"log"

	"github.com/buttairfly/goPanel/src/device"
)

func main() {
	const (
		panelLed   = 200
		bufferSize = panelLed * device.WS2801NumColor
	)
	pixelDevice, err := device.NewSpiDevice(device.Print, panelLed)
	if err != nil {
		log.Fatal(err)
	}
	defer pixelDevice.Close()

	for {
		for c := 0; c < 0x100; c++ {
			data := make([]byte, bufferSize, bufferSize)
			for i := range data {
				data[i] = byte(c)
			}
			pixelDevice.Write(data)
		}
	}
}
