package device

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/buttairfly/goPanel/pkg/arduinocom"
)

// DeviceConfig is the config of the type of device
type DeviceConfig struct {
	Type         Type                     `yaml:"type"`
	SerialConfig *arduinocom.SerialConfig `yaml:"serialConfig,omitempty"`
}

// NewDeviceConfigFromPath returns a new DeviceConfig or error
func NewDeviceConfigFromPath(path string) (*DeviceConfig, error) {
	dc := new(DeviceConfig)
	err := dc.FromYamlFile(path)
	if err != nil {
		return nil, err
	}
	return dc, nil
}

// FromYamlFile reads the config from a file at path
func (dc *DeviceConfig) FromYamlFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can not read Config file %v. error: %v", path, err)
	}
	defer f.Close()
	return dc.FromYamlReader(f)
}

// FromYamlReader decodes the config from io.Reader
func (dc *DeviceConfig) FromYamlReader(r io.Reader) error {
	dec := yaml.NewDecoder(r)
	err := dec.Decode(&*dc)
	if err != nil {
		return fmt.Errorf("can not decode json. error: %v", err)
	}
	return nil
}

// WriteToYamlFile writes the config to a file at path
func (dc *DeviceConfig) WriteToYamlFile(path string) error {
	jsonConfig, err := yaml.Marshal(dc)
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
