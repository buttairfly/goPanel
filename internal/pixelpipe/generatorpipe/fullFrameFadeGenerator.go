package generatorpipe

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
	"github.com/buttairfly/goPanel/pkg/palette"
)

type fullFrameFadeGenerator struct {
	pipe    *pipepart.Pipe
	palette palette.Palette
	logger  *zap.Logger

	numSteps int
}

// FullFrameFadeGenerator generates a new full frame fading frame stream
func FullFrameFadeGenerator(
	id pipepart.ID,
	palette palette.Palette,
	inputChan hardware.FrameSource,
	wg *sync.WaitGroup,
	logger *zap.Logger,
) pipepart.PixelPiper {
	pipepart.CheckNoPlaceholderID(id, logger)
	outputChan := make(chan hardware.Frame)

	return &fullFrameFadeGenerator{
		pipe:     pipepart.NewPipe(id, outputChan),
		logger:   logger,
		palette:  palette,
		numSteps: 100,
	}
}

func (me *fullFrameFadeGenerator) RunPipe(ctx context.Context, wg *sync.WaitGroup) {
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

func (me *fullFrameFadeGenerator) GetID() pipepart.ID {
	return me.pipe.GetID()
}

func (me *fullFrameFadeGenerator) GetPrevID() pipepart.ID {
	return me.pipe.GetPrevID()
}

func (me *fullFrameFadeGenerator) Marshal() *pipepart.Marshal {
	return pipepart.MarshalFromPixelPiperInterface(me)
}

func (me *fullFrameFadeGenerator) GetType() pipepart.PipeType {
	return pipepart.FullFrameFadeGenerator
}

// GetParams implements PixelPiper interface
func (me *fullFrameFadeGenerator) GetParams() []pipepart.PipeParam {
	pp := make([]pipepart.PipeParam, 2)
	pp[0] = pipepart.PipeParam{
		Name:  "palette",
		Type:  pipepart.NameID,
		Value: string(me.palette.GetID()),
	}
	pp[1] = pipepart.PipeParam{
		Name:  "numSteps",
		Type:  pipepart.UInteger,
		Value: fmt.Sprintf("%d", me.numSteps),
	}
	return pp
}

func (me *fullFrameFadeGenerator) GetOutput(id pipepart.ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.pipe.GetOutput(id)
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(pipepart.OutputIDMismatchError(me.GetID(), id)))
	return nil
}

func (me *fullFrameFadeGenerator) SetInput(prevID pipepart.ID, inputChan hardware.FrameSource) {
	if pipepart.IsEmptyID(prevID) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(pipepart.PipeIDEmptyError()))
	}
	me.pipe.SetInput(prevID, inputChan)
}
