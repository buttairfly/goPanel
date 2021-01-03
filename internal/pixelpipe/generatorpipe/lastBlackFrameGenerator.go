package generatorpipe

import (
	"context"
	"image/color"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

type lastBlackFrameGenerator struct {
	pipe   *pipepart.Pipe
	logger *zap.Logger
}

// NewLastBlackFramePipe generates a black frame after the input channel was closed
func NewLastBlackFramePipe(
	id pipepart.ID,
	logger *zap.Logger,
) pipepart.PixelPiper {
	pipepart.CheckNoPlaceholderID(id, logger)
	outputChan := make(chan hardware.Frame)
	return &lastBlackFrameGenerator{
		pipe:   pipepart.NewPipe(id, outputChan),
		logger: logger,
	}
}

func (me *lastBlackFrameGenerator) RunPipe(cancelCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(me.pipe.GetFullOutput())

	// first get a frame; if not, just return
	select {
	case <-cancelCtx.Done():
		me.logger.Warn("got cancelCtx.Done() before emptyframe")
		return
	case emptyFrame, ok := <-me.pipe.GetInput():
		if !ok {
			me.logger.Warn("got closed leakybuffer.GetFrameSource() before emptyframe")
			return
		}

		emptyFrame.Fill(color.Black)

		// just pass trough all colorFrames
		for colorFrame := range me.pipe.GetInput() {
			me.pipe.GetFullOutput() <- colorFrame
		}

		// input chan is closed
		me.pipe.GetFullOutput() <- emptyFrame
	}
}

func (me *lastBlackFrameGenerator) GetID() pipepart.ID {
	return me.pipe.GetID()
}

func (me *lastBlackFrameGenerator) GetPrevID() pipepart.ID {
	return me.pipe.GetPrevID()
}

func (me *lastBlackFrameGenerator) Marshal() pipepart.Marshal {
	return pipepart.Marshal{
		ID:     me.GetID(),
		PrevID: me.GetPrevID(),
		Params: me.GetParams(),
	}
}

// GetParams implements PixelPiper interface
func (me *lastBlackFrameGenerator) GetParams() []pipepart.PipeParam {
	return nil
}

func (me *lastBlackFrameGenerator) GetOutput(id pipepart.ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.pipe.GetOutput(id)
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(pipepart.OutputIDMismatchError(me.GetID(), id)))
	return nil
}

func (me *lastBlackFrameGenerator) SetInput(prevID pipepart.ID, inputChan hardware.FrameSource) {
	if pipepart.IsEmptyID(prevID) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(pipepart.PipeIDEmptyError()))
	}
	me.pipe.SetInput(prevID, inputChan)
}
