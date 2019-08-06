package arduinocom

import (
	"encoding/json"
	"fmt"
	"time"
)

// SerialConfig is the serial config
type SerialConfig struct {
	StreamConfig       *StreamConfig       `json:"streamConfig"`
	ArduinoErrorConfig *ArduinoErrorConfig `json:"arduinoErrorConfig,omitempty"`
	Verbose            bool                `json:"verbose"`
	ReadBufferSize     int                 `json:"readBufferSize"`
	InitSleepTime      time.Duration       `json:"initSleepTime,omitempty"`
	LatchSleepTime     time.Duration       `json:"latchSleepTime,omitempty"`
	CommandSleepTime   time.Duration       `json:"commandSleepTime,omitempty"`
}
type aliasSerialConfig struct {
	StreamConfig       *StreamConfig       `json:"streamConfig"`
	ArduinoErrorConfig *ArduinoErrorConfig `json:"arduinoErrorConfig,omitempty"`
	Verbose            bool                `json:"verbose"`
	ReadBufferSize     int                 `json:"readBufferSize"`
	InitSleepTime      string              `json:"initSleepTime,omitempty"`
	LatchSleepTime     string              `json:"latchSleepTime,omitempty"`
	CommandSleepTime   string              `json:"commandSleepTime,omitempty"`
}

// UnmarshalJSON unmarshals JSONDuration
func (sc *SerialConfig) UnmarshalJSON(b []byte) error {
	var tmp aliasSerialConfig
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	sc.InitSleepTime, err = time.ParseDuration(tmp.InitSleepTime)
	if err != nil {
		return fmt.Errorf("failed to parse initSleepTime '%s' to time.Duration: %v", tmp.InitSleepTime, err)
	}

	sc.LatchSleepTime, err = time.ParseDuration(tmp.LatchSleepTime)
	if err != nil {
		return fmt.Errorf("failed to parse latchSleepTime '%s' to time.Duration: %v", tmp.LatchSleepTime, err)
	}

	sc.CommandSleepTime, err = time.ParseDuration(tmp.CommandSleepTime)
	if err != nil {
		return fmt.Errorf("failed to parse commandSleepTime '%s' to time.Duration: %v", tmp.CommandSleepTime, err)
	}

	sc.StreamConfig = tmp.StreamConfig
	sc.ArduinoErrorConfig = tmp.ArduinoErrorConfig
	sc.Verbose = tmp.Verbose
	sc.ReadBufferSize = tmp.ReadBufferSize

	return nil
}

// MarshalJSON marshals SerialConfig to json or error
func (sc *SerialConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(&aliasSerialConfig{
		StreamConfig:       sc.StreamConfig,
		ArduinoErrorConfig: sc.ArduinoErrorConfig,
		Verbose:            sc.Verbose,
		ReadBufferSize:     sc.ReadBufferSize,
		InitSleepTime:      fmt.Sprintf("%s", time.Duration(sc.InitSleepTime)),
		LatchSleepTime:     fmt.Sprintf("%s", time.Duration(sc.LatchSleepTime)),
		CommandSleepTime:   fmt.Sprintf("%s", time.Duration(sc.CommandSleepTime)),
	})
}
