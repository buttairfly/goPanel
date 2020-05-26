package hardware

import (
	"fmt"
	"image"
)

// Tile hardware interface
type Tile interface {
	MapTilePixelToStripePosition(tilePixel image.Point) int
	MapTilePositionToStipePosition(tileXYPos int) int
	FramePixelToTilePixel(framePixel image.Point) (image.Point, error)
	MapFramePixelToStripePosition(framePixel image.Point) (int, error)
	Bounds() image.Rectangle
	NumHardwarePixel() int
	FramePoint(tilePoint image.Point) image.Point
	GetWidth() int
	GetHeight() int
}

type tile struct {
	numPreviousLedsOnStripe int
	connectionOrder         int
	numHardwarePixel        int
	width                   int
	height                  int
	bounds                  image.Rectangle
	// ledStripeMap maps tile image pixel position to tile relative stripe position
	// (counting starts at 0 for every tile)
	ledStripeMap map[string]int
}

// NewTile creates a new Tile
func NewTile(tileConfig *TileConfig, numPreviousLedsOnStripe int) Tile {
	return &tile{
		numPreviousLedsOnStripe: numPreviousLedsOnStripe,
		connectionOrder:         tileConfig.GetConnectionOrder(),
		numHardwarePixel:        tileConfig.NumHardwarePixel(),
		width:                   tileConfig.GetBounds().Dx(),
		height:                  tileConfig.GetBounds().Dy(),
		bounds:                  tileConfig.GetBounds(),
		ledStripeMap:            tileConfig.GetLedStripeMap(),
	}
}

// Bounds implmenents image.Bounds() interface
func (t *tile) Bounds() image.Rectangle {
	return t.bounds
}

func (t *tile) GetWidth() int {
	return t.width
}

func (t *tile) GetHeight() int {
	return t.height
}

func (t *tile) MapTilePixelToStripePosition(tilePixel image.Point) int {
	tilePos := t.width*tilePixel.Y + tilePixel.X
	return t.MapTilePositionToStipePosition(tilePos)
}

func (t *tile) MapTilePositionToStipePosition(tileXYPos int) int {
	stripePosition, ok := t.ledStripeMap[tilePositionToString(tileXYPos)]
	if ok && stripePosition >= 0 {
		return stripePosition + t.numPreviousLedsOnStripe
	}
	return -1
}

func (t *tile) NumHardwarePixel() int {
	return t.numHardwarePixel
}

func (t *tile) FramePoint(tilePoint image.Point) image.Point {
	return t.Bounds().Min.Add(tilePoint)
}

func (t *tile) isFramePointInTile(framePixel image.Point) bool {
	return framePixel.In(t.bounds)
}

func (t *tile) FramePixelToTilePixel(framePixel image.Point) (image.Point, error) {
	if !t.isFramePointInTile(framePixel) {
		return image.Pt(-1, -1), fmt.Errorf("FramePoint not in tile")
	}
	return framePixel.Sub(t.bounds.Min), nil
}

func (t *tile) MapFramePixelToStripePosition(framePixel image.Point) (int, error) {
	tilePixel, err := t.FramePixelToTilePixel(framePixel)
	if err != nil {
		return -1, err
	}
	return t.MapTilePixelToStripePosition(tilePixel), nil
}
