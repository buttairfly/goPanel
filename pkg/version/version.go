package version

import (
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Version holds all compile data
type Version struct {
	Tag         string `json:"tag"         yaml:"tag"`
	CompileDate string `json:"compileDate" yaml:"compileDate"`
	ProgramName string `json:"programName" yaml:"programName"`
	logger      *zap.Logger
}

var programVersion Version

// Init initializes the program version struct
func Init(compileDate, tag string, logger *zap.Logger) {

	programParts := strings.Split(os.Args[0], "/")
	programName := programParts[len(programParts)-1]

	programVersion.Tag = tag
	programVersion.CompileDate = compileDate
	programVersion.ProgramName = programName
	programVersion.logger = logger.With(
		zap.String("programName", programName),
		zap.String("compileDate", compileDate),
		zap.String("tag", tag),
	)
}

// Run starts a go routine to print program details in a regular manner into the log
func Run(wg *sync.WaitGroup, exitChan <-chan bool) {
	defer wg.Done()
	const channelTestInterval = 500 * time.Millisecond
	const printInterval = 10 * time.Second
	lastPrint := time.Now()
	print("version")
	for {
		select {
		case _, ok := <-exitChan:
			if !ok {
				print("exit program")
				return
			}
		default:
			now := time.Now()
			if now.Sub(lastPrint) > printInterval {
				print("version")
				lastPrint = now
			}
			time.Sleep(channelTestInterval)
		}
	}
}

func print(text string) {
	if programVersion.logger != nil {
		programVersion.logger.Info(text)
	}
}
