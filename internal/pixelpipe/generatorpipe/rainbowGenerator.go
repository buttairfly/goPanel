package generatorpipe

import (
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
	"github.com/buttairfly/goPanel/pkg/palette"
)

type rainbowGenerator struct {
	pipe     *pipepart.Pipe
	palette  palette.Palette
	wheelPos float64
	dx       float64
	dy       float64
	logger   *zap.Logger
}

// RainbowGenerator generates for each tick interval a progressing rainbow through the color palette
// the color palette should be circular to avoid hard color changes
func RainbowGenerator(
	id pipepart.ID,
	palette palette.Palette,
	dx float64,
	dy float64,
	logger *zap.Logger,
) pipepart.PixelPiper {
	if pipepart.IsPlaceholderID(id) {
		logger.Fatal("PipeIDPlaceholderError", zap.Error(pipepart.PipeIDPlaceholderError(id)))
	}
	outputChan := make(chan hardware.Frame)

	return &rainbowGenerator{
		pipe:    pipepart.NewPipe(id, outputChan),
		palette: palette,
		dx:      dx,
		dy:      dy,
		logger:  logger,
	}
}

func (me *rainbowGenerator) RunPipe(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(me.pipe.GetFullOutput())

	for frame := range me.pipe.GetInput() {
		Ny := frame.GetHeight()
		Nx := frame.GetWidth()

		wheelPos00 := me.wheelPos
		for y := 0; y < Ny; y++ {
			wheelPos0Y := wheelPos00 - float64(y)*me.dy
			for x := 0; x < Nx; x++ {
				wheelPosXY := wheelPos0Y - float64(x)*me.dx
				wheelPosXY = trimWheelPos(wheelPosXY)
				c := me.palette.Blend(wheelPosXY)
				frame.Set(x, y, c)
			}
		}
		me.wheelPos += me.dx
		me.wheelPos = trimWheelPos(me.wheelPos)

		// TODO: frame counter logic
		me.pipe.GetFullOutput() <- frame
	}
}

func trimWheelPos(w float64) float64 {
	if w > 1.0 {
		return w - 1.0
	}
	if w < 0 {
		return w + 1.0
	}
	return w
}

func (me *rainbowGenerator) GetID() pipepart.ID {
	return me.pipe.GetID()
}

func (me *rainbowGenerator) GetPrevID() pipepart.ID {
	return me.pipe.GetPrevID()
}

func (me *rainbowGenerator) GetOutput(id pipepart.ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.pipe.GetOutput(id)
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(pipepart.OutputIDMismatchError(me.GetID(), id)))
	return nil
}

func (me *rainbowGenerator) SetInput(prevID pipepart.ID, inputChan hardware.FrameSource) {
	if pipepart.IsEmptyID(prevID) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(pipepart.PipeIDEmptyError()))
	}
	me.pipe.SetInput(prevID, inputChan)
}
