package version

import (
	"context"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Version holds all compile data
type Version struct {
	Tag         string        `json:"tag"         yaml:"tag"`
	CompileDate string        `json:"compileDate" yaml:"compileDate"`
	ProgramName string        `json:"programName" yaml:"programName"`
	Interval    time.Duration `json:"interval" yaml:"interval"`
	logger      *zap.Logger
}

var programVersion Version

// New initializes returns a new Version struct
func New(compileDate, tag string, interval time.Duration, logger *zap.Logger) *Version {

	programParts := strings.Split(os.Args[0], "/")
	programName := programParts[len(programParts)-1]

	return &Version{
		Tag:         tag,
		CompileDate: compileDate,
		ProgramName: programName,
		Interval:    interval,
		logger: logger.With(
			zap.String("programName", programName),
			zap.String("compileDate", compileDate),
			zap.String("tag", tag),
		),
	}
}

// Run starts a go routine to print program details in a regular manner into the log
func (v *Version) Run(cancelCtx context.Context) {
	print("version")
	for {
		select {
		case <-cancelCtx.Done():
			{
				print("exit")
				return
			}
		case <-time.After(v.Interval):
			{
				print("version")
			}
		}
	}
}

func (v *Version) print(text string) {
	v.logger.Info(text)
}
