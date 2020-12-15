package generatorpipe

import (
	"sync"

	"github.com/lucasb-eyer/go-colorful"
	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
	"github.com/buttairfly/goPanel/pkg/palette"
)

type fullFrameFadePipe struct {
	pipe   *pipepart.Pipe
	logger *zap.Logger

	numSteps int
	palette  palette.Palette
}

// FullFrameFadePipe generates a new full frame fading frame stream
func FullFrameFadePipe(
	id pipepart.ID,
	inputChan hardware.FrameSource,
	wg *sync.WaitGroup,
	logger *zap.Logger,
) pipepart.PixelPiper {
	if pipepart.IsPlaceholderID(id) {
		logger.Fatal("PipeIDPlaceholderError", zap.Error(pipepart.PipeIDPlaceholderError(id)))
	}
	outputChan := make(chan hardware.Frame)

	palette := palette.NewPalette()
	palette.AddAt(colorful.Color{R: 0xff, G: 0, B: 0}, 0)
	palette.AddAt(colorful.Color{R: 0xff, G: 0xa5, B: 0}, 0.5)

	return &fullFrameFadePipe{
		pipe:     pipepart.NewPipe(id, outputChan),
		logger:   logger,
		numSteps: 100,
		palette:  palette,
	}
}

func (me *fullFrameFadePipe) RunPipe(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(me.pipe.GetFullOutput())

	step := 1.0 / float64(me.numSteps)
	for {
		for i := 0; i < me.numSteps; i++ {
			colorFrame, ok := <-me.pipe.GetInput()
			if !ok {
				return
			}
			color := me.palette.Blend(step * float64(i))
			colorFrame.Fill(color)
			// TODO: frame counter logic
			me.pipe.GetFullOutput() <- colorFrame
		}
	}
}

func (me *fullFrameFadePipe) GetID() pipepart.ID {
	return me.pipe.GetID()
}

func (me *fullFrameFadePipe) GetPrevID() pipepart.ID {
	return me.pipe.GetPrevID()
}

func (me *fullFrameFadePipe) GetOutput(id pipepart.ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.pipe.GetOutput(id)
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(pipepart.OutputIDMismatchError(me.GetID(), id)))
	return nil
}

func (me *fullFrameFadePipe) SetInput(prevID pipepart.ID, inputChan hardware.FrameSource) {
	if pipepart.IsEmptyID(prevID) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(pipepart.PipeIDEmptyError()))
	}
	me.pipe.SetInput(prevID, inputChan)
}
