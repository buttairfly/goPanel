package config

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/buttairfly/goPanel/src/intmath"
)

// TileConfig is the config of a tile or led module
type TileConfig interface {
	NumHardwarePixel() int
	GetBounds() image.Rectangle
	GetConnectionOrder() int
	GetLedStripeMap() map[string]int
	FromReader(r io.Reader) error
	FromFile(path string) error
	WriteToFile(path string) error
}

type tileConfig struct {
	ConnectionOrder int             `json:"connectionOrder"`
	Bounds          image.Rectangle `json:"bounds"`
	LedStripeMap    map[string]int  `json:"ledStripeMap"`
}

// NewTileConfigFromPath creates a new tile from config file path
func NewTileConfigFromPath(path string) (TileConfig, error) {
	tc := new(tileConfig)
	err := tc.FromFile(path)
	if err != nil {
		return nil, err
	}
	return tc, nil
}

// NumHardwarePixel counts the number of actual valid hardware pixels in the config
func (tc *tileConfig) NumHardwarePixel() int {
	maxX := tc.Bounds.Dx() + 1
	maxY := tc.Bounds.Dy() + 1
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
	if maxStripePos != numHardwarePixel-1 {
		log.Printf("numHardwarePixel of %d tile %d is not within max stripe pos %d", numHardwarePixel, tc.ConnectionOrder, maxStripePos)
	}
	return numHardwarePixel
}

// FromFile reads the config from a file at path
func (tc *tileConfig) FromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can not read Config file %v. error: %v", path, err)
	}
	defer f.Close()
	return tc.FromReader(f)
}

// FromReader decodes the config from io.Reader
func (tc *tileConfig) FromReader(r io.Reader) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(&*tc)
	if err != nil {
		return fmt.Errorf("can not decode json. error: %v", err)
	}
	return nil
}

// WriteToFile writes the config to a file at path
func (tc *tileConfig) WriteToFile(path string) error {
	jsonConfig, err := json.MarshalIndent(tc, "", "\t")
	if err != nil {
		return err
	}
	jsonConfig = append(jsonConfig, byte('\n'))
	return ioutil.WriteFile(path, jsonConfig, 0622)
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
	return fmt.Sprintf("%04d", pos)
}
