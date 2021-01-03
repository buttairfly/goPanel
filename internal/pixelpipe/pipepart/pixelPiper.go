package pipepart

import (
	"context"
	"sync"

	"github.com/buttairfly/goPanel/internal/hardware"
)

// PixelPiper is a interface to generate a pixelPipeline from several pipe segments
type PixelPiper interface {
	RunPipe(cancelCtx context.Context, wg *sync.WaitGroup)
	SetInput(inputID ID, inputChan hardware.FrameSource)
	GetOutput(outputID ID) hardware.FrameSource
	GetID() ID
	GetPrevID() ID
	Marshal() Marshal
	GetParams() []PipeParam
}
