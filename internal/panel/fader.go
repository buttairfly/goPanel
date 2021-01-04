package panel

import (
	"errors"

	"github.com/buttairfly/goPanel/pkg/fader"
)

// GetFaders returns the panel faders
func (me *Panel) GetFaders() map[fader.ID]*fader.Fader {
	return me.faders
}

// GetFaderByID returns the panel fader by id
func (me *Panel) GetFaderByID(id fader.ID) (*fader.Fader, error) {
	fader, exists := me.faders[id]
	if exists {
		return fader, nil
	}
	return nil, errors.New("not found")
}
