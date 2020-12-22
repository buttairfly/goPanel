package panel

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/config"
	"github.com/buttairfly/goPanel/internal/device"
	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/leakybuffer"
	"github.com/buttairfly/goPanel/internal/pixelpipe"
	"github.com/buttairfly/goPanel/internal/pixelpipe/generatorpipe"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
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
	leakySource   *leakybuffer.Source
	frameSource   hardware.FrameSource
	framePipeline *pixelpipe.FramePipeline
}

// NewPanel creates a new panel struct at a global variable
func NewPanel(config *config.MainConfig, device device.LedDevice, logger *zap.Logger) *Panel {
	frameSource := leakybuffer.NewFrameSource(config.TileConfigs.ToTileConfigs(), logger)
	emptyFramePipeID := pipepart.ID("mainPipe")
	panel = &Panel{
		config:        config,
		device:        device,
		faders:        make([]fader.Fader, 0, 1),
		palettes:      make([]palette.Palette, 0, 1),
		leakySource:   frameSource,
		frameSource:   frameSource.GetFrameSource(),
		framePipeline: pixelpipe.NewEmptyFramePipeline(emptyFramePipeID, logger),
	}
	panel.framePipeline.SetInput(pipepart.SourceID, panel.frameSource)
	device.SetInput(panel.framePipeline.GetOutput(emptyFramePipeID))
	panel.framePipeline.AddPipeBefore(emptyFramePipeID, generatorpipe.WhiteNoisePipe(pipepart.ID("whitenoise"), time.Millisecond, logger))
	return panel
}

// GetPanel returns the global panel
func GetPanel() *Panel {
	return panel
}

// Run starts the panel
func (me *Panel) Run(cancelCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	go me.leakySource.Run(cancelCtx)
	wg.Add(1)
	go me.framePipeline.RunPipe(cancelCtx, wg)
}

// GetConfig returns the global config
func (me *Panel) GetConfig() *config.MainConfig {
	return me.config
}

// GetDevice returns the LedDevice
func (me *Panel) GetDevice() device.LedDevice {
	return me.device
}

// GetFramePipeline returns the panel FramePipeline
func (me *Panel) GetFramePipeline() *pixelpipe.FramePipeline {
	return me.framePipeline
}

// GetFaders returns the panel faders
func (me *Panel) GetFaders() []fader.Fader {
	return me.faders
}
