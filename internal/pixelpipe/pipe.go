package pixelpipe

import (
	"github.com/buttairfly/goPanel/internal/hardware"
)

// Pipe is a struct that defines a basic pipe
type Pipe struct {
	id         ID
	inputChan  hardware.FrameSource
	outputChan chan hardware.Frame
}

func NewPipe(id ID, outputChan chan hardware.Frame) *Pipe {
	return &Pipe{
		id:         id,
		outputChan: outputChan,
	}
}

func (me *Pipe) GetOutput() chan hardware.Frame {
	return me.outputChan
}

func (me *Pipe) SetInput(inputChan hardware.FrameSource) {
	me.inputChan = inputChan
}

func (me *Pipe) GetInput() hardware.FrameSource {
	return me.inputChan
}

func (me *Pipe) GetID() ID {
	return me.id
}
