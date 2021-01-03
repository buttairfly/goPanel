package pipepart

import (
	"context"
	"sync"

	"github.com/buttairfly/goPanel/internal/hardware"
)

// PixelPiper is a interface to generate a pixelPipeline from several pipe segments
type PixelPiper interface {
	PixelPiperSink
	PixelPiperSource
}

// PixelPiperSink is a basic sink PixelPiper
type PixelPiperSink interface {
	PixelPiperBasic
	SetInput(inputID ID, inputChan hardware.FrameSource)
	GetPrevID() ID
}

// PixelPiperSource is a basic source PixelPiper
type PixelPiperSource interface {
	PixelPiperBasic
	GetOutput(outputID ID) hardware.FrameSource
}

// PixelPiperBasic is a base PixelPiper
type PixelPiperBasic interface {
	RunPipe(cancelCtx context.Context, wg *sync.WaitGroup)
	GetID() ID
	Marshal() Marshal
	GetParams() []PipeParam
}
