package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"

	"github.com/buttairfly/goPanel/internal/device"
	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/pkg/arduinocom"
	"github.com/buttairfly/goPanel/pkg/common"
)

// Config is the internal full config
type Config interface {
	common.JSONFileReadWriter
	json.Unmarshaler

	GetTileConfigs() hardware.TileConfigs
	GetDeviceConfig() *device.DeviceConfig
}

type config struct {
	TileConfigs  hardware.TileConfigs `json:"tileConfigs"`
	DeviceConfig *device.DeviceConfig `json:"deviceConfig"`
}

// NewConfigFromPanelConfigPath generates a new internal config struct from panel config file
func NewConfigFromPanelConfigPath(folderOffset, path string) (Config, error) {
	panelConfig, err := newPanelConfigFromPath(folderOffset, path)
	if err != nil {
		return nil, err
	}

	tileConfigs := make(hardware.TileConfigSlice, len(panelConfig.TileConfigPaths))
	for i, tileConfigPath := range panelConfig.TileConfigPaths {
		tileConfigs[i], err = hardware.NewTileConfigFromPath(folderOffset + tileConfigPath)
		if err != nil {
			return nil, err
		}
	}
	sort.Sort(tileConfigs)

	var deviceConfig *device.DeviceConfig
	deviceConfig, err = device.NewDeviceConfigFromPath(folderOffset + panelConfig.DeviceConfigPath)
	if err != nil {
		return nil, err
	}

	if deviceConfig.Type == device.Serial {
		var arduinoErrorConfig *arduinocom.ArduinoErrorConfig
		arduinoErrorConfig, err = arduinocom.NewArduinoErrorConfigFromPath(folderOffset + panelConfig.ArduinoErrorConfigPath)
		if err != nil {
			return nil, err
		}
		deviceConfig.SerialConfig.ArduinoErrorConfig = arduinoErrorConfig
	}

	return &config{
		TileConfigs:  tileConfigs,
		DeviceConfig: deviceConfig,
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

func (c *config) GetTileConfigs() hardware.TileConfigs {
	return c.TileConfigs
}

func (c *config) GetDeviceConfig() *device.DeviceConfig {
	return c.DeviceConfig
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

	if objMap["deviceConfig"] != nil {
		err = json.Unmarshal(*objMap["deviceConfig"], &c.DeviceConfig)
		if err != nil {
			return err
		}
	} else {
		return errors.New("No DeviceConfig in config file")
	}

	if objMap["tileConfigs"] != nil {
		var rawMessagesTileConfigsJSON []*json.RawMessage
		err = json.Unmarshal(*objMap["tileConfigs"], &rawMessagesTileConfigsJSON)
		if err != nil {
			return err
		}

		c.TileConfigs = make(hardware.TileConfigSlice, len(rawMessagesTileConfigsJSON))

		for i, rawMessage := range rawMessagesTileConfigsJSON {
			var tileConfig hardware.TileConfig
			err = json.Unmarshal(*rawMessage, tileConfig)
			if err != nil {
				return err
			}
			c.TileConfigs.Set(i, tileConfig)
		}
		sort.Sort(c.TileConfigs)
	} else {
		return errors.New("No TileConfigs in config file")
	}
	return nil
}
