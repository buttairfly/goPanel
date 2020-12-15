package generatorpipe

import (
	"context"
	"image/color"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/leakybuffer"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

type lastBlackFramePipe struct {
	cancelCtx context.Context
	pipe      *pipepart.Pipe
	logger    *zap.Logger
}

// NewLastBlackFramePipe generates a black frame after the input channel was closed
func NewLastBlackFramePipe(
	cancelCtx context.Context,
	id pipepart.ID,
	logger *zap.Logger,
) pipepart.PixelPiper {
	if pipepart.IsPlaceholderID(id) {
		logger.Fatal("PipeIDPlaceholderError", zap.Error(pipepart.PipeIDPlaceholderError(id)))
	}

	outputChan := make(chan hardware.Frame)
	return &lastBlackFramePipe{
		cancelCtx: cancelCtx,
		pipe:      pipepart.NewPipe(id, outputChan),
		logger:    logger,
	}
}

func (me *lastBlackFramePipe) RunPipe(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(me.pipe.GetFullOutput())

	// first get a frame; if not, just return
	select {
	case <-me.cancelCtx.Done():
		me.logger.Warn("got cancelCtx.Done() before emptyframe")
		return
	case emptyFrame, ok := <-leakybuffer.GetFrameSource():
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

func (me *lastBlackFramePipe) GetID() pipepart.ID {
	return me.pipe.GetID()
}

func (me *lastBlackFramePipe) GetPrevID() pipepart.ID {
	return me.pipe.GetPrevID()
}

func (me *lastBlackFramePipe) GetOutput(id pipepart.ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.pipe.GetOutput(id)
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(pipepart.OutputIDMismatchError(me.GetID(), id)))
	return nil
}

func (me *lastBlackFramePipe) SetInput(prevID pipepart.ID, inputChan hardware.FrameSource) {
	if pipepart.IsEmptyID(prevID) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(pipepart.PipeIDEmptyError()))
	}
	me.pipe.SetInput(prevID, inputChan)
}
