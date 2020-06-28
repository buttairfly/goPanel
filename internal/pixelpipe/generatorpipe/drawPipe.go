package generatorpipe

import (
	"sync"

	"github.com/lucasb-eyer/go-colorful"
	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/leakybuffer"
	"github.com/buttairfly/goPanel/internal/pixelpipe"
	"github.com/buttairfly/goPanel/pkg/palette"
)

//TODO: replace p and pos with real draw commands
var p palette.Palette
var pos float64

// DrawCommand is a command to draw on a frame
type DrawCommand string

type drawPipe struct {
	pipe   *pixelpipe.Pipe
	logger *zap.Logger

	commandInput <-chan DrawCommand
}

// DrawPipe generates for each command a draw step and draws a new frame
func DrawPipe(
	id pixelpipe.ID,
	logger *zap.Logger,
	commandInput <-chan DrawCommand,
) pixelpipe.PixelPiper {
	outputChan := make(chan hardware.Frame)

	//TODO: replace p and pos with real draw commands
	p = palette.NewPalette()
	p.AddAt(colorful.Color{R: 0xff, G: 0, B: 0}, 0)
	p.AddAt(colorful.Color{R: 0, G: 0, B: 0xff}, 0.5)
	pos = 0.0

	return &drawPipe{
		pipe:         pixelpipe.NewPipe(id, outputChan),
		logger:       logger,
		commandInput: commandInput,
	}
}

func (me *drawPipe) RunPipe(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(me.pipe.GetFullOutput())
	for frame := range me.pipe.GetInput() {
		isClosed := me.interpretCommand(frame)
		if !isClosed {
			leakybuffer.DumpFrame(frame)
			return
		}

		// TODO: frame counter logic
		me.pipe.GetFullOutput() <- frame
	}
}

func (me *drawPipe) GetID() pixelpipe.ID {
	return me.pipe.GetID()
}

func (me *drawPipe) GetPrevID() pixelpipe.ID {
	return me.pipe.GetPrevID()
}

func (me *drawPipe) GetOutput(id pixelpipe.ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.pipe.GetOutput(id)
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(pixelpipe.OutputIDMismatchError("simplePipeIntersection", me.GetID(), id)))
	return nil
}

func (me *drawPipe) SetInput(inputID pixelpipe.ID, inputChan hardware.FrameSource) {
	me.pipe.SetInput(inputID, inputChan)
}

func (me *drawPipe) interpretCommand(frame hardware.Frame) bool {
	_, ok := <-me.commandInput
	if !ok {
		return true
	}

	//TODO: replace p and pos with real draw commands
	pos += 0.01
	if pos > 1 {
		pos = 0
	}
	frame.Set(0, 0, p.Blend(pos))
	return false
}
