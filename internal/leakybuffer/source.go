package leakybuffer

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
)

var freeList = make(chan hardware.Frame, 100)
var outputChan = make(chan hardware.Frame)

// NewFrameSource produces a frame channel with a recycled frame or new one
func NewFrameSource(cancelCtx context.Context, tileConfigs hardware.TileConfigs, wg *sync.WaitGroup, logger *zap.Logger) {
	defer wg.Done()
	defer close(outputChan)
	for {
		var f hardware.Frame
		// Grab a buffer if available; allocate if not.
		select {
		case <-cancelCtx.Done():
			return
		case f = <-freeList:
			// Got one; nothing more to do.
		default:
			// None free, so allocate a new one.
			f = hardware.NewFrame(tileConfigs, logger)
		}
		outputChan <- f // Send to output.
	}
}

// GetFrameSource returns the frame producer chan
func GetFrameSource() hardware.FrameSource {
	return outputChan
}
