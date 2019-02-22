package hardware

import (
	"image"
	"image/color"

	"github.com/buttairfly/goPanel/src/config"
)

// Frame is a hardware frame
type Frame interface {
	image.Image
	ToLedStripe() LedStripe

	GetSumHardwarePixel() int
	SetRGBA(x, y int, c color.RGBA)
	GetWidth() int
	GetHeight() int

	getTiles() []Tile
}

type frame struct {
	image            *image.RGBA
	tiles            []Tile
	sumHardwarePixel int
	width, height    int
}

// NewFrame return new Frame
func NewFrame(tileConfigs config.TileConfigs) Frame {
	frameBounds := image.ZR
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
		width:            frameBounds.Dx() + 1,
		height:           frameBounds.Dy() + 1,
	}
}

// NewCopyFrameWithEmptyImage creates a new Frame with the reference of Tiles
// but creates a new image
func NewCopyFrameWithEmptyImage(other Frame) Frame {
	return &frame{
		image:            image.NewRGBA(other.Bounds()),
		tiles:            other.getTiles(),
		sumHardwarePixel: other.GetSumHardwarePixel(),
		width:            other.GetWidth(),
		height:           other.GetHeight(),
	}
}

func (f *frame) ToLedStripe() LedStripe {
	buffer := make([]uint8, f.sumHardwarePixel*NumBytePixel)
	for _, tile := range f.tiles {
		for x := 0; x < tile.GetWidth(); x++ {
			for y := 0; y < tile.GetHeight(); y++ {
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

// SetRGBA makes image buffer muteable
func (f *frame) SetRGBA(x, y int, c color.RGBA) {
	f.image.SetRGBA(x, y, c)
}

func (f *frame) GetSumHardwarePixel() int {
	return f.sumHardwarePixel
}

func (f *frame) GetWidth() int {
	return f.width
}

func (f *frame) GetHeight() int {
	return f.height
}

func (f *frame) getTiles() []Tile {
	return f.tiles
}
