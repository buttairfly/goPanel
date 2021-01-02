package panel

import (
	"errors"

	"github.com/buttairfly/goPanel/pkg/palette"
)

// GetMarshalledPalettes returns the marshalled panel palettes
func (me *Panel) GetMarshalledPalettes() map[string]palette.Marshal {
	p := make(map[string]palette.Marshal, len(me.palettes))
	for id, palette := range me.palettes {
		p[id] = palette.Marshal()
	}
	return p
}

// GetMarshaledPaletteByID returns the marshalled panel palette by id
func (me *Panel) GetMarshaledPaletteByID(id string) (palette.Marshal, error) {
	palette, err := me.GetPaletteByID(id)
	if err != nil {

		return nil, err
	}
	return palette.Marshal(), nil
}

// GetPaletteByID returns the panel palette by id
func (me *Panel) GetPaletteByID(id string) (palette.Palette, error) {
	palette, exists := me.palettes[id]
	if exists {
		return palette, nil
	}
	return nil, errors.New("not found")
}
