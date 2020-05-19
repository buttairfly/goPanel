package arduinocom

import (
	"time"
)

// SerialConfig is the serial config
type SerialConfig struct {
	StreamConfig       *StreamConfig       `yaml:"streamConfig"`
	Verbose            bool                `yaml:"verbose"`
	VerboseArduino     bool                `yaml:"verboseArduino"`
	ReadBufferSize     int                 `yaml:"readBufferSize"`
	ParitySeed         byte                `yaml:"paritySeed"`
	InitSleepTime      time.Duration       `yaml:"initSleepTime,omitempty"`
	LatchSleepTime     time.Duration       `yaml:"latchSleepTime,omitempty"`
	CommandSleepTime   time.Duration       `yaml:"commandSleepTime,omitempty"`
	ArduinoErrorConfig *ArduinoErrorConfig `yaml:"arduinoErrorConfig,omitempty"`
}
