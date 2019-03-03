package main

import (
	"flag"
	"image"
	"image/color"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/buttairfly/goPanel/src/config"
	"github.com/buttairfly/goPanel/src/device"
	"github.com/buttairfly/goPanel/src/hardware"
	"github.com/buttairfly/goPanel/src/palette"
)

var (
	compileDate string
	versionTag  string
	programName string
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile | log.LUTC)

	printProgramInfo()

	panelConfigPtr := flag.String("config", "main.panel.config.json", "a string")
	folderConfigPtr := flag.String("folder", "config/", "a string")
	flag.Parse()
	mainConfig, err1 := config.NewConfigFromPanelConfigPath(*folderConfigPtr, *panelConfigPtr)
	if err1 != nil {
		log.Fatal(err1)
	}

	frame := hardware.NewFrame(mainConfig.GetTileConfigs())

	pixelDevice, err := device.NewLedDevice(
		mainConfig.GetDeviceConfig(),
		frame.GetSumHardwarePixel(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer pixelDevice.Close()

	inputChan := make(chan hardware.Frame)
	pixelDevice.SetInput(inputChan)
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go pixelDevice.Run(wg)

	mainPicture := image.NewRGBA(frame.Bounds())

	colors := make([]color.Color, 0, 10)
	colors = append(colors, color.RGBA{0xff, 0, 0, 0xff})
	colors = append(colors, color.RGBA{0xff, 0xa5, 0, 0xff})
	const granularity int = 200
	const wrapping bool = true
	fader := palette.NewFader(colors, granularity, wrapping)
	increments := fader.GetIncrements()
	for {
		for _, increment := range increments {
			color := fader.Fade(increment)
			for y := 0; y < frame.GetHeight(); y++ {
				for x := 0; x < frame.GetWidth(); x++ {
					mainPicture.Set(x, y, color)
				}
			}
			colorFrame := hardware.NewCopyFrameFromImage(frame, mainPicture)
			inputChan <- colorFrame
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
