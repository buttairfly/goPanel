package hardware

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
	bounds := tileBoundsInFrame.Canon()
	if !((g.startPoint.X == 0 || g.startPoint.X == bounds.Dx()) &&
		(g.startPoint.Y == 0 || g.startPoint.Y == bounds.Dy())) {
		return nil, fmt.Errorf("Starting point %s is not a corner of tile bounds %s, %d , %d",
			g.startPoint, bounds, bounds.Dx(), bounds.Dy())
	}
	if !((g.endPoint.X == 0 || g.endPoint.X == bounds.Dx()) &&
		(g.endPoint.Y == 0 || g.endPoint.Y == bounds.Dy())) {
		return nil, fmt.Errorf("Ending point %s is not a corner of tile bounds %s",
			g.endPoint, bounds)
	}

	ledStripeMap := map[string]int{}
	pos := 0
	maxX := bounds.Dx() + 1
	maxY := bounds.Dy() + 1
	for dy := 0; dy < maxY; dy++ {
		for dx := 0; dx < maxX; dx++ {
			x := intmath.Abs(g.startPoint.X - dx)
			y := intmath.Abs(g.startPoint.Y - dy)

			// snake the pixels
			if g.direction == vertical {
				if x%2 == 1 {
					y = bounds.Dy() - y
				}
			} else {
				if y%2 == 1 {
					x = bounds.Dx() - x
				}
			}

			mapKey := ""
			if g.direction == vertical {
				mapKey = xyToPositionString(y, x, maxY)
			} else {
				mapKey = xyToPositionString(x, y, maxX)
			}
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
		Bounds:          bounds,
		LedStripeMap:    ledStripeMap,
	}, nil
}

func xyToPositionString(x, y, maxX int) string {
	return fmt.Sprintf("%04d", y*maxX+x)
}
