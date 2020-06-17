package panel

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/config"
	"github.com/buttairfly/goPanel/internal/device"
	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/leakybuffer"
	"github.com/buttairfly/goPanel/internal/pixelpipe"
	"github.com/buttairfly/goPanel/pkg/fader"
	"github.com/buttairfly/goPanel/pkg/palette"
)

var panel *Panel

// Panel is a datatype combining everything a LED panels driver needs
type Panel struct {
	config        *config.MainConfig
	device        device.LedDevice
	faders        []fader.Fader
	palettes      []palette.Palette
	frameSource   hardware.FrameSource
	framePipeline *pixelpipe.FramePipeline
}

// NewPanel creates a new panel struct at a global variable
func NewPanel(cancelCtx context.Context, config *config.MainConfig, device device.LedDevice, wg *sync.WaitGroup, logger *zap.Logger) {
	wg.Add(1)
	go leakybuffer.NewFrameSource(cancelCtx, config.TileConfigs.ToTileConfigs(), wg, logger)
	panel = &Panel{
		config:        config,
		device:        device,
		faders:        make([]fader.Fader, 0, 1),
		palettes:      make([]palette.Palette, 0, 1),
		frameSource:   leakybuffer.GetFrameSource(),
		framePipeline: pixelpipe.NewEmptyFramePipeline(cancelCtx, "pipe", logger),
	}
}

// GetPanel returns the global panel
func GetPanel() *Panel {
	return panel
}

func (me *Panel) GetConfig() *config.MainConfig {
	return me.config
}

func (me *Panel) GetDevice() device.LedDevice {
	return me.device
}

func (me *Panel) GetFramePipeline() *pixelpipe.FramePipeline {
	return me.framePipeline
}

func (me *Panel) GetFaders() []fader.Fader {
	return me.faders
}
