package fader

import "github.com/buttairfly/goPanel/pkg/palette"

// Fader is a color fader with a state
type Fader struct {
	id          ID
	palette     palette.Palette
	currentPos  float64
	granularity int
	wrapping    bool
}

// NewEmptyFader creates an empty Fader
func NewEmptyFader(id ID, start float64, granularity int, wrapping bool) *Fader {
	return &Fader{
		id:          id,
		palette:     palette.NewPalette("palette_" + palette.ID(id)),
		currentPos:  start,
		granularity: granularity,
		wrapping:    wrapping,
	}
}
