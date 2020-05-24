package device

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
)

type printDevice struct {
	inputChan   <-chan hardware.Frame
	numPix      int
	lenHex      int
	printConfig *PrintConfig
	cancelCtx   context.Context
	logger      *zap.Logger
}

// NewPrintDevice creates a new printDevice
func NewPrintDevice(cancelCtx context.Context, numPix int, printConfig *PrintConfig, logger *zap.Logger) LedDevice {
	pd := new(printDevice)
	pd.numPix = numPix
	pd.lenHex = numPix * NumBytePerColor * NumByteToRepresentHex
	pd.cancelCtx = cancelCtx
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
	return len(data), nil
}

func (pd *printDevice) SetInput(inputChan <-chan hardware.Frame) {
	pd.inputChan = inputChan
}

func (pd *printDevice) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer pd.Close()
	frameDuration := time.Second / time.Duration(pd.printConfig.FramesPerSecond)
	for frame := range pd.inputChan {

		// TODO: fix frame input
		// pd.logger.Info("receive frame", zap.Time("frameTime", frame.GetTime()))
		now := time.Now()
		sleepDuration := frameDuration - now.Sub(frame.GetTime())
		// pd.logger.Sugar().Infof("sleepDuration %d, %v, %v", runtime.NumGoroutine(), sleepDuration, frame.GetTime())

		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}
		if _, err := pd.Write(string(frame.ToLedStripe().GetBuffer())); err != nil {
			pd.logger.Fatal("write error", zap.Error(err))
		}
	}
}

func (pd *printDevice) GetType() Type {
	return Print
}
