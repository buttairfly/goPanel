package spots

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"
	"sort"
)

type InputPictureConfig struct {
	Offset     image.Point   `json:"offset"`
	TileWidth  int           `json:"tileWidth"`
	TileHeight int           `json:"tileHeight"`
	TileSpots  []image.Point `json:"tileSpots"`
	Height     int           `json:"height"`
	Width      int           `json:"width"`
}

func NewSpotsFromConfig(path string) (Spots, error) {
	var ipc InputPictureConfig
	err := ipc.FromFile(path)
	if err != nil {
		return nil, err
	}
	spots := ipc.ToSpots()
	return spots, nil
}

func (ipc *InputPictureConfig) FromReader(r io.Reader) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(&*ipc)
	if err != nil {
		return fmt.Errorf("can not decode json. error: %v", err)
	}
	return nil
}

func (ipc *InputPictureConfig) FromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can not read Config file %v. error: %v", path, err)
	}
	defer f.Close()
	return ipc.FromReader(f)
}

func (ipc *InputPictureConfig) WriteToFile(path string) error {
	jsonConfig, err := json.MarshalIndent(ipc, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, jsonConfig, 0622)
}

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
