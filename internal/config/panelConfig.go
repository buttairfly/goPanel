package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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

func newPanelConfigFromPath(file string) (*PanelConfig, error) {
	pc := new(PanelConfig)
	err := pc.FromFile(file)
	return pc, err
}

// FromFile reads the config from a file at path
func (pc *PanelConfig) FromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can not read Config file %v. error: %v", path, err)
	}
	defer f.Close()
	return pc.FromReader(f)
}

// FromReader decodes the config from io.Reader
func (pc *PanelConfig) FromReader(r io.Reader) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(&*pc)
	if err != nil {
		return fmt.Errorf("can not decode json. error: %v", err)
	}
	return nil
}

// WriteToFile writes the config to a file at path
func (pc *PanelConfig) WriteToFile(path string) error {
	jsonConfig, err := json.MarshalIndent(pc, "", "\t")
	if err != nil {
		return err
	}
	jsonConfig = append(jsonConfig, byte('\n'))
	return ioutil.WriteFile(path, jsonConfig, 0622)
}
