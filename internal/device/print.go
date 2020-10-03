package device

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/leakybuffer"
)

type printDevice struct {
	inputChan    hardware.FrameSource
	currentFrame hardware.Frame
	numPix       int
	lenHex       int
	printConfig  *PrintConfig
	logger       *zap.Logger
}

// NewPrintDevice creates a new printDevice as LedDevice
func NewPrintDevice(numPix int, printConfig *PrintConfig, logger *zap.Logger) LedDevice {
	pd := new(printDevice)
	pd.numPix = numPix
	pd.lenHex = numPix * NumBytePerColor * NumByteToRepresentHex
	pd.printConfig = printConfig
	pd.logger = logger
	return pd
}

func (pd *printDevice) Open() error {
	pd.logger.Info("open print device")
	return nil
}

func (pd *printDevice) Close() error {
	return nil
}

func (pd *printDevice) Write(data string) (int, error) {
	if !pd.printConfig.Quiet {
		hexData := fmt.Sprintf("%x", data)
		if len(hexData) != pd.lenHex {
			return 0, fmt.Errorf(
				"len write hexData %v/numBytePerColor=%v/numByteToRepresentHex=%v does not equal numPix %v",
				len(hexData),
				pd.numPix,
				NumBytePerColor,
				NumByteToRepresentHex,
			)
		}
		pd.logger.Info("printDevice", zap.String("frame", hexData))
	}
	return len(data), nil
}

func (pd *printDevice) SetInput(inputChan hardware.FrameSource) {
	pd.logger.Debug("printDevice SetInput")
	pd.inputChan = inputChan
}

func (pd *printDevice) Run(cancelCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer pd.Close()
	frameDuration := time.Second / time.Duration(pd.printConfig.FramesPerSecond)
	lastFrameTime := time.Unix(0, 0)
	for frame := range pd.inputChan {
		now := time.Now()
		pd.currentFrame = frame
		// TODO: fix frame input
		sleepDuration := frameDuration - now.Sub(lastFrameTime)
		pd.logger.Sugar().Infof("sleepDuration %d, %v, %v", runtime.NumGoroutine(), sleepDuration, lastFrameTime)

		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}
		if _, err := pd.Write(string(frame.ToLedStripe().GetBuffer())); err != nil {
			pd.logger.Fatal("write error", zap.Error(err))
		}
		leakybuffer.DumpFrame(frame)
		lastFrameTime = now
	}
}

func (pd *printDevice) GetType() Type {
	return Print
}

func (pd *printDevice) GetCurrentFrame() hardware.Frame {
	return pd.currentFrame
}
