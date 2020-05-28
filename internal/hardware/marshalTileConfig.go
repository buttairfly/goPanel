package hardware

import (
	"io"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/buttairfly/goPanel/pkg/marshal"
)

// MarshalTileConfig is a marshalable TileConfig
type MarshalTileConfig struct {
	ConnectionOrder int               `json:"connectionOrder" yaml:"connectionOrder"`
	Bounds          marshal.Rectangle `json:"bounds" yaml:"bounds"`
	LedStripeMap    map[string]int    `json:"ledStripeMap" yaml:"ledStripeMap"`
}

// NewTileConfigFromPath creates a new tile from config filePath
func NewTileConfigFromPath(filePath string, logger *zap.Logger) (*MarshalTileConfig, error) {
	mtc := new(MarshalTileConfig)
	err := mtc.FromYamlFile(filePath, logger)
	if err != nil {
		return nil, err
	}
	return mtc, nil
}

// FromYamlFile reads the config from a filePath
func (mtc *MarshalTileConfig) FromYamlFile(filePath string, logger *zap.Logger) error {
	f, err := os.Open(filePath)
	if err != nil {
		logger.Error("can not read TileConfig file", zap.String("configPath", filePath), zap.Error(err))
		return err
	}
	defer f.Close()
	return mtc.FromYamlReader(f, logger)
}

// FromYamlReader decodes the config from io.Reader
func (mtc *MarshalTileConfig) FromYamlReader(r io.Reader, logger *zap.Logger) error {
	dec := yaml.NewDecoder(r)
	err := dec.Decode(mtc)
	if err != nil {
		logger.Error("can not decode TileConfig yaml", zap.Error(err))
		return err
	}
	return nil
}

// WriteToYamlFile writes the config to a filePath
func (mtc *MarshalTileConfig) WriteToYamlFile(filePath string) error {
	yamlConfig, err := yaml.Marshal(mtc)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, yamlConfig, 0622)
}

// ToTileConfig retruns the TileConfig
func (mtc *MarshalTileConfig) ToTileConfig() *TileConfig {
	tc := new(TileConfig)
	tc.ConnectionOrder = mtc.ConnectionOrder
	tc.Bounds = mtc.Bounds.ToImageRectangle()
	tc.LedStripeMap = mtc.LedStripeMap
	return tc
}
