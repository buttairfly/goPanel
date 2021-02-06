package panel

import (
	"context"
	"fmt"
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
)

var panel *Panel

// Panel is a datatype combining everything a LED panels driver needs
type Panel struct {
	config        *config.MainConfig
	device        device.LedDevice
	faders        map[fader.ID]*fader.Fader
	palettes      map[palette.ID]palette.Palette
	leakySource   *leakybuffer.Source
	frameSource   hardware.FrameSource
	framePipeline *pixelpipe.FramePipeline
}

// NewPanel creates a new panel struct at a global variable
func NewPanel(config *config.MainConfig, device device.LedDevice, logger *zap.Logger) *Panel {
	palette.Init("config/palette/")
	frameSource := leakybuffer.NewFrameSource(config.TileConfigs.ToTileConfigs(), logger)
	emptyFramePipeID := pipepart.ID("mainPipe")

	panel = &Panel{
		config:        config,
		device:        device,
		faders:        make(map[fader.ID]*fader.Fader, 0),
		palettes:      palette.GetGlobal(),
		leakySource:   frameSource,
		frameSource:   frameSource.GetOutput(frameSource.GetID()),
		framePipeline: pixelpipe.NewEmptyFramePipeline(emptyFramePipeID, logger),
	}
	panel.framePipeline.SetInput(pipepart.SourceID, panel.frameSource)
	device.SetInput(panel.framePipeline.GetID(), panel.framePipeline.GetOutput(emptyFramePipeID))

	defaultPalette, _ := palette.GetGlobalByID(palette.DefaultID)
	// TODO: load from file
	// panel.framePipeline.AddPipeBefore(
	// 	emptyFramePipeID,
	// 	generatorpipe.WhiteNoiseGenerator(
	// 		"noise",
	// 		defaultPalette,
	// 		1,
	// 		logger,
	// 	),
	// )
	panel.framePipeline.AddPipeBefore(
		emptyFramePipeID,
		generatorpipe.RainbowGenerator(
			"rainbow",
			defaultPalette,
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
	// 		defaultPalette,
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
func (me *Panel) Marshal() *pipepart.Marshal {
	m := pipepart.MarshalFromPixelPiperBaseInterface(me)
	m.FirstPipeID = me.leakySource.GetID()
	m.LastPipeID = me.device.GetID()
	m.PixelPipes = me.GetPipes()
	return m
}

// GetPipeByID implements PixelPiperWithSubPipes interface
func (me *Panel) GetPipeByID(id pipepart.ID) (pipepart.PixelPiper, error) {
	if pipepart.IsPlaceholderID(id) {
		return nil, fmt.Errorf("can not return placeholder id '%#v'", id)
	}
	switch id {
	case me.framePipeline.GetID():
		return me.framePipeline, nil
	default:
		return me.framePipeline.GetPipeByID(id)
	}
}

// GetReservedPipes implements PixelPiperWithSubPipes interface
func (me *Panel) GetReservedPipes() []pipepart.PipeType {
	return pipepart.GetReservedPipeTypes()
}

// GetPipes implements PixelPiperWithSubPipes interface
func (me *Panel) GetPipes() pipepart.PipesMarshal {
	pipeParts := make(pipepart.PipesMarshal, 3)
	pipeParts[me.leakySource.GetID()] = me.leakySource.Marshal()
	pipeParts[me.framePipeline.GetID()] = me.framePipeline.Marshal()
	pipeParts[me.device.GetID()] = me.device.Marshal()
	return pipeParts
}

// GetType implements PixelPiper interface
func (me *Panel) GetType() pipepart.PipeType {
	return pipepart.Panel
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
