package version

import (
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type version struct {
	tag         string
	compileDate string
	programName string
	logger      *zap.Logger
}

var programVersion version

// Init initializes the program version struct
func Init(compileDate, tag string, logger *zap.Logger) {

	programParts := strings.Split(os.Args[0], "/")
	programName := programParts[len(programParts)-1]

	programVersion.tag = tag
	programVersion.compileDate = compileDate
	programVersion.programName = programName
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
