package main

import (
	"flag"
	"image"
	"image/color"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/config"
	"github.com/buttairfly/goPanel/internal/device"
	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/palette"
	"github.com/buttairfly/goPanel/pkg/log"
	"github.com/buttairfly/goPanel/pkg/version"
)

var (
	compileDate string
	versionTag  string
)

func main() {
	logger := log.NewZapDevelopLogger()
	defer logger.Sync()

	version.PrintProgramInfo(compileDate, versionTag, logger)

	mainConfigPath := *(flag.String("config", "config/main.composed.config.yaml", "path to config"))

	flag.Parse()
	mainConfig, err1 := config.NewMainConfigFromPath(mainConfigPath, logger)
	if err1 != nil {
		logger.Fatal("could not load mainConfig %e", zap.Error(err1))
	}

	frame := hardware.NewFrame(mainConfig.TileConfigs, logger)

	pixelDevice, err := device.NewLedDevice(
		mainConfig.LedDeviceConfig,
		frame.GetSumHardwarePixel(),
		logger,
	)
	if err != nil {
		logger.Fatal("could not load pixelDevice", zap.Error(err))
	}
	defer pixelDevice.Close()

	inputChan := make(chan hardware.Frame)
	pixelDevice.SetInput(inputChan)
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go pixelDevice.Run(wg)

	wg.Add(1)
	go FrameGenerator(frame, inputChan, wg, logger)
}

func FrameGenerator(frame hardware.Frame, inputChan chan<- hardware.Frame, wg *sync.WaitGroup, logger *zap.Logger) {
	wg.Add(1)
	defer wg.Done()

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
			colorFrame := hardware.NewCopyFrameFromImage(frame, mainPicture, logger)
			inputChan <- colorFrame
		}
	}
}
