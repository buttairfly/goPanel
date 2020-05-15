package device

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
)

type printDevice struct {
	input  <-chan hardware.Frame
	numPix int
	logger *zap.Logger
}

// NewPrintDevice creates a new printDevice
func NewPrintDevice(numPix int, logger *zap.Logger) LedDevice {
	pd := new(printDevice)
	pd.numPix = numPix
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
	pd.logger.Sugar().Infof("%+x", data)
	if len(data) != pd.numPix {
		return 0, fmt.Errorf(
			"len write data %v does not equal numPix %v", len(data), pd.numPix)
	}
	return len(data), nil
}

func (pd *printDevice) SetInput(input <-chan hardware.Frame) {
	pd.input = input
}

func (pd *printDevice) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer pd.Close()
	for frame := range pd.input {
		_, err := pd.Write(string(frame.ToLedStripe().GetBuffer()))
		if err != nil {
			pd.logger.Panic("write error", zap.Error(err))
		}
	}
}

func (pd *printDevice) GetType() Type {
	return Print
}
