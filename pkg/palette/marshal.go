package palette

import "github.com/lucasb-eyer/go-colorful"

// Marshal is the marshalable version of palette.palette and therefore palette.Palette
type Marshal struct {
	ID      ID         `json:"id" yaml:"id"`
	Blender string     `json:"blender" yaml:"blender"`
	Colors  []FixColor `json:"colors" yaml:"colors"`
}

// FixColor is the marshalable version of palette.fixColor
type FixColor struct {
	Pos   float64 `json:"pos" yaml:"pos"`
	Color string  `json:"color" yaml:"color"`
}

// Marshal converts a marshalable palette to palette.Marshal
func (p *palette) Marshal() *Marshal {
	fixColors := make([]FixColor, p.Len())
	for i, val := range p.fixColors {
		fixColors[i].Color = val.color.Hex()
		fixColors[i].Pos = val.pos
	}
	return &Marshal{
		ID:      p.GetID(),
		Blender: p.blender.String(),
		Colors:  fixColors,
	}
}

// Unmarshal unmarshals a Marshal json to Palette
func (m *Marshal) Unmarshal() (Palette, error) {
	fixColors := make([]fixColor, len(m.Colors))
	for i, val := range m.Colors {
		c, err := colorful.Hex(val.Color)
		if err != nil {
			return nil, err
		}
		fixColors[i].color = c
		fixColors[i].pos = val.Pos
	}
	b, errBlender := FromString(m.Blender)
	if errBlender != nil {
		return nil, errBlender
	}
	return &palette{
		id:        m.ID,
		blender:   b,
		fixColors: fixColors,
	}, nil
}

// ColorMoveMarshal is a struct to move a color fixpoint From position To position
type ColorMoveMarshal struct {
	From float64 `json:"from" yaml:"from"`
	To   float64 `json:"to" yaml:"to"`
}

// SetBlenderMarshal is used to set a blender
type SetBlenderMarshal struct {
	Blender string `json:"blender" yaml:"blender"`
}
