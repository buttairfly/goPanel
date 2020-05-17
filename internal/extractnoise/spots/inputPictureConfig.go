package spots

import (
	"image"
	"io"
	"io/ioutil"
	"os"
	"sort"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// InputPictureConfig configures the mapping for incoming pictures
type InputPictureConfig struct {
	Offset     image.Point   `yaml:"offset"`
	TileWidth  int           `yaml:"tileWidth"`
	TileHeight int           `yaml:"tileHeight"`
	TileSpots  []image.Point `yaml:"tileSpots"`
	Height     int           `yaml:"height"`
	Width      int           `yaml:"width"`
}

// NewSpotsFromConfig creates a new InputPictureConfig from config file
func NewSpotsFromConfig(filePath string, logger *zap.Logger) (Spots, error) {
	var ipc InputPictureConfig
	err := ipc.FromYamlFile(filePath, logger)
	if err != nil {
		return nil, err
	}
	spots := ipc.ToSpots()
	return spots, nil
}

// FromYamlFile reads the config from a filePath
func (ipc *InputPictureConfig) FromYamlFile(filePath string, logger *zap.Logger) error {
	f, err := os.Open(filePath)
	if err != nil {
		logger.Error("can not read InputPictureConfig file", zap.String("configPath", filePath), zap.Error(err))
		return err
	}
	defer f.Close()
	return ipc.FromYamlReader(f, logger)
}

// FromYamlReader decodes the config from io.Reader
func (ipc *InputPictureConfig) FromYamlReader(r io.Reader, logger *zap.Logger) error {
	dec := yaml.NewDecoder(r)
	err := dec.Decode(&*ipc)
	if err != nil {
		logger.Error("can not decode InputPictureConfig yaml", zap.Error(err))
		return err
	}
	return nil
}

// WriteToYamlFile writes the config to a filePath
func (ipc *InputPictureConfig) WriteToYamlFile(filePath string) error {
	jsonConfig, err := yaml.Marshal(ipc)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, jsonConfig, 0622)
}

// ToSpots transforms a InputPictureConfig to spots stuct
func (ipc *InputPictureConfig) ToSpots() Spots {
	points := make([]image.Point, ipc.Height*ipc.Width*len(ipc.TileSpots))
	i := 0
	for y := 0; y < ipc.Height; y++ {
		for x := 0; x < ipc.Width; x++ {
			for _, spot := range ipc.TileSpots {
				points[i] = image.Point{
					X: ipc.Offset.X + x*ipc.TileWidth + spot.X,
					Y: ipc.Offset.Y + y*ipc.TileHeight + spot.Y,
				}
				i++
			}
		}
	}
	spots := NewSpots(ipc.Width, points)
	sort.Sort(spots)
	return spots
}
