package device

import (
	"io"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/buttairfly/goPanel/pkg/arduinocom"
)

// LedDeviceConfig is the config of the type of device
type LedDeviceConfig struct {
	Type         Type                     `json:"type" yaml:"type"`
	SerialConfig *arduinocom.SerialConfig `json:"serialConfig,omitempty" yaml:"serialConfig,omitempty"`
	PrintConfig  *PrintConfig             `json:"printConfig,omitempty" yaml:"printConfig,omitempty"`
}

// NewDeviceConfigFromPath returns a new LedDeviceConfig or error
func NewDeviceConfigFromPath(filePath string, logger *zap.Logger) (*LedDeviceConfig, error) {
	dc := new(LedDeviceConfig)
	err := dc.FromYamlFile(filePath, logger)
	if err != nil {
		return nil, err
	}
	return dc, nil
}

// FromYamlFile reads the config from filePath
func (dc *LedDeviceConfig) FromYamlFile(filePath string, logger *zap.Logger) error {
	f, err := os.Open(filePath)
	if err != nil {
		logger.Error("can not read LedDeviceConfig file", zap.String("configPath", filePath), zap.Error(err))
		return err
	}
	defer f.Close()
	return dc.FromYamlReader(f, logger)
}

// FromYamlReader decodes the config from io.Reader
func (dc *LedDeviceConfig) FromYamlReader(r io.Reader, logger *zap.Logger) error {
	dec := yaml.NewDecoder(r)
	err := dec.Decode(&*dc)
	if err != nil {
		logger.Error("can not decode LedDeviceConfig yaml", zap.Error(err))
		return err
	}
	return nil
}

// WriteToYamlFile writes the config to filePath
func (dc *LedDeviceConfig) WriteToYamlFile(filePath string) error {
	yamlConfig, err := yaml.Marshal(dc)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, yamlConfig, 0622)
}

// Type is a LedDevice type
type Type string

const (
	// Print debug print device
	Print = Type("print")
	// Serial high level serial tty device
	Serial = Type("serial")
)
