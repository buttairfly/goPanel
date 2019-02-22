package hardware

import (
	"image"
	"image/color"
	"log"

	"github.com/buttairfly/goPanel/src/config"
)

// Frame is a hardware frame
type Frame interface {
	image.Image
	ToLedStripe() LedStripe
	GetSumHardwarePixel() int
	SetRGBA(x, y int, c color.RGBA)
	Set(x, y int, c color.Color)
	RGBAAt(x, y int) color.RGBA
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

// NewCopyFrameWithImage creates a new Frame with the reference of Tiles
// and sets
func NewCopyFrameWithImage(other Frame, picture *image.RGBA) Frame {
	return &frame{
		image:            picture,
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
				bufferPos := stripePos * NumBytePixel
				framePoint := tile.FramePoint(tilePoint)
				frameColor := f.RGBAAt(framePoint.X, framePoint.Y)
				log.Print(f, tilePoint, stripePos, framePoint, frameColor)
				buffer[bufferPos+R] = frameColor.R
				buffer[bufferPos+G] = frameColor.G
				buffer[bufferPos+B] = frameColor.B
			}
		}
	}
	log.Print(len(buffer), buffer)
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

// Set makes image buffer muteable
func (f *frame) Set(x, y int, c color.Color) {
	f.image.Set(x, y, c)
}

// SetRGBA makes image buffer muteable
func (f *frame) SetRGBA(x, y int, c color.RGBA) {
	f.image.SetRGBA(x, y, c)
}

func (f *frame) RGBAAt(x, y int) color.RGBA {
	return f.image.RGBAAt(x, y)
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