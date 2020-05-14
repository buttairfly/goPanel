package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"

	"gopkg.in/yaml.v2"

	"github.com/buttairfly/goPanel/internal/device"
	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/pkg/arduinocom"
	"github.com/buttairfly/goPanel/pkg/filereadwriter"
)

// Config is the internal full config
type Config interface {
	filereadwriter.Yaml
	//yaml.Unmarshaler

	GetTileConfigs() hardware.TileConfigs
	GetDeviceConfig() *device.DeviceConfig
}

type config struct {
	TileConfigs  hardware.TileConfigs `yaml:"tileConfigs"`
	DeviceConfig *device.DeviceConfig `yaml:"deviceConfig"`
}

// NewConfigFromPanelConfigPath generates a new internal config struct from panel config file
func NewConfigFromPanelConfigPath(file string) (Config, error) {
	panelConfig, err := newPanelConfigFromPath(file)
	if err != nil {
		return nil, err
	}

	tileConfigs := make(hardware.TileConfigSlice, len(panelConfig.TileConfigFiles))
	for i, tileConfigFile := range panelConfig.TileConfigFiles {
		tileConfigs[i], err = hardware.NewTileConfigFromPath(path.Join(panelConfig.TileConfigPath, tileConfigFile))
		if err != nil {
			return nil, err
		}
	}
	sort.Sort(tileConfigs)

	var deviceConfig *device.DeviceConfig
	deviceConfig, err = device.NewDeviceConfigFromPath(path.Join(panelConfig.DeviceConfigPath, panelConfig.DeviceConfigFile))
	if err != nil {
		return nil, err
	}

	if deviceConfig.Type == device.Serial {
		var arduinoErrorConfig *arduinocom.ArduinoErrorConfig
		arduinoErrorConfig, err = arduinocom.NewArduinoErrorConfigFromPath(path.Join(panelConfig.ArduinoErrorConfigPath, panelConfig.ArduinoErrorConfigFile))
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
	err := c.FromYamlFile(path)
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

// FromYamlFile reads the config from a file at path
func (c *config) FromYamlFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can not read Config file %v. error: %v", path, err)
	}
	defer f.Close()
	return c.FromYamlReader(f)
}

// FromYamlReader decodes the config from io.Reader
func (c *config) FromYamlReader(r io.Reader) error {
	dec := yaml.NewDecoder(r)

	log.Print("config FromReader")
	err := dec.Decode(&*c)
	if err != nil {
		return fmt.Errorf("can not decode main config yaml. error: %v", err)
	}
	return nil
}

// WriteToYamlFile writes the config to a file at path
func (c *config) WriteToYamlFile(path string) error {
	yamlConfig, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	yamlConfig = append(yamlConfig, byte('\n'))
	return ioutil.WriteFile(path, yamlConfig, 0622)
}

/*
// UnmarshalYaml unmarshals a yaml file
func (c *config) UnmarshalYAML(b []byte) error {
	var objMap yaml.MapSlice
	err := yaml.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	if objMap["deviceConfig"] != nil {
		err = yaml.Unmarshal(*objMap["deviceConfig"], &c.DeviceConfig)
		if err != nil {
			return fmt.Errorf("deviceConfig error: %s", err)
		}
	} else {
		return fmt.Errorf("No DeviceConfig in config file")
	}

	if objMap["tileConfigs"] != nil {
		var rawMessagesTileConfigsYaml []*yaml.RawMessage
		err = yaml.Unmarshal(*objMap["tileConfigs"], &rawMessagesTileConfigsYaml)
		if err != nil {
			return fmt.Errorf("tileConfigs error: %s", err)
		}

		c.TileConfigs = make(hardware.TileConfigSlice, len(rawMessagesTileConfigsYaml))

		for i, rawMessage := range rawMessagesTileConfigsYaml {
			var tileConfig hardware.TileConfig
			err = yaml.Unmarshal(*rawMessage, tileConfig)
			if err != nil {
				return fmt.Errorf("tileConfig %d\nrawMessage: %s\ntileConfig: %T\nerror: %s", i, *rawMessage, tileConfig, err)
			}
			c.TileConfigs.Set(i, tileConfig)
		}
		sort.Sort(c.TileConfigs)
	} else {
		return errors.New("No TileConfigs in main config file")
	}
	return nil
}

*/
