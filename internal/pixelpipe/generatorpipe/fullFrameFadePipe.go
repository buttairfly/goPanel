package generatorpipe

import (
	"sync"

	"github.com/lucasb-eyer/go-colorful"
	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/pixelpipe"
	"github.com/buttairfly/goPanel/pkg/palette"
)

type fullFrameFadePipe struct {
	pipe   *pixelpipe.Pipe
	logger *zap.Logger

	numSteps int
	palette  palette.Palette
}

// FullFrameFadePipe generates a new full frame fading frame stream
func FullFrameFadePipe(
	id pixelpipe.ID,
	inputChan hardware.FrameSource,
	wg *sync.WaitGroup,
	logger *zap.Logger,
) pixelpipe.PixelPiper {
	if pixelpipe.IsPlaceholderID(id) {
		logger.Fatal("PipeIDPlaceholderError", zap.Error(pixelpipe.PipeIDPlaceholderError(id)))
	}
	outputChan := make(chan hardware.Frame)

	palette := palette.NewPalette()
	palette.AddAt(colorful.Color{R: 0xff, G: 0, B: 0}, 0)
	palette.AddAt(colorful.Color{R: 0xff, G: 0xa5, B: 0}, 0.5)

	return &fullFrameFadePipe{
		pipe:     pixelpipe.NewPipe(id, outputChan),
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

func (me *fullFrameFadePipe) GetID() pixelpipe.ID {
	return me.pipe.GetID()
}

func (me *fullFrameFadePipe) GetPrevID() pixelpipe.ID {
	return me.pipe.GetPrevID()
}

func (me *fullFrameFadePipe) GetOutput(id pixelpipe.ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.pipe.GetOutput(id)
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(pixelpipe.OutputIDMismatchError(me.GetID(), id)))
	return nil
}

func (me *fullFrameFadePipe) SetInput(prevID pixelpipe.ID, inputChan hardware.FrameSource) {
	if pixelpipe.IsEmptyID(prevID) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(pixelpipe.PipeIDEmptyError()))
	}
	me.pipe.SetInput(prevID, inputChan)
}
