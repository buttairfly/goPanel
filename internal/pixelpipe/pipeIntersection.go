package pixelpipe

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
)

type simplePipeIntersection struct {
	id         ID
	prevIds    map[ID]ID
	inputs     map[ID]hardware.FrameSource
	emptyInput hardware.FrameSource
	outputs    map[ID]chan hardware.Frame
	logger     *zap.Logger
}

// NewSimplePipeIntersection creates a pipe intersection which wil put each frame from any input channel to all output channels
func NewSimplePipeIntersection(
	id ID,
	inputs map[ID]hardware.FrameSource,
	emptyInput hardware.FrameSource,
	numOutputChannels int,
	logger *zap.Logger,
) PixelPiper {
	if IsPlaceholderID(id) {
		logger.Fatal("PipeIDPlaceholderError", zap.Error(PipeIDPlaceholderError(id)))
	}
	outputs := make(map[ID]chan hardware.Frame)
	for num := 0; num < numOutputChannels; num++ {
		channelID := ID(fmt.Sprintf("%s_%d", id, num))
		outputs[channelID] = make(chan hardware.Frame)
	}
	return &simplePipeIntersection{
		id:         id,
		inputs:     inputs,
		outputs:    outputs,
		emptyInput: emptyInput,
		logger:     logger,
	}
}

func (me *simplePipeIntersection) RunPipe(wg *sync.WaitGroup) {
	defer wg.Done()
	defer me.close()
	subWg := new(sync.WaitGroup)
	subWg.Add(len(me.inputs))
	for id := range me.inputs {
		go me.runInput(id, subWg)
	}
	subWg.Wait()
}

func (me *simplePipeIntersection) runInput(id ID, wg *sync.WaitGroup) {
	defer wg.Done()
	for frame := range me.inputs[id] {
		currentOutputNum := 0
		outputFrame := frame
		for id, outputChan := range me.outputs {
			if currentOutputNum != 0 {
				var ok bool
				outputFrame, ok = <-me.emptyInput
				if !ok {
					return
				}
				outputFrame.CopyFromOther(frame)
			}
			select {
			case outputChan <- outputFrame:
				// all fine
			default:
				me.logger.Debug("simplePipeIntersection outputchan blocks", zap.String("outputId", string(id)), zap.String("id", string(me.GetID())))
			}
			currentOutputNum++
		}
	}
}

func (me *simplePipeIntersection) GetOutput(id ID) hardware.FrameSource {
	if output, ok := me.outputs[id]; ok {
		return output
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(OutputIDMismatchError(me.GetID(), id)))
	return nil
}

func (me *simplePipeIntersection) SetInput(prevID ID, inputChan hardware.FrameSource) {
	if IsEmptyID(prevID) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(PipeIDEmptyError()))
	}
	me.inputs[prevID] = inputChan
	me.prevIds[prevID] = prevID
}

func (me *simplePipeIntersection) GetID() ID {
	return me.id
}

func (me *simplePipeIntersection) GetPrevID() ID {
	// TODO fix function
	return me.id
}

func (me *simplePipeIntersection) close() {
	for _, output := range me.outputs {
		close(output)
	}
}
