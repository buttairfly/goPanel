package hardware

import (
	"image"
	"image/color"
	"time"

	"go.uber.org/zap"
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
	GetTime() time.Time
	GetLogger() *zap.Logger

	getTiles() []Tile
}

type frame struct {
	image            *image.RGBA
	tiles            []Tile
	sumHardwarePixel int
	width, height    int
	frameTime        time.Time
	logger           *zap.Logger
}

// NewFrame returns a paintable image consisting of many led-panel tiles
func NewFrame(tileConfigs TileConfigs, logger *zap.Logger) Frame {
	frameBounds := image.ZR
	tiles := make([]Tile, tileConfigs.Len())
	numPreviousLedsOnStripe := 0
	for i, tileConfig := range tileConfigs {
		frameBounds = frameBounds.Union(tileConfig.GetBounds())
		tiles[i] = NewTile(tileConfig, numPreviousLedsOnStripe)
		numPreviousLedsOnStripe += tileConfig.NumHardwarePixel()
	}
	frameBounds.Canon()
	return &frame{
		image:            image.NewRGBA(frameBounds),
		tiles:            tiles,
		sumHardwarePixel: numPreviousLedsOnStripe,
		width:            frameBounds.Dx(),
		height:           frameBounds.Dy(),
		frameTime:        time.Now(),
		logger:           logger,
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
		frameTime:        time.Now(),
		logger:           other.GetLogger(),
	}
}

// NewCopyFrameFromImage creates a new Frame with the reference of Tiles
// and copies the other image contents into the frame
func NewCopyFrameFromImage(other Frame, pictureToCopy *image.RGBA) Frame {
	if !other.Bounds().Eq(pictureToCopy.Bounds()) {
		other.GetLogger().Sugar().Fatalf("can not copy picture (%v) with different bounds as frame (%v)",
			other.Bounds(),
			pictureToCopy.Bounds(),
		)
	}
	picture := image.NewRGBA(pictureToCopy.Bounds())
	copy(picture.Pix, pictureToCopy.Pix)
	return &frame{
		image:            picture,
		tiles:            other.getTiles(),
		sumHardwarePixel: other.GetSumHardwarePixel(),
		width:            other.GetWidth(),
		height:           other.GetHeight(),
		frameTime:        time.Now(),
		logger:           other.GetLogger(),
	}
}

func (f *frame) ToLedStripe() LedStripe {
	ledStripe := NewLedStripe(f.sumHardwarePixel, f.logger)
	buffer := ledStripe.GetBuffer()
	for _, tile := range f.tiles {
		for x := 0; x < tile.GetWidth(); x++ {
			for y := 0; y < tile.GetHeight(); y++ {
				tilePoint := image.Pt(x, y)
				stripePos := tile.MapTilePixelToStripePosition(tilePoint)
				bufferPos := stripePos * NumBytePixel
				framePoint := tile.FramePoint(tilePoint)
				frameColor := f.RGBAAt(framePoint.X, framePoint.Y)
				buffer[bufferPos+R] = frameColor.R
				buffer[bufferPos+G] = frameColor.G
				buffer[bufferPos+B] = frameColor.B
			}
		}
	}
	return ledStripe
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

func (f *frame) GetTime() time.Time {
	return f.frameTime
}

func (f *frame) GetLogger() *zap.Logger {
	return f.logger
}

func (f *frame) getTiles() []Tile {
	return f.tiles
}
