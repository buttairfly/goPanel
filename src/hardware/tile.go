package hardware

import (
	"image"
	"strconv"
)

// TileConfig is the config of a tile or led module
type TileConfig struct {
	ConnectionOrder int             `json:"connectionOrder"`
	Bounds          image.Rectangle `json:"bounds"`
	LedStripeMap    map[string]int  `json:"ledStripeMap"`
}

// Tile hardware interface
type Tile interface {
	MapPixelToStripe(framePixel image.Point) int
	MapPositionToStripe(framePos int) int
	Bounds() image.Rectangle
	NumPixel() int
}

type tile struct {
	frameBounds     image.Rectangle
	connectionOrder int
	bounds          image.Rectangle
	ledStripeMap    map[string]int // maps framePosition to physical led
}

// NewTile creates a new Tile
func NewTile(tileConfig TileConfig, frameBounds image.Rectangle) Tile {
	return &tile{
		frameBounds:     frameBounds,
		connectionOrder: tileConfig.ConnectionOrder,
		bounds:          tileConfig.Bounds,
		ledStripeMap:    tileConfig.LedStripeMap,
	}
}

// Bounds implmenents image.Bounds() interface
func (t *tile) Bounds() image.Rectangle {
	return t.bounds
}

func (t *tile) MapPixelToStripe(framePixel image.Point) int {
	framePos := t.frameBounds.Dx()*framePixel.Y + framePixel.X
	return t.MapPositionToStripe(framePos)
}

func (t *tile) MapPositionToStripe(framePos int) int {
	return t.ledStripeMap[strconv.Itoa(framePos)]
}

func (t *tile) NumPixel() int {
	return t.bounds.Dx() * t.bounds.Dy()
}
