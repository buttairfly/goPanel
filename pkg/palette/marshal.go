package palette

import (
	"github.com/lucasb-eyer/go-colorful"
)

// Marshal is the marshalable version of palette.palette and therefore palette.Palette
type Marshal []ColorMarshal

// ColorMarshal is the marshalable version of palette.paletteColor
type ColorMarshal struct {
	Color colorful.Color `json:"color" yaml:"color"`
	Pos   float64        `json:"pos" yaml:"pos"`
}

// ColorMoveMarshal is a struct to move a color fixpoint From position To position
type ColorMoveMarshal struct {
	From float64 `json:"from" yaml:"from"`
	To   float64 `json:"to" yaml:"to"`
}

// ToMarshal converts a marshalable palette to palette.Marshal
func (p *palette) ToMarshal() Marshal {
	m := make(Marshal, p.Len())
	for i := 0; i < p.Len(); i++ {
		m[i].Color = p.slice()[i].color
		m[i].Pos = p.slice()[i].pos
	}
	return m
}
