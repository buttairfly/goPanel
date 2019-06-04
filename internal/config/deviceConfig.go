package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/tarm/serial"
)

// DeviceConfig is the config of the type of device
type DeviceConfig struct {
	Type         Type          `json:"type"`
	SerialConfig *SerialConfig `json:"serialConfig,omitempty"`
}

func NewDeviceConfigFromPath(path string) (*DeviceConfig, error) {
	dc := new(DeviceConfig)
	err := dc.FromFile(path)
	if err != nil {
		return nil, err
	}
	return dc, nil
}

// FromFile reads the config from a file at path
func (dc *DeviceConfig) FromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can not read Config file %v. error: %v", path, err)
	}
	defer f.Close()
	return dc.FromReader(f)
}

// FromReader decodes the config from io.Reader
func (dc *DeviceConfig) FromReader(r io.Reader) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(&*dc)
	if err != nil {
		return fmt.Errorf("can not decode json. error: %v", err)
	}
	return nil
}

// WriteToFile writes the config to a file at path
func (dc *DeviceConfig) WriteToFile(path string) error {
	jsonConfig, err := json.MarshalIndent(dc, "", "\t")
	if err != nil {
		return err
	}
	jsonConfig = append(jsonConfig, byte('\n'))
	return ioutil.WriteFile(path, jsonConfig, 0622)
}

const (
	// Print debug print device
	Print = Type("print")
	// Serial high level serial tty device
	Serial = Type("serial")
)

// Type is a LedDevice type
type Type string

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

type StreamConfig struct {
	Name        string          `json:"name"`
	Baud        int             `json:"baud"`
	Size        byte            `json:"size"`
	ReadTimeout time.Duration   `json:"readTimeout,omitempty"`
	Parity      serial.Parity   `json:"parity,omitempty"`
	StopBits    serial.StopBits `json:"stopBits,omitempty"`
}
type aliasStreamConfig struct {
	Name        string          `json:"name"`
	Baud        int             `json:"baud"`
	Size        byte            `json:"size"`
	ReadTimeout string          `json:"readTimeout,omitempty"`
	Parity      serial.Parity   `json:"parity,omitempty"`
	StopBits    serial.StopBits `json:"stopBits,omitempty"`
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
