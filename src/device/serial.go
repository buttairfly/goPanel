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
		ReadTimeout: 300 * time.Millisecond,
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
	s.init()
	return err
}

func (s *serialDevice) init() {
	log.Println("INITIALIZE", s.numLed)
	command := fmt.Sprintf("I%04x\n", s.numLed)
	s.Write([]byte(command))
}

func (s *serialDevice) read(wg *sync.WaitGroup) {
	defer wg.Done()
	defer s.Close()

	buf := make([]byte, 1024)
	for {
		_, ok := <-s.readActive
		if !ok {
			return
		}

		n, err := s.stream.Read(buf)
		if err != nil {
			if err == io.EOF {
				continue
			}
			log.Fatal(err)
		}
		lines := strings.Split(string(buf[:n]), "\n")
		for _, line := range lines {
			log.Println(line)
		}
	}
}

func (s *serialDevice) Close() error {
	return nil
}

func (s *serialDevice) Write(data []byte) (int, error) {
	log.Print("write ", string(data))
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

	for frame := range s.input {
		/*
			for pixel := 0; pixel < s.numLed; pixel++ {
				s.setPixel(pixel, frame)
			}
			s.latchFrame()
		*/
		s.shade(s.numLed, frame[0:3])
	}
}

func (s *serialDevice) GetType() Type {
	return Serial
}
