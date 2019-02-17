package main

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/buttairfly/goPanel/src/device"
)

var (
	compileDate string
	versionTag  string
	programName string
)

func main() {
	const (
		panelLed   = 200
		bufferSize = panelLed * device.NumBytePerColor
	)

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile | log.LUTC)

	printProgramInfo()

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
				data[i] = byte(c)
			}
			inputChan <- data
			//pixelDevice.Write(data)
			//time.Sleep(time.Second)
		}
	}
}

func printProgramInfo() {
	program := strings.Split(os.Args[0], "/")
	programName = program[len(program)-1]
	go func() {
		for {
			log.Printf("%s: compiled at %s with version %s", programName, compileDate, versionTag)
			time.Sleep(30 * time.Second)
		}
	}()
}
