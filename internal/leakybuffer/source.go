package leakybuffer

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

var source *Source
var freeList = make(chan hardware.Frame, 100)

// Source is a frame recycler and source of frames
type Source struct {
	tileConfigs hardware.TileConfigs
	outputChan  chan hardware.Frame
	logger      *zap.Logger
}

// NewFrameSource produces a frame channel with a recycled frame or new one
func NewFrameSource(tileConfigs hardware.TileConfigs, logger *zap.Logger) *Source {
	source = &Source{
		tileConfigs: tileConfigs,
		outputChan:  make(chan hardware.Frame),
		logger:      logger,
	}
	return source
}

// RunPipe starts the Source
func (me *Source) RunPipe(destroyCtx context.Context, wg *sync.WaitGroup) {
	// wg is only here to implement the PixelPiper interface
	// wg must not get incremented
	defer close(me.outputChan)
	for {
		var f hardware.Frame
		// Grab a buffer if available; allocate if not.
		select {
		case f = <-freeList:
			// frame is still filled with old contents
		case <-destroyCtx.Done():
			return
		default:
			// None free, so allocate a new one.
			f = hardware.NewFrame(me.tileConfigs, me.logger)
		}
		select {
		case me.outputChan <- f:
			// Send to output => will wait for ever
			continue
		case <-destroyCtx.Done():
			return
		}
	}
}

// GetOutput returns the frame producer chan
func (me *Source) GetOutput(id pipepart.ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.outputChan
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(pipepart.OutputIDMismatchError(me.GetID(), id)))
	return nil
}

// GetFrameSource returns the frame producer chan from the global var source
func GetFrameSource(id pipepart.ID) hardware.FrameSource {
	return source.GetOutput(id)
}

// GetID returns the frame producer chan
func (me *Source) GetID() pipepart.ID {
	return pipepart.SourceID
}

// GetPrevID returns an EmptyID by definition
func (me *Source) GetPrevID() pipepart.ID {
	return pipepart.EmptyID
}

// Marshal returns the Marshalled description of Source
func (me *Source) Marshal() pipepart.Marshal {
	return pipepart.Marshal{
		ID:     me.GetID(),
		PrevID: me.GetPrevID(),
		Params: me.GetParams(),
	}
}

// GetParams implements PixelPiper interface
func (me *Source) GetParams() []pipepart.PipeParam {
	return nil
}

// SetInput implements PixelPiper interface and is otherwise useless for Source
func (me *Source) SetInput(prevID pipepart.ID, inputChan hardware.FrameSource) {
	if pipepart.IsEmptyID(prevID) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(pipepart.PipeIDEmptyError()))
	}
	// do nothing
}
