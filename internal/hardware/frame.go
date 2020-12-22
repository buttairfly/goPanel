package hardware

import (
	"image"
	"image/color"
	"image/draw"

	"go.uber.org/zap"
)

// FrameFillType is an enum which holds the frame fill type
type FrameFillType int

const (
	// FillTypeFullFrame is used when a full frame needs transferring
	FillTypeFullFrame FrameFillType = iota
	// FillTypeSinglePixelChange is used when only a single pixel changed
	FillTypeSinglePixelChange
	// FillTypeSingleFillColor is used when the frame consists of one color only
	FillTypeSingleFillColor
)

// Frame is a hardware frame
type Frame interface {
	draw.Image
	CopyFromOther(other Frame)
	CopyImageFromOther(other Frame)
	ToLedStripe() LedStripe
	GetSumHardwarePixel() int
	SetRGBA(x, y int, c color.RGBA)
	FillRGBA(c color.RGBA)
	Fill(c color.Color)
	RGBAAt(x, y int) color.RGBA
	GetFillType() FrameFillType
	GetWidth() int
	GetHeight() int
	GetLogger() *zap.Logger

	getTiles() []Tile
}

type frame struct {
	image            *image.RGBA
	tiles            []Tile
	sumHardwarePixel int
	width, height    int
	numPixelChanges  int
	changedPixel     *image.Point
	fillType         FrameFillType
	logger           *zap.Logger
}

// NewFrame returns a paintable image consisting of many led-panel tiles
func NewFrame(tileConfigs TileConfigs, logger *zap.Logger) Frame {
	frameBounds := image.ZR
	tiles := make([]Tile, len(tileConfigs))
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
		fillType:         FillTypeSingleFillColor,
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
		fillType:         FillTypeFullFrame,
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
		fillType:         FillTypeFullFrame,
		logger:           other.GetLogger(),
	}
}

func (f *frame) CopyImageFromOther(other Frame) {
	for x := 0; x < other.GetWidth(); x++ {
		for y := 0; y < other.GetHeight(); y++ {
			f.image.SetRGBA(x, y, other.RGBAAt(x, y))
		}
	}
	f.fillType = other.GetFillType()
}

func (f *frame) CopyFromOther(other Frame) {
	f.CopyImageFromOther(other)
	f.tiles = other.getTiles()
	f.sumHardwarePixel = other.GetSumHardwarePixel()
	f.width = other.GetWidth()
	f.height = other.GetHeight()
	f.logger = other.GetLogger()
}

func (f *frame) ToLedStripe() LedStripe {
	ledStripe := NewLedStripe(
		f.sumHardwarePixel,
		f.numPixelChanges,
		f.mapChangedPixelToStripePosition(),
		f.fillType,
		f.logger,
	)
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

func (f *frame) RGBAAt(x, y int) color.RGBA {
	return f.image.RGBAAt(x, y)
}

func (f *frame) updateFillTypeSetPixel(x, y int) {
	if f.fillType != FillTypeFullFrame {
		oldChangedPixel := f.changedPixel
		f.changedPixel = &image.Point{x, y}
		if oldChangedPixel == nil || !oldChangedPixel.Eq(*f.changedPixel) {
			f.numPixelChanges++
			if f.numPixelChanges == 1 {
				f.fillType = FillTypeSinglePixelChange
			} else {
				f.fillType = FillTypeFullFrame
			}
		}
	}
}

// Set makes image buffer muteable
func (f *frame) Set(x, y int, c color.Color) {
	f.updateFillTypeSetPixel(x, y)
	f.image.Set(x, y, c)
}

// SetRGBA makes image buffer muteable
func (f *frame) SetRGBA(x, y int, c color.RGBA) {
	f.updateFillTypeSetPixel(x, y)
	f.image.SetRGBA(x, y, c)
}

func (f *frame) updateFillTypeFillColor() {
	f.numPixelChanges = 0
	f.changedPixel = nil
	f.fillType = FillTypeSingleFillColor
}

// Fill fills the whole image with a color
func (f *frame) Fill(c color.Color) {
	f.updateFillTypeFillColor()
	for y := 0; y < f.height; y++ {
		for x := 0; x < f.width; x++ {
			f.image.Set(x, y, c)
		}
	}
}

// FillRGBA fills the whole image with a RGBA color
func (f *frame) FillRGBA(c color.RGBA) {
	f.updateFillTypeFillColor()
	for y := 0; y < f.height; y++ {
		for x := 0; x < f.width; x++ {
			f.image.SetRGBA(x, y, c)
		}
	}
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

func (f *frame) GetFillType() FrameFillType {
	return f.fillType
}

func (f *frame) GetLogger() *zap.Logger {
	return f.logger
}

func (f *frame) getTiles() []Tile {
	return f.tiles
}

func (f *frame) mapChangedPixelToStripePosition() []int {
	var stipePositions []int
	if f.fillType == FillTypeSinglePixelChange {
		posArray := make([]int, 1)
		if f.changedPixel == nil {
			f.logger.Warn("fillType must not be FillTypeSinglePixelChange with changedPixel unset")
			f.fillType = FillTypeFullFrame
			return stipePositions
		}
		for _, tile := range f.tiles {
			pos, err := tile.MapFramePixelToStripePosition(*f.changedPixel)
			if err == nil {
				posArray = append(posArray, pos)
			}
		}
		if len(posArray) < 1 {
			f.logger.Sugar().Fatalf(
				"posArray must not be empty with fillType set to FillTypeSinglePixelChange - changedPixel was %+v",
				f.changedPixel,
			)
		}
		stipePositions = posArray
	}
	return stipePositions
}
