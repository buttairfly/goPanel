package device

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/pkg/arduinocom"
)

type serialDevice struct {
	com      *arduinocom.ArduinoCom
	numLed   int
	input    <-chan hardware.Frame
	latchEnd chan bool
	latched  int64
	logger   *zap.Logger
}

// NewSerialDevice creates a new serial device
func NewSerialDevice(numLed int, serialDeviceConfig *arduinocom.SerialConfig, logger *zap.Logger) LedDevice {
	s := new(serialDevice)
	s.com = arduinocom.NewArduinoCom(numLed, serialDeviceConfig, logger)
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

func (s *serialDevice) SetInput(input <-chan hardware.Frame) {
	s.input = input
}

func (s *serialDevice) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer s.Close()
	defer close(s.latchEnd)

	s.latchEnd = make(chan bool)
	wg.Add(1)
	go s.com.Read(wg)
	wg.Add(1)
	go s.com.PrintStats(wg)
	wg.Add(1)
	go s.printLatches(wg)

	latchDelay := s.com.Config().LatchSleepTime
	lastFrameTime := time.Now().Add(-latchDelay)

	// initialize bitbanger with number of leds
	time.Sleep(latchDelay)
	s.com.Init()

	lastLedStripe := hardware.NewLedStripe(s.numLed, s.logger)
	for frame := range s.input {
		ledStripe := frame.ToLedStripe()
		ledStripeCompare := ledStripe.Compare(lastLedStripe)
		if ledStripeCompare.HasChanged() {
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
			fullColor := ledStripeCompare.GetFullColor()
			if fullColor != nil {
				s.shade(s.numLed, fullColor.Slice())
			}
			for _, pixelIndex := range ledStripeCompare.GetOtherDiffPixels() {
				s.setPixel(pixelIndex, ledStripe.GetBuffer())
			}

			s.latchFrame()
			lastLedStripe = ledStripe
		}
	}
}

func (s *serialDevice) GetType() Type {
	return Serial
}

func (s *serialDevice) printLatches(wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	lastLapLatches := int64(0)
	for {
		select {
		case <-s.latchEnd:
			return
		default:
			time.Sleep(30 * time.Second)
			timeDiff := time.Now().Sub(start)
			s.logger.Info("latch summary",
				zap.String("frames", fmt.Sprintf("%f.2/s", float64(s.latched)*float64(time.Second)/float64(timeDiff))),
				zap.Int64("lastLap", s.latched-lastLapLatches),
				zap.Duration("lastDiff", timeDiff),
			)
			lastLapLatches = s.latched
		}
	}
}

func (s *serialDevice) setPixel(pixel int, buffer []uint8) {
	bufIndex := pixel * NumBytePerColor
	command := fmt.Sprintf("P%04x%02x%02x%02x", pixel, buffer[bufIndex+0], buffer[bufIndex+1], buffer[bufIndex+2])
	s.Write(command)
}

func (s *serialDevice) shade(pixel int, buffer []uint8) {
	command := fmt.Sprintf("S%04x%02x%02x%02x", pixel, buffer[0], buffer[1], buffer[2])
	s.Write(command)
}

func (s *serialDevice) latchFrame() {
	command := "L"
	s.Write(command)
}
