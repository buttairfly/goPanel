package panel

import (
	"errors"

	"github.com/buttairfly/goPanel/pkg/palette"
)

// GetMarshalledPalettes returns the marshalled panel palettes
func (me *Panel) GetMarshalledPalettes() map[palette.ID]*palette.Marshal {
	p := make(map[palette.ID]*palette.Marshal, len(me.palettes))
	for id, currentPalette := range me.palettes {
		p[id] = currentPalette.Marshal()
	}
	return p
}

// GetMarshaledPaletteByID returns the marshalled panel palette by id
func (me *Panel) GetMarshaledPaletteByID(id palette.ID) (*palette.Marshal, error) {
	p, err := me.GetPaletteByID(id)
	if err != nil {
		return nil, err
	}
	return p.Marshal(), nil
}

// GetPaletteByID returns the panel palette by id
func (me *Panel) GetPaletteByID(id palette.ID) (palette.Palette, error) {
	palette, exists := me.palettes[id]
	if exists {
		return palette, nil
	}
	return nil, errors.New("not found")
}
