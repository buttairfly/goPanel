package arduinocom

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tarm/serial"
)

// StreamConfig holds all config values needed for the serial stream tarm package
type StreamConfig struct {
	Name        string          `json:"name" yaml:"name"`
	Baud        int             `json:"baud" yaml:"baud"`
	Size        byte            `json:"size" yaml:"size"`
	ReadTimeout time.Duration   `json:"readTimeout,omitempty" yaml:"readTimeout,omitempty"`
	Parity      serial.Parity   `json:"parity,omitempty" yaml:"parity,omitempty"`
	StopBits    serial.StopBits `json:"stopBits,omitempty" yaml:"stopBits,omitempty"`
}
type aliasStreamConfig struct {
	Name        string          `json:"name" yaml:"name"`
	Baud        int             `json:"baud" yaml:"baud"`
	Size        byte            `json:"size" yaml:"size"`
	ReadTimeout string          `json:"readTimeout,omitempty" yaml:"readTimeout,omitempty"`
	Parity      serial.Parity   `json:"parity,omitempty" yaml:"parity,omitempty"`
	StopBits    serial.StopBits `json:"stopBits,omitempty" yaml:"stopBits,omitempty"`
}

// UnmarshalJSON unmarshals JSONDuration
func (sc *StreamConfig) UnmarshalJSON(b []byte) error {
	var tmp aliasStreamConfig
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	sc.ReadTimeout, err = time.ParseDuration(tmp.ReadTimeout)
	if err != nil {
		return fmt.Errorf("failed to parse readTimeout '%s' to time.Duration: %v", tmp.ReadTimeout, err)
	}

	sc.Name = tmp.Name
	sc.Baud = tmp.Baud
	sc.Size = tmp.Size
	sc.Parity = tmp.Parity
	sc.StopBits = tmp.StopBits

	return nil
}

// MarshalJSON marshals a StreamConfig to json or error
func (sc *StreamConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(&aliasStreamConfig{
		Name:        sc.Name,
		Baud:        sc.Baud,
		Size:        sc.Size,
		ReadTimeout: fmt.Sprintf("%s", time.Duration(sc.ReadTimeout)),
		Parity:      sc.Parity,
		StopBits:    sc.StopBits,
	})
}

// ToStreamSerialConfig serializes a StreamConfig to a tarm serial Config
func (sc *StreamConfig) ToStreamSerialConfig() *serial.Config {
	return &serial.Config{
		Name:        sc.Name,
		Baud:        sc.Baud,
		ReadTimeout: sc.ReadTimeout,
		Size:        sc.Size,
		Parity:      sc.Parity,
		StopBits:    sc.StopBits,
	}
}
