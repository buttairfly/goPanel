package device

import (
	"io"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// PrintConfig configures a print device
type PrintConfig struct {
	FramesPerSecond int  `yaml:"framesPerSecond"`
	Quiet           bool `yaml:"quiet,omitempty"`
}

// NewPrintConfigFromPath returns a new PrintConfig or error
func NewPrintConfigFromPath(filePath string, logger *zap.Logger) (*PrintConfig, error) {
	dc := new(PrintConfig)
	err := dc.FromYamlFile(filePath, logger)
	if err != nil {
		return nil, err
	}
	return dc, nil
}

// FromYamlFile reads the config from filePath
func (dc *PrintConfig) FromYamlFile(filePath string, logger *zap.Logger) error {
	f, err := os.Open(filePath)
	if err != nil {
		logger.Error("can not read PrintConfig file", zap.String("configPath", filePath), zap.Error(err))
		return err
	}
	defer f.Close()
	return dc.FromYamlReader(f, logger)
}

// FromYamlReader decodes the config from io.Reader
func (dc *PrintConfig) FromYamlReader(r io.Reader, logger *zap.Logger) error {
	dec := yaml.NewDecoder(r)
	err := dec.Decode(&*dc)
	if err != nil {
		logger.Error("can not decode PrintConfig yaml", zap.Error(err))
		return err
	}
	return nil
}

// WriteToYamlFile writes the config to filePath
func (dc *PrintConfig) WriteToYamlFile(filePath string) error {
	yamlConfig, err := yaml.Marshal(dc)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, yamlConfig, 0622)
}
