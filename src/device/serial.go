package device

import (
	"fmt"
	"log"
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
func NewSerialDevice(numLed int) *serialDevice {
	s := new(serialDevice)
	s.config = &serial.Config{
		Name:        "/dev/ttyUSB0",
		Baud:        1152000,
		ReadTimeout: 500 * time.Millisecond,
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
	command := "I00C8\n"
	_, err := s.stream.Write([]byte(command))
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)
	n, errRead := s.stream.Read(buf)
	if errRead != nil {
		log.Fatal(errRead)
	}
	log.Printf("%q", buf[:n])
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
	for buffer := range s.input {
		_, err := s.Write(buffer)
		if err != nil {
			log.Panic(err)
		}
	}
}

func (s *serialDevice) GetType() Type {
	return Serial
}
