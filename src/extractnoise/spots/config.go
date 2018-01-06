package spots

import (
	"image"
	"os"
)

type InputPictureConfig struct {
	Offset     image.Point   `json:"offset"`
	TileWidth  int           `json:"tileWidth"`
	TileHeight int           `json:"tileHeight"`
	TileSpots  []image.Point `json:"tileSpots"`
	Height     int           `json:"height"`
	Width      int           `json:"width"`
}

func DecodeConfig(filename string) (*InputPictureConfig, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return nil, nil
}

func (ipc *InputPictureConfig) ToSpots() Spots {
	points := make([]image.Point, 0)
	for y := 0; y < ipc.Height; y++ {
		for x := 0; x < ipc.Width; x++ {
			for spot := range ipc.TileSpots {
				p := image.Point{
					X: ipc.Offset.X + x*ipc.TileWidth + spot.X,
					Y: ipc.Offset.Y + y*ipc.TileHeight + spot.Y,
				}
				points = append(spots, p)
			}
		}
	}
	spots := new(spots)
	spots.Width = ipc.Width
	spots.spots = points
	spots.Sort()
	return spots
}
