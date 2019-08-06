package arduinocom

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tarm/serial"
)

// ArduinoCom is a serial communication with arduino friendly protocol
//
// Set numLed to 0, when arduinoCom does not use led length
type ArduinoCom struct {
	config     *SerialConfig
	stream     *serial.Port
	readActive chan bool
	initDone   chan bool
	stats      chan *stats
	latched    int64
	paritySeed byte
	numLed     int
}

// NewArduinoCom creates a new serial arduino communication
//
// Set numLed to 0 when not needed as configureable one time parameter
func NewArduinoCom(numLed int, sc *SerialConfig) *ArduinoCom {
	a := new(ArduinoCom)
	a.config = sc
	a.readActive = make(chan bool)
	a.initDone = make(chan bool)
	a.stats = make(chan *stats, 10)
	return a
}

// Open opens port for aruduino serial connection
func (a *ArduinoCom) Open() error {
	var err error
	a.stream, err = serial.OpenPort(a.config.StreamConfig.ToStreamSerialConfig())
	return err
}

// Close removes serial arduino stream and all helper channels
func (a *ArduinoCom) Close() error {
	close(a.readActive)
	close(a.stats)
	return a.stream.Close()
}

// Init initializes a arduino serial
func (a *ArduinoCom) Init() {
	defer func() {
		if !a.config.Verbose {
			a.sendInitComand("Q0fff")
		}
	}()

	for {
		select {
		case <-a.initDone:
			return
		default:
			a.stream.Flush()
			a.sendInitComand("Q0000")
			a.sendInitComand("V")
			if a.numLed > 0 {
				a.sendInitComand(fmt.Sprintf("I%04x", a.numLed))
			} else {
				close(a.initDone)
			}
		}
	}
}

// Read is the function to handle arduinoCom reads
//
// Read prints ArduinoErrors, when neccesary and shows debug information when configured
func (a *ArduinoCom) Read(wg *sync.WaitGroup) {
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

func (a *ArduinoCom) Write(data string) (int, error) {
	if a.config.Verbose {
		log.Printf("Command %s", string(data))
	}
	n, err := a.stream.Write([]byte(data))
	if err != nil {
		log.Fatal(err)
	}
	return n, err
}

// CalcParityAndWrite calculates the parity and writes it to serial
func (a *ArduinoCom) CalcParityAndWrite(command string) (int, error) {
	return a.Write(CalcHexParity(command, a.paritySeed))
}

func (a *ArduinoCom) sendInitComand(command string) {
	time.Sleep(a.config.InitSleepTime)
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
				if int(initLed) == a.numLed {
					close(a.initDone)
				}
			}
		}
	}
}

func (a *ArduinoCom) needsInit() bool {
	select {
	case <-a.initDone:
		return false
	default:
		return true
	}
}
