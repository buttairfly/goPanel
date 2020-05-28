package hardware

import (
	"fmt"
	"image"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/intmath"
	"github.com/buttairfly/goPanel/pkg/marshal"
)

// MapFormatString is the string to format the map[string]int led position mapping
const MapFormatString = "%2d"

// TileConfig is a struct of a config of one led panel tile
type TileConfig struct {
	ConnectionOrder int
	Bounds          image.Rectangle
	LedStripeMap    map[string]int
}

// NumHardwarePixel counts the number of actual valid hardware pixels in the config
func (tc *TileConfig) NumHardwarePixel() int {
	maxX := tc.Bounds.Dx()
	maxY := tc.Bounds.Dy()
	maxPixel := maxX * maxY
	numHardwarePixel := 0
	maxStripePos := 0
	for tilePos := 0; tilePos < maxPixel; tilePos++ {
		stripePos, ok := tc.LedStripeMap[tilePositionToString(tilePos)]
		stripePoint := image.Point{X: tilePos % maxX, Y: tilePos / maxY}
		maxStripePos = intmath.Max(stripePos, maxStripePos)
		if ok &&
			stripePos >= 0 &&
			stripePoint.X >= 0 &&
			stripePoint.Y >= 0 &&
			stripePoint.X < maxX &&
			stripePoint.Y < maxY {
			numHardwarePixel++
		}
	}
	if maxStripePos > numHardwarePixel-1 {
		zap.L().Sugar().Infof(
			"numHardwarePixel (%d) of tile %d is not within max stripe pos %d",
			numHardwarePixel,
			tc.ConnectionOrder,
			maxStripePos,
		)
	}
	return numHardwarePixel
}

// GetBounds retruns the tile image rectangle
func (tc *TileConfig) GetBounds() image.Rectangle {
	return tc.Bounds
}

// GetConnectionOrder retruns the tile connection order
func (tc *TileConfig) GetConnectionOrder() int {
	return tc.ConnectionOrder
}

// GetLedStripeMap retruns the tile led stripe map
func (tc *TileConfig) GetLedStripeMap() map[string]int {
	return tc.LedStripeMap
}

func tilePointxyToString(x, y, maxX int) string {
	return tilePositionToString(y*maxX + x)
}

func tilePositionToString(pos int) string {
	return fmt.Sprintf(MapFormatString, pos)
}

// ToMarshalTileConfig retruns the marshalable tile config
func (tc *TileConfig) ToMarshalTileConfig() *MarshalTileConfig {
	mtc := new(MarshalTileConfig)
	mtc.ConnectionOrder = tc.ConnectionOrder
	mtc.Bounds = marshal.FromImageRectangle(tc.Bounds)
	mtc.LedStripeMap = tc.LedStripeMap
	return mtc
}

/*
// MarshalJSON implements the json.Marshaller interface
func (tc *TileConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(tc.ToMarshalTileConfig())
}

// UnmarshalJSON implements the json.Unmarshaller interface
func (tc *TileConfig) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, tc.ToMarshalTileConfig())
}

// MarshalYAML implements the yaml.Marshaller interface
func (tc *TileConfig) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(tc.ToMarshalTileConfig())
}

// UnmarshalYAML implements the yaml.UnUnmarshaller interface
func (tc *TileConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(tc.ToMarshalTileConfig())
}
*/
