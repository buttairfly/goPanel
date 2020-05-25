package arduinocom

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tarm/serial"
	"go.uber.org/zap"
)

// ArduinoCom is a serial communication with arduino friendly protocol
//
// Set numLed to 0, when arduinoCom does not use led length
type ArduinoCom struct {
	config     *SerialConfig
	stream     *serial.Port
	initDone   chan bool
	stats      chan *Stat
	latched    int64
	paritySeed byte
	numLed     int
	logger     *zap.Logger
}

// NewArduinoCom creates a new serial arduino communication
//
// Set numLed to 0 when not needed as configureable one time parameter
func NewArduinoCom(numLed int, sc *SerialConfig, logger *zap.Logger) *ArduinoCom {
	a := new(ArduinoCom)
	a.config = sc
	a.initDone = make(chan bool)
	a.stats = make(chan *Stat, 10)
	a.numLed = numLed
	a.paritySeed = sc.ParitySeed
	a.logger = logger
	return a
}

// AddStat adds a stat on the channel
func (a *ArduinoCom) AddStat(stat *Stat) {
	a.stats <- stat
}

// Config reuturns the serial config
func (a *ArduinoCom) Config() *SerialConfig {
	return a.config
}

// Open opens port for aruduino serial connection
func (a *ArduinoCom) Open() error {
	var err error
	a.stream, err = serial.OpenPort(a.config.StreamConfig.ToStreamSerialConfig())
	return err
}

// Close removes serial arduino stream and all helper channels
func (a *ArduinoCom) Close() error {
	err := a.stream.Close()
	return err
}

// Init initializes a arduino serial
func (a *ArduinoCom) Init() {
	defer func() {
		if !a.config.VerboseArduino {
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
func (a *ArduinoCom) Read(cancelCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	lastLine := ""
	buf := make([]byte, a.config.ReadBufferSize)

	for {
		select {
		case <-cancelCtx.Done():
			return
		default:
			n, err := a.stream.Read(buf)
			timeStamp := time.Now()
			if err != nil {
				if err == io.EOF {
					continue
				}
				a.logger.Fatal("stream read", zap.Error(err))
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
					a.logger.Info("arduinoRead", zap.String("line", line))
					a.checkInitDone(line)
					stat := &Stat{
						Event:     PrintStatType,
						TimeStamp: timeStamp,
						Message:   line,
					}
					if IsArduinoError(line) {
						arduinoError, err := NewArduinoError(a.config.ArduinoErrorConfig, line)
						if err != nil {
							a.logger.Warn("arduino", zap.Error(err))
							continue
						}
						stat.Event = ArdoinoErrorStatType
						stat.Message = arduinoError.Error()
					}
					a.AddStat(stat)
				}
			}
		}
	}
}

func (a *ArduinoCom) Write(data string) (int, error) {
	if a.config.Verbose {
		a.logger.Info("verbose", zap.String("command", data))
	}
	n, err := a.stream.Write([]byte(data))
	if err != nil {
		a.logger.Fatal("stream write", zap.Error(err))
	}
	return n, err
}

// CalcParityAndWrite calculates the parity and writes it to serial
func (a *ArduinoCom) CalcParityAndWrite(command string) (int, error) {
	return a.Write(CalcHexParity(command, a.paritySeed))
}

func (a *ArduinoCom) sendInitComand(command string) {
	a.CalcParityAndWrite(command)
	time.Sleep(a.config.InitSleepTime)
}

// PrintStats prints all stats for running serial connection
func (a *ArduinoCom) PrintStats(cancelCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-cancelCtx.Done():
			return
		case stat := <-a.stats:
			{
				if stat.Event != LatchStatType {
					a.logger.Info(
						"stats logger",
						zap.String("arduino", string(stat.Event)),
						zap.String("eventTine", fmt.Sprintf("%02d.%06d", stat.TimeStamp.Second(), stat.TimeStamp.Nanosecond()/int(time.Microsecond))),
						zap.String("message", stat.Message),
					)
				} else {
					a.latched++
				}
			}
		}
	}
}

func (a *ArduinoCom) checkInitDone(line string) {
	if a.needsInit() {
		parts := strings.Split(line, " ")
		if len(parts) == 2 && parts[0] == "Init" {
			initLed, err := strconv.ParseInt(parts[1], 16, 16)
			if err != nil {
				a.logger.Warn("not initialized", zap.Error(err))
			} else {
				if int(initLed) == a.numLed {
					a.logger.Info("initialized", zap.Int("numLed", a.numLed))
					close(a.initDone)
				} else {
					a.logger.Warn("not initizalized: ", zap.Int("numLed", a.numLed), zap.Int("initLed", int(initLed)))
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
