package hardware

import (
	"image"
	"image/color"
)

// Frame is a hardware frame
type Frame interface {
	image.Image
	ToLedStripe() LedStripe

	getTiles() []Tile
	getNumHardwarePixel() int
}

type frame struct {
	image            *image.RGBA
	tiles            []Tile
	numHardwarePixel int
}

// NewFrame return new Frame
func NewFrame(tileConfigs []TileConfig) Frame {
	frameBounds := image.ZR
	tiles := make([]Tile, len(tileConfigs))
	for _, tileConfig := range tileConfigs {
		frameBounds = frameBounds.Union(tileConfig.Bounds)
	}
	for i, tileConfig := range tileConfigs {
		tiles[i] = NewTile(tileConfig, frameBounds)
	}
	var numHardwarePixel int
	for _, tile := range tiles {
		numHardwarePixel += tile.NumPixel()
	}
	return &frame{
		image:            image.NewRGBA(frameBounds),
		tiles:            tiles,
		numHardwarePixel: numHardwarePixel,
	}
}

// NewCopyFrameWithEmptyImage creates a new Frame with the reference of Tiles
// but creates a new image
func NewCopyFrameWithEmptyImage(other Frame) Frame {
	return &frame{
		image:            image.NewRGBA(other.Bounds()),
		tiles:            other.getTiles(),
		numHardwarePixel: other.getNumHardwarePixel(),
	}
}

func (f *frame) ToLedStripe() LedStripe {
	return &ledStripe{}
}

// implements image interface
func (f *frame) ColorModel() color.Model {
	return f.image.ColorModel()
}

// implements image interface
func (f *frame) Bounds() image.Rectangle {
	return f.image.Bounds()
}

// implements image interface
func (f *frame) At(x, y int) color.Color {
	return f.image.At(x, y)
}

func (f *frame) getTiles() []Tile {
	return f.tiles
}

func (f *frame) getNumHardwarePixel() int {
	return f.numHardwarePixel
}
