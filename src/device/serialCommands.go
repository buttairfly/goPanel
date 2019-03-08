package device

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

func (s *serialDevice) init() {
	defer func() {
		if !s.config.Verbose {
			s.sendInitComand("Q0fff\n")
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
					s.checkInitDone(line)
					stat := &stats{
						event:     printType,
						timeStamp: timeStamp,
						message:   line,
					}
					if IsArduinoError(line) {
						arduinoError, err := NewArduinoError(s.config.ArduinoErrorConfig, line)
						if err != nil {
							log.Print(err)
							continue
						}
						stat.event = ardoinoErrorType
						stat.message = arduinoError.Error()
					}
					s.stats <- stat
				}
			}
		}
	}
}

func (s *serialDevice) printLatches(wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	for {
		select {
		case <-s.latchEnd:
			return
		default:
			timeDiff := time.Now().Sub(start) / time.Second
			log.Printf("Latched frames: %f/s last diff: %v", float64(s.latched)/float64(timeDiff), timeDiff)
			time.Sleep(30 * time.Second)
		}
	}
}

func (s *serialDevice) printStats(wg *sync.WaitGroup) {
	defer wg.Done()

	for stat := range s.stats {
		if stat.event != latchType {
			timeStamp := fmt.Sprintf("%02d.%06d", stat.timeStamp.Second(), stat.timeStamp.Nanosecond()/1000)
			log.Printf("%s %s: %s", stat.event, timeStamp, stat.message)
		} else {
			s.latched++
		}
	}
}

func (s *serialDevice) checkInitDone(line string) {
	if s.needsInit() {
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

func (s *serialDevice) needsInit() bool {
	select {
	case <-s.initDone:
		return false
	default:
		return true
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
