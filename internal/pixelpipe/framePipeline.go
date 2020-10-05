package pixelpipe

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
)

// FramePipeline is a struct which defines a source
type FramePipeline struct {
	destroyCtx       context.Context
	running          bool
	rebuild          chan bool
	frameWg          *sync.WaitGroup
	internalSource   hardware.FrameSource
	internalLastPipe *Pipe
	logger           *zap.Logger
	pixelPipes       map[ID]PixelPiper
	lastPipeID       ID
	firstPipeID      ID
	prevID           ID
}

// NewEmptyFramePipeline creates a new, empty FramePipeline which can hold multiple pipes end-to-end connected to each other
func NewEmptyFramePipeline(destroyCtx context.Context, id ID, logger *zap.Logger) *FramePipeline {
	if IsPlaceholderID(id) {
		logger.Fatal("PipeIDPlaceholderError", zap.Error(PipeIDPlaceholderError(id)))
	}
	pixelPipes := make(map[ID]PixelPiper)
	rebuild := make(chan bool)
	outputChan := make(chan hardware.Frame)
	pipe := NewPipe(id, outputChan)

	return &FramePipeline{
		destroyCtx:       destroyCtx,
		rebuild:          rebuild,
		frameWg:          new(sync.WaitGroup),
		internalLastPipe: pipe,
		logger:           logger,
		pixelPipes:       pixelPipes,
		firstPipeID:      EmptyID,
		lastPipeID:       EmptyID,
		prevID:           EmptyID,
	}
}

// RunPipe implements PixelPiper interface
func (me *FramePipeline) RunPipe(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(me.internalLastPipe.outputChan)
	me.running = true

	me.startPipePieces(me.frameWg)

	for {
		if rebuildInProgress := me.runInternalPipe(); rebuildInProgress {
			select {
			case <-me.rebuild:
				continue
			default:
				me.frameWg.Wait()
				me.running = false
				me.frameWg = new(sync.WaitGroup)
				return
			}
		} else {
			return
		}

		//TODO wait until pipeline rebuild is ready
	}
}

func (me *FramePipeline) runInternalPipe() bool {
	defer func() {
		me.rebuild = make(chan bool)
	}()
	for {
		if IsEmptyID(me.lastPipeID) || me.internalSource == nil {
			return false
		}
		var sourceChan hardware.FrameSource
		if IsEmptyID(me.lastPipeID) {
			sourceChan = me.pixelPipes[me.lastPipeID].GetOutput(me.lastPipeID)
		} else {
			sourceChan = me.internalSource
		}
		sourceFrame, ok := <-sourceChan
		if !ok {
			return true
		}
		me.internalLastPipe.applyFrame(sourceFrame)
	}
}

func (me *FramePipeline) startPipePieces(wg *sync.WaitGroup) {
	wg.Add(len(me.pixelPipes))
	for _, pipe := range me.pixelPipes {
		go pipe.RunPipe(wg)
	}
}

// GetOutput implements PixelPiper interface
func (me *FramePipeline) GetOutput(id ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.internalLastPipe.GetOutput(id)
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(OutputIDMismatchError(me.GetID(), id)))
	return nil
}

// SetInput implements PixelPiper interface
func (me *FramePipeline) SetInput(prevID ID, inputChan hardware.FrameSource) {
	if IsEmptyID(prevID) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(PipeIDEmptyError()))
	}
	me.internalSource = inputChan
	if me.firstPipeID == EmptyID {
		me.internalLastPipe.SetInput(prevID, inputChan)
	} else {
		me.pixelPipes[me.firstPipeID].SetInput(prevID, inputChan)
	}
	me.prevID = prevID
}

// GetID implements PixelPiper interface
func (me *FramePipeline) GetID() ID {
	return me.internalLastPipe.GetID()
}

// AddPipeBefore adds a pipe segment before id
func (me *FramePipeline) AddPipeBefore(id ID, newPipe PixelPiper) {
	if IsEmptyID(id) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(PipeIDEmptyError()))
	}
	if newPipe.GetID() == me.GetID() {
		me.logger.Fatal("PipeIDNotUniqueError", zap.Error(PipeIDNotUniqueError(me.GetID())))
	}
	if me.running {
		// stop pipeline and wait until all frames are empty
		close(me.rebuild)
		for me.running {
			// TODO: what to do when
		}
	}
	me.addPipeBefore(id, newPipe)

}

func (me *FramePipeline) addPipeBefore(id ID, newPipe PixelPiper) {

	// actually insert pipe or preplace itself

	afterPipe, ok := me.pixelPipes[id]
	if !ok {
		if id != me.GetID() {
			me.logger.Fatal("PipeIDMismatchError afterPipe", zap.Error(PipeIDMismatchError(id, me.GetID())))
		} else {
			afterPipe = me.internalLastPipe
			me.lastPipeID = newPipe.GetID()
		}
	}
	prevID := afterPipe.GetPrevID()

	if !IsEmptyID(prevID) {
		if me.prevID == prevID {
			newPipe.SetInput(prevID, me.internalSource)
		} else {
			var prevPipe PixelPiper
			prevPipe, ok = me.pixelPipes[prevID]
			if !ok {
				me.logger.Fatal("PipeIDMismatchError prevPipe", zap.Error(PipeIDMismatchError(prevID, me.prevID)))
			}
			if prevID == me.firstPipeID {
				me.firstPipeID = newPipe.GetID()
			}
			newPipe.SetInput(prevID, prevPipe.GetOutput(prevID))
		}
	}

	afterPipe.SetInput(newPipe.GetID(), newPipe.GetOutput(newPipe.GetID()))
	me.pixelPipes[newPipe.GetID()] = newPipe
}
