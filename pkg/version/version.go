package version

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
)

// PrintProgramInfo starts a go routine to print program details in a regular manner into the log
func PrintProgramInfo(compileDate, versionTag string, logger *zap.Logger) {
	const intervalSeconds = 30
	program := strings.Split(os.Args[0], "/")
	programName := program[len(program)-1]
	progLogger := logger.With(
		zap.String("program", programName),
		zap.String("compileDate", compileDate),
		zap.String("versionTag", versionTag),
	)
	go func() {
		for {
			progLogger.Info("program version info")
			time.Sleep(intervalSeconds * time.Second)
		}
	}()
}
