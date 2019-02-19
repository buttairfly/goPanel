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

	frame := make([]byte, bufferSize, bufferSize) // one buffer only as state
	for {
		for c := 0; c < device.NumBytePerColor; c++ {
			for p := 0; p < panelLed; p++ {
				frame = setStripPixelToColor(frame, p, rgbToColor(c%3 * 0xff, c%, b))

				data := make([]byte, bufferSize, bufferSize)
				copy(data, frame)
				inputChan <- data

			}
		}
	}
}

func rgbToColor(r, g, b byte) int {
	return (int(r)<<16 | int(g)<<8 | int(b)) & 0xffffff
}

func setStripPixelToColor(frame []byte, posOnStrip int, color int) []byte {
	pos := posOnStrip * device.NumBytePerColor
	frame[pos+0] = byte(color >> 16)
	frame[pos+1] = byte(color >> 8)
	frame[pos+2] = byte(color)
	return frame
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
