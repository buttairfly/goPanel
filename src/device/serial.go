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
	config *serial.Config
	stream *serial.Port
	numLed int
	input  <-chan []byte
}

// NewSerialDevice creates a new serial device
func NewSerialDevice(numLed int) LedDevice {
	s := new(serialDevice)
	s.config = &serial.Config{
		Name:        "/dev/ttyUSB0",
		Baud:        1152000,
		ReadTimeout: 2000 * time.Millisecond,
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
	command := fmt.Sprintf("I%4x\n", s.numLed)
	_, err := s.stream.Write([]byte(command))
	if err != nil {
		log.Fatal(err)
	}
}

func (s *serialDevice) read() {
	buf := make([]byte, 1024)
	n, err := s.stream.Read(buf)
	if err != nil {
		if err == io.EOF {
			return
		}
		log.Fatal(err)
	}
	lines := strings.Split(string(buf[:n]), "\n")
	for _, line := range lines {
		log.Println(line)
	}
}

func (s *serialDevice) Close() error {
	return nil
}

func (s *serialDevice) Write(data []byte) (int, error) {
	if len(data) != s.numLed*NumBytePerColor {
		return 0, fmt.Errorf(
			"could not write %v bytes of data, %v is needed",
			len(data), s.numLed*NumBytePerColor)
	}
	n, err := s.stream.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	return n, err
}

func (s *serialDevice) SetInput(input <-chan []byte) {
	s.input = input
}

func (s *serialDevice) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer s.Close()
	bufferSize := s.numLed * NumBytePerColor * NumByteToRepresentHex
	buffer := make([]byte, bufferSize, bufferSize)
	for frame := range s.input {
		for i, b := range frame {
			s := fmt.Sprintf("%2x", b)
			bufIndex := i * NumByteToRepresentHex
			buffer[bufIndex+0] = s[0]
			buffer[bufIndex+1] = s[1]
		}
		for pixel := 0; pixel < s.numLed; pixel++ {
			bufIndex := pixel * NumBytePerColor
			command := fmt.Sprintf("P%4x%2x%2x%2x\n", pixel, buffer[bufIndex+0], buffer[bufIndex+1], buffer[bufIndex+2])
			_, err := s.stream.Write([]byte(command))
			if err != nil {
				log.Fatal(err)
			}
			command = "S\n"
			_, err = s.stream.Write([]byte(command))
			if err != nil {
				log.Fatal(err)
			}
		}

		s.read()
	}
}

func (s *serialDevice) GetType() Type {
	return Serial
}
