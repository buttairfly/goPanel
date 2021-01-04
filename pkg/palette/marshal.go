package palette

// Marshal is the marshalable version of palette.palette and therefore palette.Palette
type Marshal struct {
	ID        ID
	FixColors []FixColor
}

// FixColor is the marshalable version of palette.fixColor
type FixColor struct {
	Color string  `json:"color" yaml:"color"`
	Pos   float64 `json:"pos" yaml:"pos"`
}

// Marshal converts a marshalable palette to palette.Marshal
func (p *palette) Marshal() Marshal {
	fixColors := make([]FixColor, p.Len())
	for i := 0; i < p.Len(); i++ {
		fixColors[i].Color = p.slice()[i].color.Hex()
		fixColors[i].Pos = p.slice()[i].pos
	}
	return Marshal{
		FixColors: fixColors,
		ID:        p.GetID(),
	}
}

// ColorMoveMarshal is a struct to move a color fixpoint From position To position
type ColorMoveMarshal struct {
	From float64 `json:"from" yaml:"from"`
	To   float64 `json:"to" yaml:"to"`
}
