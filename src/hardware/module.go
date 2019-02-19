package hardware

import "image"

// Module hardware struct
type Module struct {
	connectionOrder int
	bounds          image.Rectangle
	ledStripeMap    map[image.Point]int
}

// Bounds implmenents image.Bounds() interface
func (m *Module) Bounds() image.Rectangle {
	return m.bounds
}
