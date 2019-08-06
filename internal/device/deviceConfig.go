package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	arduinocomconfig "github.com/buttairfly/goPanel/pkg/arduinocom/config"
)

// DeviceConfig is the config of the type of device
type DeviceConfig struct {
	Type         Type                           `json:"type"`
	SerialConfig *arduinocomconfig.SerialConfig `json:"serialConfig,omitempty"`
}

// NewDeviceConfigFromPath returns a new DeviceConfig or error
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

// Type is a LedDevice type
type Type string

const (
	// Print debug print device
	Print = Type("print")
	// Serial high level serial tty device
	Serial = Type("serial")
)
