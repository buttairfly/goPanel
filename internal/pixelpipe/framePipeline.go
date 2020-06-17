package pixelpipe

import (
	"context"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
)

// FramePipeline is a struct which defines a source
type FramePipeline struct {
	destroyCtx        context.Context
	count             int32
	rebuild           chan bool
	rebuildDone       chan bool
	emptyPipeline     chan bool
	internalInputChan chan hardware.Frame

	pipe       *Pipe
	logger     *zap.Logger
	pixelPipes map[string]PixelPiper
}

// NewEmptyFramePipeline creates a new, empty FramePipeline
func NewEmptyFramePipeline(destroyCtx context.Context, id ID, logger *zap.Logger) *FramePipeline {
	pixelPipes := make(map[string]PixelPiper)
	rebuild := make(chan bool)
	rebuildDone := make(chan bool)
	emptyPipeline := make(chan bool)
	outputChan := make(chan hardware.Frame)
	internalInputChan := make(chan hardware.Frame)

	return &FramePipeline{
		destroyCtx:        destroyCtx,
		rebuild:           rebuild,
		emptyPipeline:     emptyPipeline,
		rebuildDone:       rebuildDone,
		internalInputChan: internalInputChan,
		pipe:              NewPipe(id, outputChan),
		logger:            logger,
		pixelPipes:        pixelPipes,
	}
}

// RunPipe implements PixelPiper interface
func (me *FramePipeline) RunPipe(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(me.pipe.outputChan)
	subWg := new(sync.WaitGroup)

	for {
		if inputChannelClosed := me.runInternalPipe(); inputChannelClosed {
			<-me.emptyPipeline
			return
		}
		<-me.emptyPipeline
		me.rebuild = make(chan bool)

		// rebuild pipeline
	}

	subWg.Wait()
}

func (me *FramePipeline) runInternalPipe() bool {
	defer close(me.rebuild)
	for {
		select {
		case <-me.rebuild:
			return false
		default:
		}
		if inputChannelClosed := me.processFrame(); inputChannelClosed {
			return true
		}
	}
}

func (me *FramePipeline) processFrame() bool {
	defer me.decFrameCount()
	me.incFrameCount()
	inputFrame, ok := <-me.pipe.inputChan
	if !ok {
		close(me.internalInputChan)
		return true
	}
	me.internalInputChan <- inputFrame
	return false
}

// GetOutput implements PixelPiper interface
func (me *FramePipeline) GetOutput(id ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.pipe.GetOutput()
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(OutputIDMismatchError("simplePipeIntersection", me.GetID(), id)))
	return nil
}

// SetInput implements PixelPiper interface
func (me *FramePipeline) SetInput(inputID ID, inputChan hardware.FrameSource) {
	me.pipe.SetInput(inputChan)
}

// GetID implements PixelPiper interface
func (me *FramePipeline) GetID() ID {
	return me.pipe.GetID()
}

func (me *FramePipeline) incFrameCount() {
	atomic.AddInt32(&me.count, 1)
}

func (me *FramePipeline) decFrameCount() {
	atomic.AddInt32(&me.count, -1)
	switch {
	case <-me.rebuild:
		if atomic.CompareAndSwapInt32(&me.count, 0, 0) {
			me.emptyPipeline <- true
		}
	default:
	}
}
