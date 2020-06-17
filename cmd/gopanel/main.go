package main

import (
	"context"
	"flag"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/config"
	"github.com/buttairfly/goPanel/internal/device"
	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/http"
	"github.com/buttairfly/goPanel/internal/panel"
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

	mainVersion := version.New("main", version.GetProgramName(), compileDate, versionTag, 10*time.Second, logger)
	go mainVersion.Run(cancelCtx)

	mainConfigPath := flag.String("config", "config/main.composed.config.yaml", "path to config")
	flag.Parse()

	mainConfig, err1 := config.NewMainConfigFromPath(*mainConfigPath, logger)
	if err1 != nil {
		logger.Fatal("could not load mainConfig", zap.Error(err1))
	}

	frame := hardware.NewFrame(mainConfig.TileConfigs.ToTileConfigs(), logger)

	pixelDevice, err := device.NewLedDevice(
		mainConfig.LedDeviceConfig,
		frame.GetSumHardwarePixel(),
		logger,
	)
	if err != nil {
		logger.Fatal("could not load pixelDevice", zap.Error(err))
	}
	defer pixelDevice.Close()

	wg := new(sync.WaitGroup)
	panel.NewPanel(cancelCtx, mainConfig, pixelDevice, wg, logger)

	inputChan := make(chan hardware.Frame)
	// inputChan is closed in LastBlackFrameGenerator

	pixelDevice.SetInput(inputChan)

	wg.Add(1)
	go pixelDevice.Run(cancelCtx, wg)

	go http.RunHTTPServer(logger)

	wg.Wait()

	logger.Info("successfully stopped")
}
