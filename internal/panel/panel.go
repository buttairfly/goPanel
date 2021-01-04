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
	"github.com/buttairfly/goPanel/internal/pixelpipe/alphablender"
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
	faders        map[fader.ID]fader.Fader
	palettes      map[palette.ID]palette.Palette
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
		faders:        make(map[fader.ID]fader.Fader, 0),
		palettes:      make(map[palette.ID]palette.Palette, 0),
		leakySource:   frameSource,
		frameSource:   frameSource.GetOutput(frameSource.GetID()),
		framePipeline: pixelpipe.NewEmptyFramePipeline(emptyFramePipeID, logger),
	}

	// TODO: move to file
	fire := palette.NewPalette("fire")
	fire.PutAt(colorful.Color{R: 0.1, G: 0, B: 0}, 0)
	fire.PutAt(colorful.Color{R: 0.5, G: 0.1, B: 0}, 1.0/3)
	fire.PutAt(colorful.Color{R: 0.3, G: 0, B: 0}, 2.0/3)
	fire.PutAt(colorful.Color{R: 0.4, G: 0.1, B: 0}, 1.0)
	panel.palettes[fire.GetID()] = fire

	// TODO: move to file
	const c1 = float64(1.0)
	const c0 = float64(1.0 / 255.0)
	rainbowPalette := palette.NewPalette("rainbow")
	rainbowPalette.PutAt(colorful.Color{R: c1, G: c0, B: c0}, 0)
	rainbowPalette.PutAt(colorful.Color{R: c0, G: c1, B: c0}, 1.0/3)
	rainbowPalette.PutAt(colorful.Color{R: c0, G: c0, B: c1}, 2.0/3)
	rainbowPalette.PutAt(colorful.Color{R: c1, G: c0, B: c0}, 1.0)
	panel.palettes[rainbowPalette.GetID()] = rainbowPalette

	panel.framePipeline.SetInput(pipepart.SourceID, panel.frameSource)

	device.SetInput(panel.framePipeline.GetID(), panel.framePipeline.GetOutput(emptyFramePipeID))

	// TODO: load from file
	panel.framePipeline.AddPipeBefore(
		emptyFramePipeID,
		generatorpipe.RainbowGenerator(
			"rainbow",
			rainbowPalette,
			0.009,
			0.02,
			logger,
		),
	)

	panel.framePipeline.AddPipeBefore(
		emptyFramePipeID,
		alphablender.NewClockBlender(
			"24h_clock",
			0.1,
			0.9,
			logger,
		),
	)

	// panel.framePipeline.AddPipeBefore(
	// 	emptyFramePipeID,
	// 	generatorpipe.SnakeGenerator(
	// 		"rainbow",
	// 		rainbowPalette,
	// 		195.0/200,
	// 		logger,
	// 	),
	// )
	return panel
}

// GetPanel returns the global panel
func GetPanel() *Panel {
	return panel
}

// RunPipe starts the panel
func (me *Panel) RunPipe(cancelCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	// do not increment wg for leakySource
	go me.leakySource.RunPipe(cancelCtx, wg)
	wg.Add(1)
	go me.framePipeline.RunPipe(cancelCtx, wg)
}

// Marshal returns the Marshalled Panel and implements PixelPiper interface
func (me *Panel) Marshal() pipepart.Marshal {
	pp := make(map[pipepart.ID]pipepart.Marshal)
	pp[me.leakySource.GetID()] = me.leakySource.Marshal()
	pp[me.framePipeline.GetID()] = me.framePipeline.Marshal()
	pp[me.device.GetID()] = me.device.Marshal()
	return pipepart.Marshal{
		ID:          me.GetID(),
		PrevID:      pipepart.EmptyID,
		FirstPipeID: me.leakySource.GetID(),
		LastPipeID:  me.device.GetID(),
		PixelPipes:  pp,
		Params:      me.GetParams(),
	}
}

// GetID implements PixelPiper interface
func (me *Panel) GetID() pipepart.ID {
	return pipepart.PanelID
}

// GetParams implements PixelPiper interface
func (me *Panel) GetParams() []pipepart.PipeParam {
	return nil
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
