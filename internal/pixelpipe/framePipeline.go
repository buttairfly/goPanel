package pixelpipe

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

// FramePipeline is a struct which defines a
type FramePipeline struct {
	running          bool
	rebuild          chan bool
	frameWg          *sync.WaitGroup
	internalSource   hardware.FrameSource
	outputFrameChan  chan hardware.Frame
	internalLastPipe *pipepart.Pipe
	pixelPipes       map[pipepart.ID]pipepart.PixelPiper
	lastPipeID       pipepart.ID
	firstPipeID      pipepart.ID
	prevID           pipepart.ID
	logger           *zap.Logger
}

// NewEmptyFramePipeline creates a new, empty FramePipeline which can hold multiple pipes end-to-end connected to each other
func NewEmptyFramePipeline(id pipepart.ID, logger *zap.Logger) *FramePipeline {
	if pipepart.IsPlaceholderID(id) {
		logger.Fatal("PipeIDPlaceholderError", zap.Error(pipepart.PipeIDPlaceholderError(id)))
	}
	pixelPipes := make(map[pipepart.ID]pipepart.PixelPiper)
	rebuild := make(chan bool)
	outputFrameChan := make(chan hardware.Frame)
	internalLastPipe := pipepart.NewPipe(id, outputFrameChan)

	return &FramePipeline{
		rebuild:          rebuild,
		frameWg:          new(sync.WaitGroup),
		internalLastPipe: internalLastPipe,
		internalSource:   nil,
		logger:           logger,
		pixelPipes:       pixelPipes,
		outputFrameChan:  outputFrameChan,
		firstPipeID:      pipepart.EmptyID,
		lastPipeID:       pipepart.EmptyID,
		prevID:           pipepart.EmptyID,
	}
}

// RunPipe implements PixelPiper interface
func (me *FramePipeline) RunPipe(destroyCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(me.internalLastPipe.GetFullOutput())
	me.running = true

	me.startPipePieces(me.frameWg)

	for {
		if rebuildInProgress := me.runInternalPipe(destroyCtx); rebuildInProgress {
			select {
			case <-destroyCtx.Done():
				me.drain()
				return
			case <-me.rebuild:
				me.running = true
				me.startPipePieces(me.frameWg)
				continue
			default:
				me.drain()
				me.frameWg.Wait()
				me.running = false
				me.rebuild = make(chan bool)
				me.frameWg = new(sync.WaitGroup)

			}
		} else {
			me.drain()
			return
		}

		//TODO wait until pipeline rebuild is ready
	}
}

func (me *FramePipeline) drain() {
	if !pipepart.IsEmptyID(me.firstPipeID) {
		fakeInput := make(chan hardware.Frame)
		me.pixelPipes[me.firstPipeID].SetInput(pipepart.ID("Drain"), fakeInput)
		close(fakeInput)
	}
}

func (me *FramePipeline) runInternalPipe(destroyCtx context.Context) bool {
	for {
		if me.internalSource == nil {
			return false
		}
		var sourceChan hardware.FrameSource
		if !pipepart.IsEmptyID(me.lastPipeID) {
			sourceChan = me.pixelPipes[me.lastPipeID].GetOutput(me.lastPipeID)
		} else {
			sourceChan = me.internalSource
		}
		select {
		case <-destroyCtx.Done():
			return false
		case sourceFrame, ok := <-sourceChan:
			if !ok {
				return true
			}
			me.outputFrameChan <- sourceFrame
		}
	}
}

func (me *FramePipeline) startPipePieces(wg *sync.WaitGroup) {
	wg.Add(len(me.pixelPipes))
	for _, pipe := range me.pixelPipes {
		go pipe.RunPipe(wg)
	}
}

// GetOutput implements PixelPiper interface
func (me *FramePipeline) GetOutput(id pipepart.ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.internalLastPipe.GetOutput(id)
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(pipepart.OutputIDMismatchError(me.GetID(), id)))
	return nil
}

// SetInput implements PixelPiper interface
func (me *FramePipeline) SetInput(prevID pipepart.ID, inputChan hardware.FrameSource) {
	if pipepart.IsEmptyID(prevID) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(pipepart.PipeIDEmptyError()))
	}
	me.internalSource = inputChan
	if me.firstPipeID == pipepart.EmptyID {
		me.internalLastPipe.SetInput(prevID, inputChan)
	} else {
		me.pixelPipes[me.firstPipeID].SetInput(prevID, inputChan)
	}
	me.prevID = prevID
}

// GetID implements PixelPiper interface
func (me *FramePipeline) GetID() pipepart.ID {
	return me.internalLastPipe.GetID()
}

// AddPipeBefore adds a pipe segment before id
func (me *FramePipeline) AddPipeBefore(id pipepart.ID, newPipe pipepart.PixelPiper) {
	if pipepart.IsEmptyID(id) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(pipepart.PipeIDEmptyError()))
	}
	if newPipe.GetID() == me.GetID() {
		me.logger.Fatal("PipeIDNotUniqueError", zap.Error(pipepart.PipeIDNotUniqueError(me.GetID())))
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

func (me *FramePipeline) addPipeBefore(id pipepart.ID, newPipe pipepart.PixelPiper) {

	// actually insert pipe or preplace itself

	afterPipe, ok := me.pixelPipes[id]
	if !ok {
		if id != me.GetID() {
			me.logger.Fatal("PipeIDMismatchError afterPipe", zap.Error(pipepart.PipeIDMismatchError(id, me.GetID())))
		} else {
			afterPipe = me.internalLastPipe
			me.lastPipeID = newPipe.GetID()
		}
	}
	prevID := afterPipe.GetPrevID()

	if !pipepart.IsEmptyID(prevID) {
		if me.prevID == prevID {
			newPipe.SetInput(prevID, me.internalSource)
		} else {
			var prevPipe pipepart.PixelPiper
			prevPipe, ok = me.pixelPipes[prevID]
			if !ok {
				me.logger.Fatal("PipeIDMismatchError prevPipe", zap.Error(pipepart.PipeIDMismatchError(prevID, me.prevID)))
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
