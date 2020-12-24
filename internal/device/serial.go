package device

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/leakybuffer"
	"github.com/buttairfly/goPanel/pkg/arduinocom"
)

type serialDevice struct {
	com          *arduinocom.ArduinoCom
	numLed       int
	inputChan    hardware.FrameSource
	currentFrame hardware.Frame
	latched      int64
	logger       *zap.Logger
}

// NewSerialDevice creates a new serial device
func NewSerialDevice(numLed int, serialDeviceConfig *arduinocom.SerialConfig, logger *zap.Logger) LedDevice {
	s := new(serialDevice)
	s.com = arduinocom.NewArduinoCom(numLed, serialDeviceConfig, logger)
	s.numLed = numLed
	s.logger = logger
	return s
}

func (s *serialDevice) Open() error {
	return s.com.Open()
}

func (s *serialDevice) Close() error {
	return s.com.Close()
}

func (s *serialDevice) Write(command string) (int, error) {
	return s.com.CalcParityAndWrite(command)
}

func (s *serialDevice) SetInput(inputChan hardware.FrameSource) {
	s.inputChan = inputChan
}

func (s *serialDevice) Run(cancelCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer s.Close()

	subWg := new(sync.WaitGroup)
	subWg.Add(4)
	go s.com.Read(cancelCtx, subWg)
	go s.com.PrintStats(cancelCtx, subWg)
	go s.printLatches(cancelCtx, subWg)
	go s.runFrameProcessor(subWg)

	subWg.Wait()
}

func (s *serialDevice) runFrameProcessor(wg *sync.WaitGroup) {
	defer wg.Done()

	latchDelay := s.com.Config().LatchSleepTime
	lastFrameTime := time.Now().Add(-latchDelay)

	// initialize bitbanger with number of leds
	time.Sleep(latchDelay)
	s.com.Init()

	s.logger.Sugar().Infof("numLed ledStripe %d", s.numLed)
	for frame := range s.inputChan {
		s.currentFrame = frame // this is unsafe
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
			s.com.AddStat(stat)
			if sleepDuration > 0 {
				time.Sleep(sleepDuration)
			}
			lastFrameTime = now

			if ledStripeAction.IsFullFrame() {
				s.rawFrame(s.numLed, ledStripe.GetBuffer())
			} else {
				fillColor := ledStripeAction.GetFillColor()
				if fillColor != nil {
					s.shade(s.numLed, fillColor.Slice())
				}
				for _, pixelIndex := range ledStripeAction.GetOtherDiffPixels() {
					s.setPixel(pixelIndex, ledStripe.GetBuffer())
				}
				s.latchFrame()
			}
		}
		leakybuffer.DumpFrame(frame)
	}
}

func (s *serialDevice) GetType() Type {
	return Serial
}

func (s *serialDevice) GetCurrentFrame() hardware.Frame {
	return s.currentFrame // this is unsafe
}

func (s *serialDevice) printLatches(cancelCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer s.logger.Info("printLatches done")

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
				currentLapLatched := s.latched - lastLapLatches
				s.logger.Info("latch summary",
					zap.String("frames", fmt.Sprintf("%d (%2.3f/s)", s.latched, float64(s.latched)/float64(timeDiff.Seconds()))),
					zap.String("lap frames", fmt.Sprintf("%d (%2.3f/s)", currentLapLatched, float64(currentLapLatched)/float64(timerDuration.Seconds()))),
					zap.Duration("diff", timeDiff),
				)
				lastLapLatches = s.latched
			}
		}
	}
}

func (s *serialDevice) setPixel(pixelNum int, buffer []uint8) {
	bufIndex := pixelNum * NumBytePerColor
	command := fmt.Sprintf("P%04x%02x%02x%02x", pixelNum, buffer[bufIndex+0], buffer[bufIndex+1], buffer[bufIndex+2])
	s.Write(command)
	time.Sleep(s.com.Config().LatchSleepTime)
}

func (s *serialDevice) shade(pixel int, buffer []uint8) {
	command := fmt.Sprintf("S%04x%02x%02x%02x", pixel, buffer[0], buffer[1], buffer[2])
	s.Write(command)
	time.Sleep(s.com.Config().LatchSleepTime)
}

func (s *serialDevice) rawFrame(pixel int, frameBuffer []uint8) {
	currentRawFramePartNumLed := s.com.Config().RawFramePartNumLed
	if currentRawFramePartNumLed == 0 {
		currentRawFramePartNumLed = pixel
	}
	maxRawFrameParts := pixel / currentRawFramePartNumLed
	for currentRawFramePart := 0; currentRawFramePart < maxRawFrameParts; currentRawFramePart++ {
		pixelOffset := currentRawFramePart * currentRawFramePartNumLed
		s.rawFramePart(pixel, pixelOffset, currentRawFramePart, currentRawFramePartNumLed, frameBuffer)
	}
	remainingRawFramePartNumLed := pixel % currentRawFramePartNumLed
	if remainingRawFramePartNumLed > 0 {
		pixelOffset := maxRawFrameParts * currentRawFramePartNumLed
		s.rawFramePart(pixel, pixelOffset, maxRawFrameParts, remainingRawFramePartNumLed, frameBuffer)
	}
	time.Sleep(s.com.Config().LatchSleepTime)
	s.latchFrame()
}

func (s *serialDevice) rawFramePart(pixel, pixelOffset, currentRawFramePart, currentRawFramePartNumLed int, frameBuffer []uint8) {
	frameString := ""
	for p := 0; p < currentRawFramePartNumLed; p++ {
		bufIndex := (pixelOffset + p) * NumBytePerColor
		color := fmt.Sprintf("%02x%02x%02x", frameBuffer[bufIndex], frameBuffer[bufIndex+1], frameBuffer[bufIndex+2])
		frameString += color
	}
	command := fmt.Sprintf("W%04x%02x%02x%s", pixel, currentRawFramePart, currentRawFramePartNumLed, frameString)
	time.Sleep(s.com.Config().LatchSleepTime)
	s.Write(command)
}

func (s *serialDevice) latchFrame() {
	s.latched++
	command := "L"
	s.Write(command)
}
