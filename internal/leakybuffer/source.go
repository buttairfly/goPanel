package leakybuffer

import (
	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

var source *Source
var freeList = make(chan hardware.Frame, 100)

// Source is a frame recycler and source of frames
type Source struct {
	tileConfigs hardware.TileConfigs
	outputChan  chan hardware.Frame
	logger      *zap.Logger
}

// NewFrameSource produces a frame channel with a recycled frame or new one
func NewFrameSource(tileConfigs hardware.TileConfigs, logger *zap.Logger) *Source {
	source = &Source{
		tileConfigs: tileConfigs,
		outputChan:  make(chan hardware.Frame),
		logger:      logger,
	}
	return source
}

// Run starts the Source
func (me *Source) Run() {
	defer close(me.outputChan)
	for {
		var f hardware.Frame
		// Grab a buffer if available; allocate if not.
		select {
		case f = <-freeList:
			// Got one; nothing more to do.
		default:
			// None free, so allocate a new one.
			f = hardware.NewFrame(me.tileConfigs, me.logger)
		}
		me.outputChan <- f // Send to output => will wait for ever
	}
}

// GetFrameSource returns the frame producer chan
func (me *Source) GetFrameSource() hardware.FrameSource {
	return me.outputChan
}

// GetFrameSource returns the frame producer chan from the global var source
func GetFrameSource() hardware.FrameSource {
	return source.outputChan
}

// GetID returns the frame producer chan
func (me *Source) GetID() pipepart.ID {
	return pipepart.SourceID
}
