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
					s.checkInitDone(line)
				}
			}
		}
	}
}

func (s *serialDevice) printStats(wg *sync.WaitGroup) {
	defer wg.Done()

	for stat := range s.stats {
		log.Printf("%v: %s", stat.timeStamp, stat.message)
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
