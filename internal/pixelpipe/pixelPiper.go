package pixelpipe

import (
	"sync"

	"github.com/buttairfly/goPanel/internal/hardware"
)

// PixelPiper is a interface to generate a pixelPipeline from several pipe segments
type PixelPiper interface {
	RunPipe(wg *sync.WaitGroup)
	SetInput(inputID ID, inputChan hardware.FrameSource)
	GetOutput(outputID ID) hardware.FrameSource
	GetID() ID
	GetPrevID() ID
}
