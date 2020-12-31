package version

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Version holds all compile data
type Version struct {
	Name        string `json:"name"        yaml:"name"`
	ProgramName string `json:"programName" yaml:"programName"`
	Tag         string `json:"tag"         yaml:"tag"`
	CompileDate string `json:"compileDate" yaml:"compileDate"`
	interval    time.Duration
	logger      *zap.Logger
}

var versions []Version

// New returns a new Version struct and adds it to the Versions slice
func New(name, programName, compileDate, tag string, interval time.Duration, logger *zap.Logger) *Version {
	newVersion := Version{
		Name:        name,
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
	versions = append(versions, newVersion)
	return &newVersion
}

// GetVersions returns the applied version array
func GetVersions() []Version {
	return versions
}

// GetVersionByName returns the Version with name or error
func GetVersionByName(name string) (*Version, error) {
	for _, current := range versions {
		if current.Name == name {
			return &current, nil
		}
	}
	return nil, fmt.Errorf("unknown version %s", name)
}

// Log loggs the logger once
func (v *Version) Log() {
	v.logger.Info("version")
}

// Run starts a go routine to print program details in a regular manner into the log
func (v *Version) Run(cancelCtx context.Context) {
	ticker := time.NewTicker(v.interval)
	defer ticker.Stop()

	if v.interval > 0 {
		for {
			select {
			case <-cancelCtx.Done():
				{
					v.logger.Info("version exit")
					return
				}
			case <-ticker.C:
				{
					v.Log()
				}
			}
		}
	}
}
