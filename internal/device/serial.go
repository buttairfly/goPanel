package device

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/buttairfly/goPanel/internal/config"
	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/pkg/com"
)

type serialDevice struct {
	serialCom  com.serialCom
	numLed     int
	input      <-chan hardware.Frame
	latchEnd   chan bool
	latched    int64
}

// NewSerialDevice creates a new serial device
func NewSerialDevice(numLed int, serialDeviceConfig *config.SerialConfig) LedDevice {
	s := new(serialDevice)
	s.config = serialDeviceConfig
	return s
}

func (s *serialDevice) Open() error {
	var err error
	s.stream, err = serial.OpenPort(s.serialCom.config.StreamConfig.ToStreamSerialConfig())
	return err
}

func (s *serialDevice) Close() error {
	return s.serialCom.stream.Close()
}

func (s *serialDevice) Write(data []byte) (int, error) {
	if s.serialCom.config.Verbose {
		log.Printf("Command %s", string(data))
	}
	n, err := s.serialCom.stream.Write(data)
	time.Sleep(s.config.CommandSleepTime)
	if err != nil {
		log.Fatal(err)
	}
	return n, err
}

func (s *serialDevice) SetInput(input <-chan hardware.Frame) {
	s.input = input
}

func (s *serialDevice) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer s.Close()
	defer close(s.readActive)
	defer close(s.stats)
	defer close(s.latchEnd)

	s.readActive = make(chan bool)
	s.initDone = make(chan bool)
	s.latchEnd = make(chan bool)
	s.stats = make(chan *stats, 10)
	wg.Add(1)
	go s.read(wg)
	wg.Add(1)
	go s.printStats(wg)
	wg.Add(1)
	go s.printLatches(wg)

	latchDelay := s.config.LatchSleepTime
	lastFrameTime := time.Now().Add(-latchDelay)

	// initialize bitbanger with number of leds
	time.Sleep(latchDelay)
	s.init()

	lastLedStripe := hardware.NewLedStripe(s.numLed)
	for frame := range s.input {
		ledStripe := frame.ToLedStripe()
		ledStripeCompare := ledStripe.Compare(lastLedStripe)
		if ledStripeCompare.HasChanged() {
			now := time.Now()
			sleepDuration := latchDelay - (now.Sub(lastFrameTime))
			s.stats <- &stats{
				event:     latchType,
				timeStamp: now,
				message:   fmt.Sprintf("%v", sleepDuration),
			}
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

func (s *serialDevice) GetType() config.Type {
	return config.Serial
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
			log.Printf("Latched frames: %f.2/s last lap: %d last diff: %v",
				float64(s.latched)*float64(time.Second)/float64(timeDiff),
				s.latched-lastLapLatches,
				timeDiff,
			)
			lastLapLatches = s.latched
		}
	}
}

func (s *serialDevice) setPixel(pixel int, buffer []uint8) {
	bufIndex := pixel * NumBytePerColor
	command := fmt.Sprintf("P%04x%02x%02x%02x\n", pixel, buffer[bufIndex+0], buffer[bufIndex+1], buffer[bufIndex+2])
	s.Write([]byte(command))
}

func (s *serialDevice) shade(pixel int, buffer []uint8) {
	command := fmt.Sprintf("S%04x%02x%02x%02x\n", pixel, buffer[0], buffer[1], buffer[2])
	s.Write([]byte(command))
}

func (s *serialDevice) latchFrame() {
	command := "L\n"
	s.Write([]byte(command))
}

