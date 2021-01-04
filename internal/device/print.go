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
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

type printDevice struct {
	inputChan    hardware.FrameSource
	currentFrame hardware.Frame
	params       []pipepart.PipeParam
	prevID       pipepart.ID
	numPix       int
	lenHex       int
	printConfig  *PrintConfig
	logger       *zap.Logger
}

// NewPrintDevice creates a new printDevice as LedDevice
func NewPrintDevice(numPix int, printConfig *PrintConfig, logger *zap.Logger) LedDevice {
	me := new(printDevice)
	me.numPix = numPix
	me.lenHex = numPix * NumBytePerColor * NumByteToRepresentHex
	me.printConfig = printConfig
	params := make([]pipepart.PipeParam, 1)
	params[0] = pipepart.PipeParam{
		Name:     "type",
		Type:     pipepart.NameID,
		Value:    string(me.GetDeviceType()),
		Readonly: true,
	}
	me.params = params
	me.logger = logger
	return me
}

func (me *printDevice) Open() error {
	me.logger.Info("open print device")
	return nil
}

func (me *printDevice) Close() error {
	return nil
}

func (me *printDevice) Write(data string) (int, error) {
	if !me.printConfig.Quiet {
		hexData := fmt.Sprintf("%x", data)
		if len(hexData) != me.lenHex {
			return 0, fmt.Errorf(
				"len write hexData %v/numBytePerColor=%v/numByteToRepresentHex=%v does not equal numPix %v",
				len(hexData),
				me.numPix,
				NumBytePerColor,
				NumByteToRepresentHex,
			)
		}
		me.logger.Info("printDevice", zap.String("frame", hexData))
	}
	return len(data), nil
}

func (me *printDevice) SetInput(prevID pipepart.ID, inputChan hardware.FrameSource) {
	me.logger.Debug("printDevice SetInput")
	me.inputChan = inputChan
	me.prevID = prevID
}

func (me *printDevice) RunPipe(cancelCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer me.Close()
	frameDuration := time.Second / time.Duration(me.printConfig.FramesPerSecond)
	lastFrameTime := time.Unix(0, 0)
	for frame := range me.inputChan {
		now := time.Now()
		me.currentFrame = frame
		// TODO: fix frame input
		sleepDuration := frameDuration - now.Sub(lastFrameTime)
		me.logger.Sugar().Infof("sleepDuration %d, %v, %v", runtime.NumGoroutine(), sleepDuration, lastFrameTime)

		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}
		if _, err := me.Write(string(frame.ToLedStripe().GetBuffer())); err != nil {
			me.logger.Fatal("write error", zap.Error(err))
		}
		leakybuffer.DumpFrame(frame)
		lastFrameTime = now
	}
}

func (me *printDevice) GetID() pipepart.ID {
	return pipepart.SinkID
}

func (me *printDevice) GetPrevID() pipepart.ID {
	return me.prevID
}

func (me *printDevice) Marshal() *pipepart.Marshal {
	return pipepart.MarshalFromPixelPiperSinkInterface(me)
}

func (me *printDevice) GetParams() []pipepart.PipeParam {
	return me.params
}

func (me *printDevice) GetDeviceType() Type {
	return Print
}

func (me *printDevice) GetType() pipepart.PipeType {
	return pipepart.Sink
}

func (me *printDevice) GetCurrentFrame() hardware.Frame {
	return me.currentFrame
}
