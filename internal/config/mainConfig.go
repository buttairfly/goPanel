package config

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/buttairfly/goPanel/internal/device"
	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/pkg/arduinocom"
)

// MainConfig is the whole program config
type MainConfig struct {
	TileConfigs  hardware.TileConfigs `yaml:"tileConfigs"`
	DeviceConfig *device.DeviceConfig `yaml:"deviceConfig"`
}

// NewMainConfigFromPanelConfigPath generates a new internal MainConfig struct from PanelConfig file
func NewMainConfigFromPanelConfigPath(filePath string, logger *zap.Logger) (*MainConfig, error) {
	panelConfig, err := newPanelConfigFromPath(filePath, logger)
	if err != nil {
		return nil, err
	}

	tileConfigs := make(hardware.TileConfigs, len(panelConfig.TileConfigFiles))
	for i, tileConfigFile := range panelConfig.TileConfigFiles {
		tileConfigs[i], err = hardware.NewTileConfigFromPath(
			path.Join(panelConfig.TileConfigPath, tileConfigFile),
			logger,
		)
		if err != nil {
			return nil, err
		}
	}
	sort.Sort(tileConfigs)

	var deviceConfig *device.DeviceConfig
	deviceConfig, err = device.NewDeviceConfigFromPath(
		path.Join(panelConfig.DeviceConfigPath, panelConfig.DeviceConfigFile),
		logger,
	)
	if err != nil {
		return nil, err
	}

	if deviceConfig.Type == device.Serial {
		var arduinoErrorConfig *arduinocom.ArduinoErrorConfig
		arduinoErrorConfig, err = arduinocom.NewArduinoErrorConfigFromPath(
			path.Join(panelConfig.ArduinoErrorConfigPath, panelConfig.ArduinoErrorConfigFile),
			logger,
		)
		if err != nil {
			return nil, err
		}
		deviceConfig.SerialConfig.ArduinoErrorConfig = arduinoErrorConfig
	}

	return &MainConfig{
		TileConfigs:  tileConfigs,
		DeviceConfig: deviceConfig,
	}, nil
}

// NewMainConfigFromPath gets a MainConfig from filePath
func NewMainConfigFromPath(filePath string, logger *zap.Logger) (*MainConfig, error) {
	c := new(MainConfig)
	err := c.FromYamlFile(filePath, logger)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// FromYamlFile reads the config from filePath
func (c *MainConfig) FromYamlFile(filePath string, logger *zap.Logger) error {
	f, err := os.Open(filePath)
	if err != nil {
		logger.Error("can not read config file", zap.String("configPath", filePath), zap.Error(err))
		return err
	}
	defer f.Close()
	return c.FromYamlReader(f, logger)
}

// FromYamlReader decodes the config from io.Reader
func (c *MainConfig) FromYamlReader(r io.Reader, logger *zap.Logger) error {
	dec := yaml.NewDecoder(r)
	err := dec.Decode(&*c)
	if err != nil {
		logger.Error("can not decode panelConfig yaml", zap.Error(err))
		return err
	}
	return nil
}

// WriteToYamlFile writes the config to filePath
func (c *MainConfig) WriteToYamlFile(filePath string) error {
	yamlConfig, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, yamlConfig, 0622)
}
