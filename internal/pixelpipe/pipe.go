package pixelpipe

import (
	"sync"

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
		me.GetFullOutput() <- frame
	}
}

// GetOutput returns the framesource output
func (me *Pipe) GetOutput(id ID) hardware.FrameSource {
	return me.outputChan
}

// GetFullOutput returns the full output
func (me *Pipe) GetFullOutput() chan hardware.Frame {
	return me.outputChan
}

// SetInput sets the input
func (me *Pipe) SetInput(id ID, inputChan hardware.FrameSource) {
	me.prevID = id
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
