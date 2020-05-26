package version

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// Version holds all compile data
type Version struct {
	ProgramName string `json:"programName" yaml:"programName"`
	Tag         string `json:"tag"         yaml:"tag"`
	CompileDate string `json:"compileDate" yaml:"compileDate"`
	interval    time.Duration
	logger      *zap.Logger
}

// Versions are the
var Versions []Version

// New returns a new Version struct and adds it to the Versions slice
func New(programName, compileDate, tag string, interval time.Duration, logger *zap.Logger) *Version {
	newVersion := Version{
		Tag:         tag,
		CompileDate: compileDate,
		ProgramName: programName,
		interval:    interval,
		logger: logger.With(
			zap.String("programName", programName),
			zap.String("compileDate", compileDate),
			zap.String("tag", tag),
		),
	}
	Versions = append(Versions, newVersion)
	return &newVersion
}

// Run starts a go routine to print program details in a regular manner into the log
func (v *Version) Run(cancelCtx context.Context) {
	for {
		select {
		case <-cancelCtx.Done():
			{
				v.logger.Info("exit")
				return
			}
		case <-time.After(v.interval):
			{
				v.logger.Info("version")
			}
		}
	}
}
