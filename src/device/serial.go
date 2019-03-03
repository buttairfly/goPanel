package device

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
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

func (s *serialDevice) init() {
	defer func() {
		if s.config.Verbose {
			s.sendInitComand("Q00ff\n")
		}
	}()

	for {
		select {
		case <-s.initDone:
			return
		default:
			s.stream.Flush()
			s.sendInitComand("Q0000\n")
			s.sendInitComand("V\n")
			s.sendInitComand(fmt.Sprintf("I%04x\n", s.numLed))
		}
	}
}

func (s *serialDevice) sendInitComand(command string) {
	s.Write([]byte(command))
	time.Sleep(s.config.InitSleepTime)
}

func (s *serialDevice) read(wg *sync.WaitGroup) {
	lastLine := ""
	defer wg.Done()
	defer log.Println(lastLine)

	buf := make([]byte, s.config.ReadBufferSize)
	for {
		select {
		case <-s.readActive:
			return
		default:
			n, err := s.stream.Read(buf)
			if err != nil {
				if err == io.EOF {
					continue
				}
				log.Fatal(err)
			}
			read := lastLine + string(buf[:n])
			lines := strings.Split(read, "\n")
			if read[len(read)-1] != '\n' {
				numLines := len(lines) - 1
				lastLine = lines[numLines]
				lines = lines[:numLines]
			} else {
				lastLine = ""
			}
			for _, line := range lines {
				if len(line) > 0 {
					log.Println(line)
					if s.needsInit() {
						parts := strings.Split(line, " ")
						if len(parts) == 2 && parts[0] == "Init" {
							initLed, err := strconv.ParseInt(parts[1], 16, 0))
							if err != nil {
								log.Print( err)
							} else {
								if initLed == s.numLed {
									close(s.initDone)
								}
							}
						}
					}
				}
			}
		}
	}
}

func (s *serialDevice) needsInit() bool {
	select {
	case <-s.initDone:
		return false
	default:
		return true
	}
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

func (s *serialDevice) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer s.Close()
	defer close(s.readActive)

	s.readActive = make(chan bool)
	s.initDone = make(chan bool)
	wg.Add(1)
	go s.read(wg)

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
