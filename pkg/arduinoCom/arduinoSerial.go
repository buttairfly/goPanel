package com

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/buttairfly/goPanel/internal/config"
	"github.com/tarm/serial"
)

// ArduinoCom is a serial communication with arduino friendly protocol
type ArduinoCom struct {
	config     *config.SerialConfig
	stream     *serial.Port
	readActive chan bool
	initDone   chan bool
	stats      chan *stats
	latchEnd   chan bool
	latched    int64
}

// NewArduinoCom creates a new serial arduino communication
func NewArduinoCom(numLed int, sc *config.SerialConfig) *ArduinoCom {
	a := new(ArduinoCom)
	a.config = sc
	return a
}

// Init initializes a arduino serial
func (a *ArduinoCom) Init(numLed int) {
	defer func() {
		if !a.config.Verbose {
			a.sendInitComand("Q0fff\n")
		}
	}()

	for {
		select {
		case <-a.initDone:
			return
		default:
			a.stream.Flush()
			a.sendInitComand("Q0000\n")
			a.sendInitComand("V\n")
			if numLed > 0 {
				a.sendInitComand(fmt.Sprintf("I%04x\n", numLed))
			}
		}
	}
}

func (a *ArduinoCom) sendInitComand(command string) {
	a.Write([]byte(command))
	time.Sleep(a.config.InitSleepTime)
}

func (a *ArduinoCom) read(wg *sync.WaitGroup) {
	lastLine := ""
	defer wg.Done()
	defer log.Println(lastLine)

	buf := make([]byte, a.config.ReadBufferSize)
	for {
		select {
		case <-a.readActive:
			return
		default:
			n, err := a.stream.Read(buf)
			timeStamp := time.Now()
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
					a.checkInitDone(line)
					stat := &stats{
						event:     printType,
						timeStamp: timeStamp,
						message:   line,
					}
					if IsArduinoError(line) {
						arduinoError, err := NewArduinoError(a.config.ArduinoErrorConfig, line)
						if err != nil {
							log.Print(err)
							continue
						}
						stat.event = ardoinoErrorType
						stat.message = arduinoError.Error()
					}
					a.stats <- stat
				}
			}
		}
	}
}

func (a *ArduinoCom) printStats(wg *sync.WaitGroup) {
	defer wg.Done()

	for stat := range a.stats {
		if stat.event != latchType {
			timeStamp := fmt.Sprintf("%02d.%06d", stat.timeStamp.Second(), stat.timeStamp.Nanosecond()/int(time.Microsecond))
			log.Printf("%s %s: %s", stat.event, timeStamp, stat.message)
		} else {
			a.latched++
		}
	}
}

func (a *ArduinoCom) checkInitDone(line string) {
	if a.needsInit() {
		parts := strings.Split(line, " ")
		if len(parts) == 2 && parts[0] == "Init" {
			initLed, err := strconv.ParseInt(parts[1], 16, 16)
			if err != nil {
				log.Print(err)
			} else {
				if int(initLed) == s.numLed {
					close(s.initDone)
				}
			}
		}
	}
}

func (s *SerialCom) needsInit() bool {
	select {
	case <-s.initDone:
		return false
	default:
		return true
	}
}
