package main

import (
	"context"
	"flag"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/config"
	"github.com/buttairfly/goPanel/internal/device"
	"github.com/buttairfly/goPanel/internal/generator"
	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/http"
	"github.com/buttairfly/goPanel/pkg/log"
	"github.com/buttairfly/goPanel/pkg/routine"
	"github.com/buttairfly/goPanel/pkg/version"
)

var (
	compileDate string
	versionTag  string
)

func main() {
	logger := log.NewZapDevelopLogger()
	defer logger.Sync()
	ctx := context.Background()
	cancelCtx := routine.DetectExit(ctx)

	mainVersion := version.New(compileDate, versionTag, 10*time.Second, logger)
	go mainVersion.Run(cancelCtx)

	mainConfigPath := flag.String("config", "config/main.composed.config.yaml", "path to config")
	flag.Parse()

	mainConfig, err1 := config.NewMainConfigFromPath(*mainConfigPath, logger)
	if err1 != nil {
		logger.Fatal("could not load mainConfig", zap.Error(err1))
	}

	frame := hardware.NewFrame(mainConfig.TileConfigs, logger)

	pixelDevice, err := device.NewLedDevice(
		cancelCtx,
		mainConfig.LedDeviceConfig,
		frame.GetSumHardwarePixel(),
		logger,
	)
	if err != nil {
		logger.Fatal("could not load pixelDevice", zap.Error(err))
	}
	defer pixelDevice.Close()

	inputChan := make(chan hardware.Frame)
	// inputChan is closed in LastBlackFrameFrameGenerator

	pixelDevice.SetInput(inputChan)

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go pixelDevice.Run(wg)

	wg.Add(1)
	go generator.LastBlackFrameFrameGenerator(cancelCtx, frame, inputChan, wg, logger)

	wg.Add(1)
	go generator.FrameGenerator(cancelCtx, frame, inputChan, wg, logger)

	go http.RunHTTPServer(logger)

	wg.Wait()

	logger.Info("successfully stopped")
}
