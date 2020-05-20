package main

import (
	"flag"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/config"
	"github.com/buttairfly/goPanel/internal/device"
	"github.com/buttairfly/goPanel/internal/generator"
	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/http"
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
	version.Init(compileDate, versionTag, logger)
	exitChan := make(chan bool)
	wg := new(sync.WaitGroup)
	go version.Run(wg, exitChan)

	mainConfigPath := flag.String("config", "config/main.composed.config.yaml", "path to config")
	flag.Parse()

	mainConfig, err1 := config.NewMainConfigFromPath(*mainConfigPath, logger)
	if err1 != nil {
		logger.Fatal("could not load mainConfig", zap.Error(err1))
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
	defer close(inputChan)

	pixelDevice.SetInput(inputChan)

	wg.Add(1)
	go pixelDevice.Run(wg)

	wg.Add(1)
	go generator.FrameGenerator(frame, inputChan, wg, logger)

	wg.Add(1)
	go http.RunHTTPServer(wg, logger)

	wg.Wait()
}
