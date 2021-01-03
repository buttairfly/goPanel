package device

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/leakybuffer"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
	"github.com/buttairfly/goPanel/pkg/arduinocom"
)

type serialDevice struct {
	com          *arduinocom.ArduinoCom
	numLed       int
	inputChan    hardware.FrameSource
	prevID       pipepart.ID
	currentFrame hardware.Frame
	latched      int64
	params       []pipepart.PipeParam
	logger       *zap.Logger
}

// NewSerialDevice creates a new serial device
func NewSerialDevice(numLed int, serialDeviceConfig *arduinocom.SerialConfig, logger *zap.Logger) LedDevice {
	me := new(serialDevice)
	me.com = arduinocom.NewArduinoCom(numLed, serialDeviceConfig, logger)
	me.numLed = numLed
	params := make([]pipepart.PipeParam, 1)
	params[0] = pipepart.PipeParam{
		Name:     "type",
		Type:     pipepart.NameID,
		Value:    string(me.GetType()),
		Readonly: true,
	}
	me.params = params
	me.logger = logger
	return me
}

func (me *serialDevice) Open() error {
	return me.com.Open()
}

func (me *serialDevice) Close() error {
	return me.com.Close()
}

func (me *serialDevice) Write(command string) (int, error) {
	return me.com.CalcParityAndWrite(command)
}

func (me *serialDevice) SetInput(prevID pipepart.ID, inputChan hardware.FrameSource) {
	me.inputChan = inputChan
	me.prevID = prevID
}

func (me *serialDevice) RunPipe(cancelCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer me.Close()

	subWg := new(sync.WaitGroup)
	subWg.Add(4)
	go me.com.Read(cancelCtx, subWg)
	go me.com.PrintStats(cancelCtx, subWg)
	go me.printLatches(cancelCtx, subWg)
	go me.runFrameProcessor(subWg)

	subWg.Wait()
}

func (me *serialDevice) runFrameProcessor(wg *sync.WaitGroup) {
	defer wg.Done()

	latchDelay := me.com.Config().LatchSleepTime
	lastFrameTime := time.Now().Add(-latchDelay)

	// initialize bitbanger with number of leds
	time.Sleep(latchDelay)
	me.com.Init()

	me.logger.Sugar().Infof("numLed ledStripe %d", me.numLed)
	for frame := range me.inputChan {
		//TODO: hot area
		if me.currentFrame != nil {
			leakybuffer.DumpFrame(me.currentFrame)
		}
		me.currentFrame = frame // this is unsafe
		// TODO: exit hot area

		ledStripe := frame.ToLedStripe()

		ledStripeAction := ledStripe.GetAction()
		if ledStripeAction.HasChanged() {
			now := time.Now()
			sleepDuration := latchDelay - (now.Sub(lastFrameTime))
			stat := &arduinocom.Stat{
				Event:     arduinocom.LatchStatType,
				TimeStamp: now,
				Message:   fmt.Sprintf("%v", sleepDuration),
			}
			me.com.AddStat(stat)
			if sleepDuration > 0 {
				time.Sleep(sleepDuration)
			}
			lastFrameTime = now

			if ledStripeAction.IsFullFrame() {
				me.rawFrame(me.numLed, ledStripe.GetBuffer())
			} else {
				fillColor := ledStripeAction.GetFillColor()
				if fillColor != nil {
					me.shade(me.numLed, fillColor.Slice())
				}
				for _, pixelIndex := range ledStripeAction.GetOtherDiffPixels() {
					me.setPixel(pixelIndex, ledStripe.GetBuffer())
				}
				me.latchFrame()
			}
		}
	}
}

func (me *serialDevice) GetType() Type {
	return Serial
}

func (me *serialDevice) GetID() pipepart.ID {
	return pipepart.SinkID
}

func (me *serialDevice) GetPrevID() pipepart.ID {
	return me.prevID
}

func (me *serialDevice) Marshal() pipepart.Marshal {
	return pipepart.Marshal{
		ID:     me.GetID(),
		PrevID: me.GetPrevID(),
		Params: me.GetParams(),
	}
}

func (me *serialDevice) GetParams() []pipepart.PipeParam {
	return me.params
}

func (me *serialDevice) GetCurrentFrame() hardware.Frame {
	return me.currentFrame // this is unsafe
}

func (me *serialDevice) printLatches(cancelCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer me.logger.Info("printLatches done")

	start := time.Now()
	lastLapLatches := int64(0)

	timerDuration := 30 * time.Second
	ticker := time.NewTicker(timerDuration)
	defer ticker.Stop()

	for {
		select {
		case <-cancelCtx.Done():
			return
		case now := <-ticker.C:
			{
				timeDiff := now.Sub(start)
				currentLapLatched := me.latched - lastLapLatches
				me.logger.Info("latch summary",
					zap.String("frames", fmt.Sprintf("%d (%2.3f/s)", me.latched, float64(me.latched)/float64(timeDiff.Seconds()))),
					zap.String("lap frames", fmt.Sprintf("%d (%2.3f/s)", currentLapLatched, float64(currentLapLatched)/float64(timerDuration.Seconds()))),
					zap.Duration("diff", timeDiff),
				)
				lastLapLatches = me.latched
			}
		}
	}
}

func (me *serialDevice) setPixel(pixelNum int, buffer []uint8) {
	bufIndex := pixelNum * NumBytePerColor
	command := fmt.Sprintf("P%04x%02x%02x%02x", pixelNum, buffer[bufIndex+0], buffer[bufIndex+1], buffer[bufIndex+2])
	me.Write(command)
	time.Sleep(me.com.Config().LatchSleepTime)
}

func (me *serialDevice) shade(pixel int, buffer []uint8) {
	command := fmt.Sprintf("S%04x%02x%02x%02x", pixel, buffer[0], buffer[1], buffer[2])
	me.Write(command)
	time.Sleep(me.com.Config().LatchSleepTime)
}

func (me *serialDevice) rawFrame(pixel int, frameBuffer []uint8) {
	currentRawFramePartNumLed := me.com.Config().RawFramePartNumLed
	if currentRawFramePartNumLed == 0 {
		currentRawFramePartNumLed = pixel
	}
	maxRawFrameParts := pixel / currentRawFramePartNumLed
	for currentRawFramePart := 0; currentRawFramePart < maxRawFrameParts; currentRawFramePart++ {
		pixelOffset := currentRawFramePart * currentRawFramePartNumLed
		me.rawFramePart(pixel, pixelOffset, currentRawFramePart, currentRawFramePartNumLed, frameBuffer)
	}
	remainingRawFramePartNumLed := pixel % currentRawFramePartNumLed
	if remainingRawFramePartNumLed > 0 {
		pixelOffset := maxRawFrameParts * currentRawFramePartNumLed
		me.rawFramePart(pixel, pixelOffset, maxRawFrameParts, remainingRawFramePartNumLed, frameBuffer)
	}
	time.Sleep(me.com.Config().LatchSleepTime)
	me.latchFrame()
}

func (me *serialDevice) rawFramePart(pixel, pixelOffset, currentRawFramePart, currentRawFramePartNumLed int, frameBuffer []uint8) {
	frameString := ""
	for p := 0; p < currentRawFramePartNumLed; p++ {
		bufIndex := (pixelOffset + p) * NumBytePerColor
		color := fmt.Sprintf("%02x%02x%02x", frameBuffer[bufIndex], frameBuffer[bufIndex+1], frameBuffer[bufIndex+2])
		frameString += color
	}
	command := fmt.Sprintf("W%04x%02x%02x%s", pixel, currentRawFramePart, currentRawFramePartNumLed, frameString)
	time.Sleep(me.com.Config().LatchSleepTime)
	me.Write(command)
}

func (me *serialDevice) latchFrame() {
	me.latched++
	command := "L"
	me.Write(command)
}
