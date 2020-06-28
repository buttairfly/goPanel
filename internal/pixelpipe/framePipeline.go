package pixelpipe

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
)

// FramePipeline is a struct which defines a source
type FramePipeline struct {
	destroyCtx        context.Context
	running           bool
	rebuild           chan bool
	frameWg           *sync.WaitGroup
	internalInputChan chan hardware.Frame

	pipe       *Pipe
	logger     *zap.Logger
	pixelPipes map[ID]PixelPiper
	lastPipe   ID
}

// NewEmptyFramePipeline creates a new, empty FramePipeline
func NewEmptyFramePipeline(destroyCtx context.Context, id ID, logger *zap.Logger) *FramePipeline {
	pixelPipes := make(map[ID]PixelPiper)
	rebuild := make(chan bool)
	outputChan := make(chan hardware.Frame)

	return &FramePipeline{
		destroyCtx:        destroyCtx,
		rebuild:           rebuild,
		frameWg:           new(sync.WaitGroup),
		internalInputChan: outputChan,
		pipe:              NewPipe(id, outputChan),
		logger:            logger,
		pixelPipes:        pixelPipes,
		lastPipe:          EmptyID,
	}
}

// RunPipe implements PixelPiper interface
func (me *FramePipeline) RunPipe(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(me.pipe.outputChan)
	me.running = true

	subWg := new(sync.WaitGroup)
	defer subWg.Wait()

	subWg.Add(1)

	go me.processOutgoingFrames(context.TODO(), subWg)

	me.startPipePieces(subWg)

	for {
		if inputClosed := me.runInternalPipe(); inputClosed {
			return
		}
		me.frameWg = new(sync.WaitGroup)
		me.rebuild = make(chan bool)

		// wait until pipeline rebuild is ready
	}
}

func (me *FramePipeline) runInternalPipe() bool {
	defer me.frameWg.Wait()
	for {
		select {
		case <-me.rebuild:
			return false
		default:
		}
		if inputClosed := me.processIncommingFrame(); inputClosed {
			return true
		}
	}
}

func (me *FramePipeline) processIncommingFrame() bool {
	me.frameWg.Add(1)
	inputFrame, ok := <-me.pipe.GetInput()
	if !ok {
		close(me.internalInputChan)
		return true
	}
	me.internalInputChan <- inputFrame
	return false
}

func (me *FramePipeline) startPipePieces(wg *sync.WaitGroup) {
	wg.Add(len(me.pixelPipes))
	for _, pipe := range me.pixelPipes {
		go pipe.RunPipe(wg)
	}
}

func (me *FramePipeline) processOutgoingFrames(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for outputFrame := range me.pixelPipes[me.lastPipe].GetOutput(me.lastPipe) {
		me.pipe.GetFullOutput() <- outputFrame
		me.frameWg.Done()
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

// GetOutput implements PixelPiper interface
func (me *FramePipeline) GetOutput(id ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.pipe.GetOutput(id)
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(OutputIDMismatchError("simplePipeIntersection", me.GetID(), id)))
	return nil
}

// SetInput implements PixelPiper interface
func (me *FramePipeline) SetInput(inputID ID, inputChan hardware.FrameSource) {
	me.pipe.SetInput(inputID, inputChan)
}

// GetID implements PixelPiper interface
func (me *FramePipeline) GetID() ID {
	return me.pipe.GetID()
}

// AddPipeBefore adds a pipe segment before id
func (me *FramePipeline) AddPipeBefore(id ID, newPipe PixelPiper) {
	if me.running {
		// stop pipeline and wait until all frames are empty
		close(me.rebuild)
		me.frameWg.Wait()
	}
	me.addPipeBefore(id, newPipe)
}

func (me *FramePipeline) addPipeBefore(id ID, newPipe PixelPiper) {
	if newPipe.GetID() != id {
		// actually insert pipe
		basePipe, ok := me.pixelPipes[id]
		if !ok {
			basePipe = me.pipe
		}
		prevID := EmptyID
		if id != EmptyID {
			prevID = basePipe.GetPrevID()
		}
		if prevID != EmptyID {
			newPipe.SetInput(prevID, me.pixelPipes[prevID].GetOutput(newPipe.GetID()))
		} else {
			newPipe.SetInput(EmptyID, me.internalInputChan)
		}

		if me.lastPipe != EmptyID {
			basePipe.SetInput(newPipe.GetID(), newPipe.GetOutput(newPipe.GetID()))
		} else {
			me.pipe.SetInput(newPipe.GetID(), newPipe.GetOutput(newPipe.GetID()))
		}

		me.pixelPipes[newPipe.GetID()] = newPipe
		return
	}

	// replace pipe

}
