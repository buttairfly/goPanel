package config

import (
	"io"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// PanelConfig is the global panel config
type PanelConfig struct {
	TileConfigPath         string   `json:"tileConfigPath"`
	DeviceConfigPath       string   `json:"deviceConfigPath"`
	ArduinoErrorConfigPath string   `json:"arduinoErrorConfigPath,omitempty"`
	TileConfigFiles        []string `json:"tileConfigFiles"`
	DeviceConfigFile       string   `json:"deviceConfigFile"`
	ArduinoErrorConfigFile string   `json:"arduinoErrorConfigFile,omitempty"`
}

func newPanelConfigFromPath(filePath string, logger *zap.Logger) (*PanelConfig, error) {
	pc := new(PanelConfig)
	err := pc.FromYamlFile(filePath, logger)
	return pc, err
}

// FromYamlFile reads the config from filePath
func (pc *PanelConfig) FromYamlFile(filePath string, logger *zap.Logger) error {
	f, err := os.Open(filePath)
	if err != nil {
		logger.Error("can not read panelConfig file", zap.String("configPath", filePath), zap.Error(err))
		return err
	}
	defer f.Close()
	return pc.FromYamlReader(f, logger)
}

// FromYamlReader decodes the config from io.Reader
func (pc *PanelConfig) FromYamlReader(r io.Reader, logger *zap.Logger) error {
	dec := yaml.NewDecoder(r)
	err := dec.Decode(&*pc)
	if err != nil {
		logger.Error("can not decode panelConfig yaml", zap.Error(err))
		return err
	}
	return nil
}

// WriteToYamlFile writes the config to filePath
func (pc *PanelConfig) WriteToYamlFile(filePath string) error {
	yamlConfig, err := yaml.Marshal(pc)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, yamlConfig, 0622)
}
