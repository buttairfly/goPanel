package arduinocom

import (
	"time"
)

// SerialConfig is the serial config
type SerialConfig struct {
	StreamConfig       *StreamConfig       `json:"streamConfig" yaml:"streamConfig"`
	Verbose            bool                `json:"verbose" yaml:"verbose"`
	VerboseArduino     bool                `json:"verboseArduino" yaml:"verboseArduino"`
	ReadBufferSize     int                 `json:"readBufferSize" yaml:"readBufferSize"`
	RawFramePartNumLed int                 `json:"rawFramePartNumLed" yaml:"rawFramePartNumLed"`
	ParitySeed         byte                `json:"paritySeed" yaml:"paritySeed"`
	InitSleepTime      time.Duration       `json:"initSleepTime,omitempty" yaml:"initSleepTime,omitempty"`
	CmdSleepTime       time.Duration       `json:"cmdSleepTime,omitempty" yaml:"cmdSleepTime,omitempty"`
	ArduinoErrorConfig *ArduinoErrorConfig `json:"arduinoErrorConfig,omitempty" yaml:"arduinoErrorConfig,omitempty"`
}
