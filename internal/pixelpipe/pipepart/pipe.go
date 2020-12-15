package pipepart

import (
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
)

// Pipe is a struct that defines a basic pipe
type Pipe struct {
	id         ID
	prevID     ID
	inputChan  hardware.FrameSource
	outputChan chan hardware.Frame
}

// NewPipe returns a new Pipe
func NewPipe(id ID, outputChan chan hardware.Frame) *Pipe {
	if IsPlaceholderID(id) {
		zap.L().Fatal("PipeIDPlaceholderError", zap.Error(PipeIDPlaceholderError(id)))
	}
	return &Pipe{
		id:         id,
		outputChan: outputChan,
	}
}

// RunPipe implements PixelPiper interface, but is not useable
func (me *Pipe) RunPipe(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(me.GetFullOutput())
	// input chan is closed
	for frame := range me.inputChan {
		me.applyFrame(frame)
	}
}

func (me *Pipe) applyFrame(frame hardware.Frame) {
	me.GetFullOutput() <- frame
}

// GetOutput returns the framesource output
func (me *Pipe) GetOutput(id ID) hardware.FrameSource {
	if id != me.GetID() {
		zap.L().Fatal("OutputIDMismatchError", zap.Error(OutputIDMismatchError(me.GetID(), id)))
	}
	return me.outputChan
}

// GetFullOutput returns the full output
func (me *Pipe) GetFullOutput() chan hardware.Frame {
	return me.outputChan
}

// SetInput sets the input
func (me *Pipe) SetInput(prevID ID, inputChan hardware.FrameSource) {
	if IsEmptyID(prevID) {
		zap.L().Fatal("PipeIDEmptyError", zap.Error(PipeIDEmptyError()))
	}
	me.prevID = prevID
	me.inputChan = inputChan
}

// GetInput gets the input
func (me *Pipe) GetInput() hardware.FrameSource {
	return me.inputChan
}

// GetPrevID gets the input
func (me *Pipe) GetPrevID() ID {
	return me.prevID
}

// GetID returns the id
func (me *Pipe) GetID() ID {
	return me.id
}
