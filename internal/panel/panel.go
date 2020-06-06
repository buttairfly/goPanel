package panel

import (
	"github.com/buttairfly/goPanel/internal/config"
	"github.com/buttairfly/goPanel/internal/device"
	"github.com/buttairfly/goPanel/pkg/fader"
	"github.com/buttairfly/goPanel/pkg/palette"
)

var panel *Panel

// Panel is a datatype combining everything a LED panels driver needs
type Panel struct {
	config   *config.MainConfig
	device   device.LedDevice
	faders   []fader.Fader
	palettes []palette.Palette
}

// NewPanel creates a new panel struct at a global variable
func NewPanel(config *config.MainConfig, device device.LedDevice) {
	panel = &Panel{
		config:   config,
		device:   device,
		faders:   make([]fader.Fader, 0, 1),
		palettes: make([]palette.Palette, 0, 1),
	}
}

// GetPanel returns the global panel
func GetPanel() *Panel {
	return panel
}
