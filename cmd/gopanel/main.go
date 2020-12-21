package main

import (
	"context"
	"flag"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/config"
	"github.com/buttairfly/goPanel/internal/device"
	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/http"
	"github.com/buttairfly/goPanel/internal/panel"
	"github.com/buttairfly/goPanel/pkg/exit"
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
	ctx := context.Background()
	cancelCtx := exit.DetectSignal(ctx, logger)
	exit.GracefulExit(cancelCtx, 4, 10*time.Second, 100*time.Millisecond, logger)

	goVersion := version.New("golang", "goVersion", compileDate, runtime.Version(), 0, logger)
	goVersion.Log()
	gracePeriod := 5 * time.Second
	mainVersion := version.New("main", version.GetProgramName(), compileDate, versionTag, gracePeriod, logger)
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
	panel := panel.NewPanel(cancelCtx, mainConfig, pixelDevice, logger)

	wg.Add(1)
	go pixelDevice.Run(cancelCtx, wg)

	wg.Add(1)
	go panel.Run(cancelCtx, wg)

	wg.Add(1)
	go http.RunHTTPServer(cancelCtx, wg, gracePeriod-time.Second, logger)

	wg.Wait()

	time.Sleep(2 * gracePeriod)
	logger.Panic("main not gracefully exited")
}
