package device

import (
	"log"
	"sync"
	"time"

	"github.com/buttairfly/goPanel/src/config"
	"github.com/buttairfly/goPanel/src/hardware"
	"github.com/tarm/serial"
)

type serialDevice struct {
	config     *config.SerialConfig
	stream     *serial.Port
	numLed     int
	input      <-chan hardware.Frame
	readActive chan bool
	initDone   chan bool
	stats      chan stats
}

// NewSerialDevice creates a new serial device
func NewSerialDevice(numLed int, serialDeviceConfig *config.SerialConfig) LedDevice {
	s := new(serialDevice)
	s.config = serialDeviceConfig
	s.numLed = numLed
	return s
}

func (s *serialDevice) Open() error {
	var err error
	s.stream, err = serial.OpenPort(s.config.StreamConfig.ToStreamSerialConfig())
	return err
}

func (s *serialDevice) Close() error {
	return s.stream.Close()
}

func (s *serialDevice) Write(data []byte) (int, error) {
	if s.config.Verbose {
		log.Printf("Command %s", string(data))
	}
	n, err := s.stream.Write(data)
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

	s.readActive = make(chan bool)
	s.initDone = make(chan bool)
	s.stats = make(chan stats)
	wg.Add(1)
	go s.read(wg)

	wg.Add(1)
	go s.printStats(wg)

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
			log.Println(sleepDuration, now.Sub(lastFrameTime))
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
