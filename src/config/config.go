package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
)

// Config is the internal full config
type Config interface {
	FromFile(path string) error
	FromReader(r io.Reader) error
	WriteToFile(path string) error
	GetTileConfigs() TileConfigs
	json.Unmarshaler
}

type config struct {
	TileConfigs TileConfigs `json:"tileConfigs"`
	// DeviceCondig device.Config
}

// NewConfigFromPanelConfigPath generates a new internal config struct from panel config file
func NewConfigFromPanelConfigPath(path string) (Config, error) {
	panelConfig, err := newPanelConfigFromPath(path)
	if err != nil {
		return nil, err
	}
	tileConfigs := make(tileConfigs, len(panelConfig.TileConfigPaths))
	for i, tileConfigPath := range panelConfig.TileConfigPaths {
		tileConfigs[i], err = NewTileConfigFromPath(tileConfigPath)
		if err != nil {
			return nil, err
		}
	}
	sort.Sort(tileConfigs)
	return &config{
		TileConfigs: tileConfigs,
	}, nil
}

func newConfigFromPath(path string) (Config, error) {
	c := new(config)
	err := c.FromFile(path)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *config) GetTileConfigs() TileConfigs {
	return c.TileConfigs
}

// FromFile reads the config from a file at path
func (c *config) FromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can not read Config file %v. error: %v", path, err)
	}
	defer f.Close()
	return c.FromReader(f)
}

// FromReader decodes the config from io.Reader
func (c *config) FromReader(r io.Reader) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(&*c)
	if err != nil {
		return fmt.Errorf("can not decode json. error: %v", err)
	}
	return nil
}

// WriteToFile writes the config to a file at path
func (c *config) WriteToFile(path string) error {
	jsonConfig, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	jsonConfig = append(jsonConfig, byte('\n'))
	return ioutil.WriteFile(path, jsonConfig, 0622)
}

// UnmarshalJSON unmarshals
func (c *config) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var rawMessagesTileConfigsJSON []*json.RawMessage
	err = json.Unmarshal(*objMap["tileConfigs"], &rawMessagesTileConfigsJSON)
	if err != nil {
		return err
	}

	c.TileConfigs = make(tileConfigs, len(rawMessagesTileConfigsJSON))

	for i, rawMessage := range rawMessagesTileConfigsJSON {
		var tileConfig tileConfig
		err = json.Unmarshal(*rawMessage, &tileConfig)
		if err != nil {
			return err
		}
		c.TileConfigs.Set(i, &tileConfig)
	}
	sort.Sort(c.TileConfigs)
	return nil
}
