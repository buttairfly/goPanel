package hardware

import (
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/buttairfly/goPanel/internal/intmath"
	"github.com/buttairfly/goPanel/pkg/filereadwriter"
)

// MapFormatString is the string to format the map[string]int led position mapping
const MapFormatString = "%2d"

// TileConfig is the config of a tile or led module
type TileConfig interface {
	filereadwriter.Yaml
	NumHardwarePixel() int
	GetBounds() image.Rectangle
	GetConnectionOrder() int
	GetLedStripeMap() map[string]int
}

type tileConfig struct {
	ConnectionOrder int             `yaml:"connectionOrder"`
	Bounds          image.Rectangle `yaml:"bounds"`
	LedStripeMap    map[string]int  `yaml:"ledStripeMap"`
}

// NewTileConfigFromPath creates a new tile from config file path
func NewTileConfigFromPath(path string) (TileConfig, error) {
	tc := new(tileConfig)
	err := tc.FromYamlFile(path)
	if err != nil {
		return nil, err
	}
	return tc, nil
}

// NumHardwarePixel counts the number of actual valid hardware pixels in the config
func (tc *tileConfig) NumHardwarePixel() int {
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
		log.Printf("numHardwarePixel (%d) of tile %d is not within max stripe pos %d", numHardwarePixel, tc.ConnectionOrder, maxStripePos)
	}
	return numHardwarePixel
}

// FromYamlFile reads the config from a file at path
func (tc *tileConfig) FromYamlFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can not read tileConfig file %v. error: %v", path, err)
	}
	defer f.Close()
	return tc.FromYamlReader(f)
}

// FromYamlReader decodes the config from io.Reader
func (tc *tileConfig) FromYamlReader(r io.Reader) error {
	dec := yaml.NewDecoder(r)

	log.Print("tileconfig fromReader")
	err := dec.Decode(&*tc)
	if err != nil {
		return fmt.Errorf("can not decode tileConfig yaml. error: %v", err)
	}
	return nil
}

// WriteToYamlFile writes the config to a file at path
func (tc *tileConfig) WriteToYamlFile(path string) error {
	yamlConfig, err := yaml.Marshal(tc)
	if err != nil {
		return err
	}

	yamlConfig = append(yamlConfig, byte('\n'))
	return ioutil.WriteFile(path, yamlConfig, 0622)
}

// Bounds retruns the tile image rectangle
func (tc *tileConfig) GetBounds() image.Rectangle {
	return tc.Bounds
}

// GetConnectionOrder retruns the tile connection order
func (tc *tileConfig) GetConnectionOrder() int {
	return tc.ConnectionOrder
}

// GetLedStripeMap retruns the tile led stripe map
func (tc *tileConfig) GetLedStripeMap() map[string]int {
	return tc.LedStripeMap
}

func tilePointxyToString(x, y, maxX int) string {
	return tilePositionToString(y*maxX + x)
}

func tilePositionToString(pos int) string {
	return fmt.Sprintf(MapFormatString, pos)
}
