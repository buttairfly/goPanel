package panel

import (
	"errors"

	"github.com/buttairfly/goPanel/pkg/palette"
)

// GetPalettes returns the panel palettes
func (me *Panel) GetPalettes() map[string]palette.Marshal {
	p := make(map[string]palette.Marshal, len(me.palettes))
	for id, palette := range me.palettes {
		p[id] = palette.ToMarshal()
	}
	return p
}

// GetPaletteByID returns the panel palette by id
func (me *Panel) GetPaletteByID(id string) (palette.Marshal, error) {
	palette, exists := me.palettes[id]
	if exists {
		return palette.ToMarshal(), nil
	}
	return nil, errors.New("not found")
}
