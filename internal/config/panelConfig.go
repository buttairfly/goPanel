package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
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

func newPanelConfigFromPath(file string, logger *zap.Logger) (*PanelConfig, error) {
	pc := new(PanelConfig)
	err := pc.FromYamlFile(file, logger)
	return pc, err
}

// FromYamlFile reads the config from a file at path
func (pc *PanelConfig) FromYamlFile(path string, logger *zap.Logger) error {
	f, err := os.Open(path)
	if err != nil {
		logger.Error("can not read panelConfig file", zap.String("configPath", path), zap.Error(err))
		return err
	}
	defer f.Close()
	return pc.FromYamlReader(f, logger)
}

// FromYamlReader decodes the config from io.Reader
func (pc *PanelConfig) FromYamlReader(r io.Reader, logger *zap.Logger) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(&*pc)
	if err != nil {
		logger.Error("can not decode panelConfig yaml", zap.Error(err))
		return err
	}
	return nil
}

// WriteToYamlFile writes the config to a file at path
func (pc *PanelConfig) WriteToYamlFile(path string) error {
	jsonConfig, err := json.MarshalIndent(pc, "", "\t")
	if err != nil {
		return err
	}
	jsonConfig = append(jsonConfig, byte('\n'))
	return ioutil.WriteFile(path, jsonConfig, 0622)
}
