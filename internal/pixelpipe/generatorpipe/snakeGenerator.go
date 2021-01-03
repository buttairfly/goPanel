package generatorpipe

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
	"github.com/buttairfly/goPanel/pkg/palette"
)

type snakeGenerator struct {
	pipe      *pipepart.Pipe
	palette   palette.Palette
	picture   hardware.Frame
	wheelPos  float64
	snakePos  int
	colorDiff float64
	logger    *zap.Logger
}

// SnakeGenerator generates for each frame a new snake part progressing through the color palette
// the color palette should be circular to avoid hard color changes
func SnakeGenerator(
	id pipepart.ID,
	palette palette.Palette,
	colorPosDiff float64,
	logger *zap.Logger,
) pipepart.PixelPiper {

	pipepart.CheckNoPlaceholderID(id, logger)
	outputChan := make(chan hardware.Frame)

	return &snakeGenerator{
		pipe:      pipepart.NewPipe(id, outputChan),
		palette:   palette,
		colorDiff: colorPosDiff,
		logger:    logger,
	}
}

func (me *snakeGenerator) RunPipe(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(me.pipe.GetFullOutput())
	me.picture = <-me.pipe.GetInput()

	for frame := range me.pipe.GetInput() {
		Ny := frame.GetHeight()
		Nx := frame.GetWidth()

		x := me.snakePos % Nx
		y := (me.snakePos / Nx) % Ny
		if y%2 == 1 {
			x = Nx - x - 1
		}
		c := me.palette.Blend(me.wheelPos)
		me.picture.SetFillTypeDoNothing()
		me.picture.Set(x, y, c)
		frame.CopyImageFromOther(me.picture)

		me.wheelPos += me.colorDiff
		if me.wheelPos > 1.0 {
			me.wheelPos -= 1.0
		}
		if me.wheelPos < 0 {
			me.wheelPos += 1.0
		}
		me.snakePos++
		if me.snakePos > Nx*Ny {
			me.snakePos = 0
		}

		// TODO: frame counter logic
		me.pipe.GetFullOutput() <- frame
	}
}

func (me *snakeGenerator) GetID() pipepart.ID {
	return me.pipe.GetID()
}

func (me *snakeGenerator) GetPrevID() pipepart.ID {
	return me.pipe.GetPrevID()
}

func (me *snakeGenerator) GetOutput(id pipepart.ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.pipe.GetOutput(id)
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(pipepart.OutputIDMismatchError(me.GetID(), id)))
	return nil
}

func (me *snakeGenerator) Marshal() pipepart.Marshal {
	return pipepart.Marshal{
		ID:     me.GetID(),
		PrevID: me.GetPrevID(),
		Params: me.GetParams(),
	}
}

// GetParams implements PixelPiper interface
func (me *snakeGenerator) GetParams() []pipepart.PipeParam {
	pp := make([]pipepart.PipeParam, 2)
	pp[0] = pipepart.PipeParam{
		Name:  "palette",
		Type:  pipepart.NameID,
		Value: me.palette.GetName(),
	}
	pp[1] = pipepart.PipeParam{
		Name:  "colorDiff",
		Type:  pipepart.Gauge0to1,
		Value: fmt.Sprintf("%f", me.colorDiff),
	}
	return pp
}

func (me *snakeGenerator) SetInput(prevID pipepart.ID, inputChan hardware.FrameSource) {
	if pipepart.IsEmptyID(prevID) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(pipepart.PipeIDEmptyError()))
	}
	me.pipe.SetInput(prevID, inputChan)
}
