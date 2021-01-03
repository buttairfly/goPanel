package generatorpipe

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/leakybuffer"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
	"github.com/buttairfly/goPanel/pkg/palette"
)

//TODO: replace pos with real draw commands
var pos float64

// DrawCommand is a command to draw on a frame
type DrawCommand string

type drawGenerator struct {
	pipe    *pipepart.Pipe
	logger  *zap.Logger
	palette palette.Palette

	commandInput <-chan DrawCommand
}

// DrawGenerator generates for each command a draw step and draws a new frame
func DrawGenerator(
	id pipepart.ID,
	palette palette.Palette,
	logger *zap.Logger,
	commandInput <-chan DrawCommand,
) pipepart.PixelPiper {
	if pipepart.IsPlaceholderID(id) {
		logger.Fatal("PipeIDPlaceholderError", zap.Error(pipepart.PipeIDPlaceholderError(id)))
	}
	outputChan := make(chan hardware.Frame)

	return &drawGenerator{
		pipe:         pipepart.NewPipe(id, outputChan),
		palette:      palette,
		logger:       logger,
		commandInput: commandInput,
	}
}

func (me *drawGenerator) RunPipe(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(me.pipe.GetFullOutput())
	for frame := range me.pipe.GetInput() {
		isClosed := me.interpretCommand(frame)
		if isClosed {
			leakybuffer.DumpFrame(frame)
			return
		}

		// TODO: frame counter logic
		me.pipe.GetFullOutput() <- frame
	}
}

func (me *drawGenerator) GetID() pipepart.ID {
	return me.pipe.GetID()
}

func (me *drawGenerator) GetPrevID() pipepart.ID {
	return me.pipe.GetPrevID()
}

func (me *drawGenerator) Marshal() pipepart.Marshal {
	return pipepart.Marshal{
		ID:     me.GetID(),
		PrevID: me.GetPrevID(),
		Params: me.GetParams(),
	}
}

// GetParams implements PixelPiper interface
func (me *drawGenerator) GetParams() []pipepart.PipeParam {
	pp := make([]pipepart.PipeParam, 1)
	pp[0] = pipepart.PipeParam{
		Name:  "palette",
		Type:  pipepart.NameID,
		Value: me.palette.GetName(),
	}
	return pp
}

func (me *drawGenerator) GetOutput(id pipepart.ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.pipe.GetOutput(id)
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(pipepart.OutputIDMismatchError(me.GetID(), id)))
	return nil
}

func (me *drawGenerator) SetInput(prevID pipepart.ID, inputChan hardware.FrameSource) {
	if pipepart.IsEmptyID(prevID) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(pipepart.PipeIDEmptyError()))
	}
	me.pipe.SetInput(prevID, inputChan)
}

func (me *drawGenerator) interpretCommand(frame hardware.Frame) bool {
	_, ok := <-me.commandInput
	if !ok {
		return true
	}

	//TODO: replace pos with real draw commands, use xy properly
	pos += 0.01
	if pos > 1 {
		pos = 0
	}
	frame.Set(0, 0, me.palette.Blend(pos))
	return false
}
