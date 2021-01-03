package fader

import "github.com/buttairfly/goPanel/pkg/palette"

// Fader is a color fader with a state
type Fader struct {
	name        string
	palette     palette.Palette
	currentPos  float64
	granularity int
	wrapping    bool
}

func NewEmptyFader(name string, start float64, granularity int, wrapping bool) *Fader {
	return &Fader{
		name:        name,
		palette:     palette.NewPalette(name),
		currentPos:  start,
		granularity: granularity,
		wrapping:    wrapping,
	}
}
