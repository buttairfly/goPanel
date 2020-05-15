package main

import (
	"flag"
	"image"
	"image/color"
	"sync"

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
	sugar := logger.Sugar()

	version.PrintProgramInfo(compileDate, versionTag, logger)

	panelConfigPtr := flag.String("config", "config/main.panel.config.yaml", "path to config")

	flag.Parse()
	mainConfig, err1 := config.NewMainConfigFromPanelConfigPath(*panelConfigPtr, logger)
	if err1 != nil {
		sugar.Fatalf("Could not load mainConfig %e", err1)
	}

	frame := hardware.NewFrame(mainConfig.TileConfigs, logger)

	pixelDevice, err := device.NewLedDevice(
		mainConfig.DeviceConfig,
		frame.GetSumHardwarePixel(),
		logger,
	)
	if err != nil {
		sugar.Fatalf("Could not load pixelDevice %e", err)
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
			colorFrame := hardware.NewCopyFrameFromImage(frame, mainPicture, logger)
			inputChan <- colorFrame
		}
	}

}
