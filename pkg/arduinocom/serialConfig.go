package arduinocom

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v2"
)

// SerialConfig is the serial config
type SerialConfig struct {
	StreamConfig       *StreamConfig       `yaml:"streamConfig"`
	ArduinoErrorConfig *ArduinoErrorConfig `yaml:"arduinoErrorConfig,omitempty"`
	Verbose            bool                `yaml:"verbose"`
	ReadBufferSize     int                 `yaml:"readBufferSize"`
	InitSleepTime      time.Duration       `yaml:"initSleepTime,omitempty"`
	LatchSleepTime     time.Duration       `yaml:"latchSleepTime,omitempty"`
	CommandSleepTime   time.Duration       `yaml:"commandSleepTime,omitempty"`
}
type aliasSerialConfig struct {
	StreamConfig       *StreamConfig       `yaml:"streamConfig"`
	ArduinoErrorConfig *ArduinoErrorConfig `yaml:"arduinoErrorConfig,omitempty"`
	Verbose            bool                `yaml:"verbose"`
	ReadBufferSize     int                 `yaml:"readBufferSize"`
	InitSleepTime      string              `yaml:"initSleepTime,omitempty"`
	LatchSleepTime     string              `yaml:"latchSleepTime,omitempty"`
	CommandSleepTime   string              `yaml:"commandSleepTime,omitempty"`
}

// UnmarshalYAML unmarshals YAMLDuration
func (sc *SerialConfig) UnmarshalYAML(b []byte) error {
	var tmp aliasSerialConfig
	err := yaml.Unmarshal(b, &tmp)
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

// MarshalYAML marshals SerialConfig to yaml or error
func (sc *SerialConfig) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(&aliasSerialConfig{
		StreamConfig:       sc.StreamConfig,
		ArduinoErrorConfig: sc.ArduinoErrorConfig,
		Verbose:            sc.Verbose,
		ReadBufferSize:     sc.ReadBufferSize,
		InitSleepTime:      fmt.Sprintf("%s", time.Duration(sc.InitSleepTime)),
		LatchSleepTime:     fmt.Sprintf("%s", time.Duration(sc.LatchSleepTime)),
		CommandSleepTime:   fmt.Sprintf("%s", time.Duration(sc.CommandSleepTime)),
	})
}
