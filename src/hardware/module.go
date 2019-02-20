package hardware

import "image"

type module struct {
	connectionOrder int
	bounds          image.Rectangle
	ledStripeMap    map[string]int
}

// Module hardware interface
type Module interface {
	MapToStripe(p image.Point) int
	Bounds() image.Rectangle
}

// Bounds implmenents image.Bounds() interface
func (m *module) Bounds() image.Rectangle {
	return m.bounds
}

func (m *module) MapToStripe(p image.Point) int {
	return m.ledStripeMap[p.String()]
}
