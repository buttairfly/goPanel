package version


import (
	"log"
	"time"
	"os"
	"strings"
)

// PrintProgramInfo starts a go routine to print program details in a regular manner into the log
func PrintProgramInfo(compileDate, versionTag string) {
	const intervalSeconds = 30
	program := strings.Split(os.Args[0], "/")
	programName := program[len(program)-1]
	go func() {
		for {
			log.Printf("%s: compiled at %s with version %s", programName, compileDate, versionTag)
			time.Sleep(intervalSeconds * time.Second)
		}
	}()
}
