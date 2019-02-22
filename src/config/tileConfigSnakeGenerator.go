package config

import (
	"fmt"
	"image"

	"github.com/buttairfly/goPanel/src/intmath"
)

type direction int

const (
	horizontal direction = iota
	vertical
)

const (
	even = 0
	odd  = 1
)

// TileConfigSnakeGenerator struct to generate a tile config with snake pattern
type TileConfigSnakeGenerator struct {
	connectionOrder int
	startPoint      image.Point
	endPoint        image.Point
	direction       direction
}

// NewTileConfigSnakeMapFile creates a new snake tile config
func NewTileConfigSnakeMapFile(g TileConfigSnakeGenerator) (TileConfig, error) {
	tileBoundsInFrame := image.Rectangle{
		Min: g.startPoint,
		Max: g.endPoint,
	}
	boundsInFrame := tileBoundsInFrame.Canon()
	tileStart := g.startPoint.Sub(boundsInFrame.Min)
	tileEnd := g.endPoint.Sub(boundsInFrame.Min)
	bounds := image.Rectangle{Min: tileStart, Max: tileEnd}.Canon()
	if !((tileStart.X == 0 || tileStart.X == bounds.Dx()) &&
		(tileStart.Y == 0 || tileStart.Y == bounds.Dy())) {
		return nil, fmt.Errorf("Tile start point %s is not a corner of tile bounds %s, %d , %d",
			tileStart, bounds, bounds.Dx(), bounds.Dy())
	}
	if !((tileEnd.X == 0 || tileEnd.X == bounds.Dx()) &&
		(tileEnd.Y == 0 || tileEnd.Y == bounds.Dy())) {
		return nil, fmt.Errorf("TIle end point %s is not a corner of tile bounds %s",
			tileEnd, bounds)
	}

	ledStripeMap := map[string]int{}
	pos := 0
	maxX := bounds.Dx() + 1
	maxY := bounds.Dy() + 1
	snakeParity := odd
	if (g.direction == horizontal && tileStart.Y != 0 && tileStart.Y%2 == 1) ||
		(g.direction == vertical && tileStart.X != 0 && tileStart.X%2 == 1) {
		snakeParity = even
	}
	stride := maxX
	if g.direction == vertical {
		stride = maxY
	}
	for dy := 0; dy < maxY; dy++ {
		for dx := 0; dx < maxX; dx++ {
			x := dx
			y := dy
			if g.direction == vertical {
				x, y = y, x
			}
			x = intmath.Abs(tileStart.X - x)
			y = intmath.Abs(tileStart.Y - y)

			// snake the pixels
			if g.direction == vertical {
				if x%2 == snakeParity {
					y = stride - 1 - y
				}
			} else {
				if y%2 == snakeParity {
					x = stride - 1 - x
				}
			}
			mapKey := tilePointxyToString(x, y, stride)
			prevValue, ok := ledStripeMap[mapKey]
			if ok {
				return nil, fmt.Errorf("Duplicate stripe map x: %d, y: %d, pos: %d, mapKey: %s, prevValue: %d", x, y, pos, mapKey, prevValue)
			}
			ledStripeMap[mapKey] = pos
			pos++
		}
	}
	return &tileConfig{
		ConnectionOrder: g.connectionOrder,
		Bounds:          boundsInFrame,
		LedStripeMap:    ledStripeMap,
	}, nil
}
