package device

import (
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/tarm/serial"
)

type serialDevice struct {
	config     *serial.Config
	stream     *serial.Port
	numLed     int
	input      <-chan []byte
	readActive chan bool
}

// NewSerialDevice creates a new serial device
func NewSerialDevice(numLed int) LedDevice {
	s := new(serialDevice)
	s.config = &serial.Config{
		Name:        "/dev/ttyUSB0",
		Baud:        1152000,
		ReadTimeout: 1000 * time.Millisecond,
		Size:        8,
	}
	s.numLed = numLed
	return s
}

func (s *serialDevice) Open() error {
	var err error
	s.stream, err = serial.OpenPort(s.config)
	if err != nil {
		return err
	}
	return err
}

func (s *serialDevice) init() {
	log.Println("Arduino version:")
	command := "V\n"
	s.Write([]byte(command))
	command = fmt.Sprintf("I%04x\n", s.numLed)
	s.Write([]byte(command))
}

func (s *serialDevice) read(wg *sync.WaitGroup) {
	lastLine := ""
	defer wg.Done()
	defer log.Println(lastLine)

	buf := make([]byte, 1024)
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
			s := lastLine + string(buf[:n])
			lines := strings.Split(s, "\n")
			if s[len(s)-1] != '\n' {
				numLines := len(lines) - 1
				lastLine = lines[numLines]
				lines = lines[:numLines]
			} else {
				lastLine = ""
			}
			for _, line := range lines {
				if len(line) > 0 {
					log.Println(line)
				}
			}
		}
	}

}

func (s *serialDevice) Close() error {
	return nil
}

func (s *serialDevice) Write(data []byte) (int, error) {
	// log.Print("write ", string(data))
	n, err := s.stream.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	return n, err
}

func (s *serialDevice) SetInput(input <-chan []byte) {
	s.input = input
}

func (s *serialDevice) setPixel(pixel int, buffer []byte) {
	bufIndex := pixel * NumBytePerColor
	command := fmt.Sprintf("P%04x%02x%02x%02x\n", pixel, buffer[bufIndex+0], buffer[bufIndex+1], buffer[bufIndex+2])
	s.Write([]byte(command))
}

func (s *serialDevice) shade(pixel int, buffer []byte) {
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
	wg.Add(1)
	go s.read(wg)

	const latchDelay = 11 * time.Millisecond
	lastFrameTime := time.Now().Add(-latchDelay)

	// initialize bitbanger with number of leds
	time.Sleep(latchDelay)
	s.init()

	for frame := range s.input {
		now := time.Now()
		sleepDuration := latchDelay - (now.Sub(lastFrameTime))
		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}
		// log.Print("frame ", frame)
		for pixel := 0; pixel < s.numLed; pixel++ {
			s.setPixel(pixel, frame)
			time.Sleep(100 * time.Microsecond)
		}
		s.latchFrame()
		//s.shade(s.numLed, frame[0:3])
	}
}

func (s *serialDevice) GetType() Type {
	return Serial
}
