package hardware

import (
	"image"
	"image/color"
	"sort"
)

// Frame is a hardware frame
type Frame interface {
	image.Image
	ToLedStripe() LedStripe

	getTiles() []Tile
	getSumHardwarePixel() int
}

type frame struct {
	image            *image.RGBA
	tiles            []Tile
	sumHardwarePixel int
}

// NewFrame return new Frame
func NewFrame(tileConfigs TileConfigs) Frame {
	frameBounds := image.ZR
	sort.Sort(tileConfigs)
	tiles := make([]Tile, tileConfigs.Len())
	numPreviousLedsOnStripe := 0
	for i, tileConfig := range tileConfigs.GetSlice() {
		frameBounds = frameBounds.Union(tileConfig.GetBounds())
		tiles[i] = NewTile(tileConfig, numPreviousLedsOnStripe)
		numPreviousLedsOnStripe += tileConfig.NumHardwarePixel()
	}
	return &frame{
		image:            image.NewRGBA(frameBounds),
		tiles:            tiles,
		sumHardwarePixel: numPreviousLedsOnStripe,
	}
}

// NewCopyFrameWithEmptyImage creates a new Frame with the reference of Tiles
// but creates a new image
func NewCopyFrameWithEmptyImage(other Frame) Frame {
	return &frame{
		image:            image.NewRGBA(other.Bounds()),
		tiles:            other.getTiles(),
		sumHardwarePixel: other.getSumHardwarePixel(),
	}
}

func (f *frame) ToLedStripe() LedStripe {
	buffer := make([]uint8, f.sumHardwarePixel*NumBytePixel)
	for _, tile := range f.tiles {
		for x := 0; x < tile.Bounds().Dx(); x++ {
			for y := 0; y < tile.Bounds().Dy(); y++ {
				tilePoint := image.Pt(x, y)
				stripePos := tile.MapTilePixelToStripePosition(tilePoint)
				framePoint := tile.FramePoint(tilePoint)
				frameColor := f.image.RGBAAt(framePoint.X, framePoint.Y)
				buffer[stripePos+R] = frameColor.R
				buffer[stripePos+G] = frameColor.G
				buffer[stripePos+B] = frameColor.B
			}
		}
	}
	return &ledStripe{
		buffer:      buffer,
		pixelLength: f.sumHardwarePixel,
	}
}

// ColorModel implements image interface
func (f *frame) ColorModel() color.Model {
	return f.image.ColorModel()
}

// Bounds implements image interface
func (f *frame) Bounds() image.Rectangle {
	return f.image.Bounds()
}

// At implements image interface
func (f *frame) At(x, y int) color.Color {
	return f.image.At(x, y)
}

func (f *frame) getTiles() []Tile {
	return f.tiles
}

func (f *frame) getSumHardwarePixel() int {
	return f.sumHardwarePixel
}
