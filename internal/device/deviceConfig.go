package device

import (
	"io"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/buttairfly/goPanel/pkg/arduinocom"
)

// DeviceConfig is the config of the type of device
type DeviceConfig struct {
	Type         Type                     `yaml:"type"`
	SerialConfig *arduinocom.SerialConfig `yaml:"serialConfig,omitempty"`
}

// NewDeviceConfigFromPath returns a new DeviceConfig or error
func NewDeviceConfigFromPath(path string, logger *zap.Logger) (*DeviceConfig, error) {
	dc := new(DeviceConfig)
	err := dc.FromYamlFile(path, logger)
	if err != nil {
		return nil, err
	}
	return dc, nil
}

// FromYamlFile reads the config from a file at path
func (dc *DeviceConfig) FromYamlFile(path string, logger *zap.Logger) error {
	f, err := os.Open(path)
	if err != nil {
		logger.Error("can not read DeviceConfig file", zap.String("configPath", path), zap.Error(err))
		return err
	}
	defer f.Close()
	return dc.FromYamlReader(f, logger)
}

// FromYamlReader decodes the config from io.Reader
func (dc *DeviceConfig) FromYamlReader(r io.Reader, logger *zap.Logger) error {
	dec := yaml.NewDecoder(r)
	err := dec.Decode(&*dc)
	if err != nil {
		logger.Error("can not decode DeviceConfig yaml", zap.Error(err))
		return err
	}
	return nil
}

// WriteToYamlFile writes the config to a file at path
func (dc *DeviceConfig) WriteToYamlFile(path string) error {
	yamlConfig, err := yaml.Marshal(dc)
	if err != nil {
		return err
	}
	yamlConfig = append(yamlConfig, byte('\n'))
	return ioutil.WriteFile(path, yamlConfig, 0622)
}

// Type is a LedDevice type
type Type string

const (
	// Print debug print device
	Print = Type("print")
	// Serial high level serial tty device
	Serial = Type("serial")
)
