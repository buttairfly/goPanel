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
	"github.com/buttairfly/goPanel/internal/pixelpipe/generatorpipe"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
	"github.com/buttairfly/goPanel/pkg/fader"
	"github.com/buttairfly/goPanel/pkg/palette"
	"github.com/lucasb-eyer/go-colorful"
)

var panel *Panel

// Panel is a datatype combining everything a LED panels driver needs
type Panel struct {
	config        *config.MainConfig
	device        device.LedDevice
	faders        map[string]fader.Fader
	palettes      map[string]palette.Palette
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
		faders:        make(map[string]fader.Fader, 0),
		palettes:      make(map[string]palette.Palette, 0),
		leakySource:   frameSource,
		frameSource:   frameSource.GetFrameSource(),
		framePipeline: pixelpipe.NewEmptyFramePipeline(emptyFramePipeID, logger),
	}

	// TODO: move to file
	fire := palette.NewPalette()
	fire.AddAt(colorful.Color{R: 0.1, G: 0, B: 0}, 0)
	fire.AddAt(colorful.Color{R: 0.5, G: 0.1, B: 0}, 1.0/3)
	fire.AddAt(colorful.Color{R: 0.3, G: 0, B: 0}, 2.0/3)
	fire.AddAt(colorful.Color{R: 0.4, G: 0.1, B: 0}, 1.0)
	panel.palettes["fire"] = fire

	// TODO: move to file
	rainbowPalette := palette.NewPalette()
	rainbowPalette.AddAt(colorful.Color{R: 1, G: 0, B: 0}, 0)
	rainbowPalette.AddAt(colorful.Color{R: 0, G: 1, B: 0}, 1.0/3)
	rainbowPalette.AddAt(colorful.Color{R: 0, G: 0, B: 1}, 2.0/3)
	rainbowPalette.AddAt(colorful.Color{R: 1, G: 0, B: 0}, 1.0)
	panel.palettes["rainbow"] = rainbowPalette

	panel.framePipeline.SetInput(pipepart.SourceID, panel.frameSource)

	device.SetInput(panel.framePipeline.GetOutput(emptyFramePipeID))

	// TODO: load from file
	panel.framePipeline.AddPipeBefore(
		emptyFramePipeID,
		generatorpipe.RainbowGenerator(
			"rainbow",
			rainbowPalette,
			0.005,
			0.01,
			logger,
		),
	)
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
