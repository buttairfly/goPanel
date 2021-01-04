package pipepart

import (
	"context"
	"sync"

	"github.com/buttairfly/goPanel/internal/hardware"
)

// PixelPiper is a interface to generate a pixelPipeline from several pipe segments
type PixelPiper interface {
	PixelPiperBase
	pixelPiperSink
	pixelPiperSource
}

// PixelPiperSink is a basic sink PixelPiper
type PixelPiperSink interface {
	PixelPiperBase
	pixelPiperSink
}

type pixelPiperSink interface {
	SetInput(inputID ID, inputChan hardware.FrameSource)
	GetPrevID() ID
}

// PixelPiperSource is a basic source PixelPiper
type PixelPiperSource interface {
	PixelPiperBase
	pixelPiperSource
}

type pixelPiperSource interface {
	GetOutput(outputID ID) hardware.FrameSource
}

// PixelPiperBase is a base PixelPiper
type PixelPiperBase interface {
	pixelPiperMarshaller
	RunPipe(cancelCtx context.Context, wg *sync.WaitGroup)
	GetType() PipeType
	GetParams() []PipeParam
	// SetParam(param PipeParam) ([]PipeParam, error)
}

// PixelPiperAddableSubPipe is a PixelPiperWithSubPipes which can add new PixelPiper
type PixelPiperAddableSubPipe interface {
	PixelPiperWithSubPipes
	AddPipeBefore(id ID, newPipe PixelPiper)
}

// PixelPiperWithSubPipes is a PixelPiper with SubPipes
type PixelPiperWithSubPipes interface {
	pixelPiperMarshaller
	GetPipes() PipesMarshal
	GetPipeByID(id ID) (PixelPiper, error)
}

type pixelPiperMarshaller interface {
	Marshal() *Marshal
	GetID() ID
}
