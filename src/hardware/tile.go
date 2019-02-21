package hardware

import (
	"image"
	"strconv"
)

// Tile hardware interface
type Tile interface {
	MapTilePixelToStripePosition(tilePixel image.Point) int
	MapTilePositionToStipePosition(tileXYPos int) int
	Bounds() image.Rectangle
	NumHardwarePixel() int
	FramePoint(tilePoint image.Point) image.Point
}

type tile struct {
	numPreviousLedsOnStripe int
	connectionOrder         int
	numHardwarePixel        int
	bounds                  image.Rectangle
	// ledStripeMap maps tile led position to tile relative stripe position
	// (counting starts at 0 for every tile)
	ledStripeMap map[string]int
}

// NewTile creates a new Tile
func NewTile(tileConfig TileConfig, numPreviousLedsOnStripe int) Tile {
	return &tile{
		numPreviousLedsOnStripe: numPreviousLedsOnStripe,
		connectionOrder:         tileConfig.GetConnectionOrder(),
		numHardwarePixel:        tileConfig.NumHardwarePixel(),
		bounds:                  tileConfig.GetBounds(),
		ledStripeMap:            tileConfig.GetLedStripeMap(),
	}
}

// Bounds implmenents image.Bounds() interface
func (t *tile) Bounds() image.Rectangle {
	return t.bounds
}

func (t *tile) MapTilePixelToStripePosition(tilePixel image.Point) int {
	tilePos := t.bounds.Dx()*tilePixel.Y + tilePixel.X
	return t.MapTilePositionToStipePosition(tilePos)
}

func (t *tile) MapTilePositionToStipePosition(tileXYPos int) int {
	stripePosition, ok := t.ledStripeMap[strconv.Itoa(tileXYPos)]
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
