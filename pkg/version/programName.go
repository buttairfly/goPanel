package version

import (
	"os"
	"strings"
)

// GetProgramName return the name of the current program
func GetProgramName() string {
	programParts := strings.Split(os.Args[0], "/")
	programName := programParts[len(programParts)-1]
	return programName
}
