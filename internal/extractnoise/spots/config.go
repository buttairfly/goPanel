package spots

import (
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"
	"sort"

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
func NewSpotsFromConfig(path string) (Spots, error) {
	var ipc InputPictureConfig
	err := ipc.FromYamlFile(path)
	if err != nil {
		return nil, err
	}
	spots := ipc.ToSpots()
	return spots, nil
}

// FromYamlReader decodes the config from io.Reader
func (ipc *InputPictureConfig) FromYamlReader(r io.Reader) error {
	dec := yaml.NewDecoder(r)
	err := dec.Decode(&*ipc)
	if err != nil {
		return fmt.Errorf("can not decode json. error: %v", err)
	}
	return nil
}

// FromYamlFile reads the config from a file at path
func (ipc *InputPictureConfig) FromYamlFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can not read Config file %v. error: %v", path, err)
	}
	defer f.Close()
	return ipc.FromYamlReader(f)
}

// WriteToYamlFile writes the config to a file at path
func (ipc *InputPictureConfig) WriteToYamlFile(path string) error {
	jsonConfig, err := yaml.Marshal(ipc)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, jsonConfig, 0622)
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
