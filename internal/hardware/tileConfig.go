package hardware

import (
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/buttairfly/goPanel/internal/intmath"
	"github.com/buttairfly/goPanel/pkg/marshal"
)

// MapFormatString is the string to format the map[string]int led position mapping
const MapFormatString = "%2d"

// TileConfig is a struct of a config of one led panel tile
type TileConfig struct {
	ConnectionOrder int             `json:"connectionOrder" yaml:"connectionOrder"`
	Bounds          image.Rectangle `json:"bounds" yaml:"bounds"`
	LedStripeMap    map[string]int  `json:"ledStripeMap" yaml:"ledStripeMap"`
}

// MarshalTileConfig is a marshalable TileConfig
type MarshalTileConfig struct {
	ConnectionOrder int               `json:"connectionOrder" yaml:"connectionOrder"`
	Bounds          marshal.Rectangle `json:"bounds" yaml:"bounds"`
	LedStripeMap    map[string]int    `json:"ledStripeMap" yaml:"ledStripeMap"`
}

// NewTileConfigFromPath creates a new tile from config filePath
func NewTileConfigFromPath(filePath string, logger *zap.Logger) (*TileConfig, error) {
	mtc := new(MarshalTileConfig)
	err := mtc.FromYamlFile(filePath, logger)
	if err != nil {
		return nil, err
	}
	return mtc.ToTileConfig(), nil
}

// NumHardwarePixel counts the number of actual valid hardware pixels in the config
func (tc *TileConfig) NumHardwarePixel() int {
	maxX := tc.Bounds.Dx()
	maxY := tc.Bounds.Dy()
	maxPixel := maxX * maxY
	numHardwarePixel := 0
	maxStripePos := 0
	for tilePos := 0; tilePos < maxPixel; tilePos++ {
		stripePos, ok := tc.LedStripeMap[tilePositionToString(tilePos)]
		stripePoint := image.Point{X: tilePos % maxX, Y: tilePos / maxY}
		maxStripePos = intmath.Max(stripePos, maxStripePos)
		if ok &&
			stripePos >= 0 &&
			stripePoint.X >= 0 &&
			stripePoint.Y >= 0 &&
			stripePoint.X < maxX &&
			stripePoint.Y < maxY {
			numHardwarePixel++
		}
	}
	if maxStripePos > numHardwarePixel-1 {
		zap.L().Sugar().Infof(
			"numHardwarePixel (%d) of tile %d is not within max stripe pos %d",
			numHardwarePixel,
			tc.ConnectionOrder,
			maxStripePos,
		)
	}
	return numHardwarePixel
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
func (tc *TileConfig) WriteToYamlFile(filePath string) error {
	yamlConfig, err := yaml.Marshal(tc.ToMarshalTileConfig())
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, yamlConfig, 0622)
}

// ToMarshalTileConfig retruns the marshalable tile config
func (tc *TileConfig) ToMarshalTileConfig() *MarshalTileConfig {
	mtc := new(MarshalTileConfig)
	mtc.ConnectionOrder = tc.ConnectionOrder
	mtc.Bounds = marshal.FromImageRectangle(tc.Bounds)
	mtc.LedStripeMap = tc.LedStripeMap
	return mtc
}

// ToTileConfig converts the marshalable tile config into a TileConfig
func (mtc *MarshalTileConfig) ToTileConfig() *TileConfig {
	tc := new(TileConfig)
	tc.ConnectionOrder = mtc.ConnectionOrder
	tc.Bounds = mtc.Bounds.ToImageRectangle()
	tc.LedStripeMap = mtc.LedStripeMap
	return tc
}

// GetBounds retruns the tile image rectangle
func (tc *TileConfig) GetBounds() image.Rectangle {
	return tc.Bounds
}

// GetConnectionOrder retruns the tile connection order
func (tc *TileConfig) GetConnectionOrder() int {
	return tc.ConnectionOrder
}

// GetLedStripeMap retruns the tile led stripe map
func (tc *TileConfig) GetLedStripeMap() map[string]int {
	return tc.LedStripeMap
}

func tilePointxyToString(x, y, maxX int) string {
	return tilePositionToString(y*maxX + x)
}

func tilePositionToString(pos int) string {
	return fmt.Sprintf(MapFormatString, pos)
}
