package main

import (
	"log"

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

	for {
		for c := 0; c < 0x100; c++ {
			data := make([]byte, bufferSize, bufferSize)
			for i := range data {
				data[i] = byte(i + c)
			}
			//pixelDevice.Write(data)
			//time.Sleep(time.Second)
		}
	}
}
