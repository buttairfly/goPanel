package generatorpipe

import (
	"context"
	"fmt"
	"image/color"
	"math/rand"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/leakybuffer"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
	"github.com/buttairfly/goPanel/pkg/palette"
)

type whiteNoiseGenerator struct {
	pipe     *pipepart.Pipe
	newPixel int
	palette  palette.Palette
	picture  hardware.Frame
	logger   *zap.Logger
}

// WhiteNoiseGenerator generates for each tick interval a random pixel is drawn with a random color of the palette
func WhiteNoiseGenerator(
	id pipepart.ID,
	palette palette.Palette,
	newPixel int,
	logger *zap.Logger,
) pipepart.PixelPiper {
	pipepart.CheckNoPlaceholderID(id, logger)
	outputChan := make(chan hardware.Frame)

	return &whiteNoiseGenerator{
		pipe:     pipepart.NewPipe(id, outputChan),
		newPixel: newPixel,
		palette:  palette,
		picture:  nil,
		logger:   logger,
	}
}

func (me *whiteNoiseGenerator) RunPipe(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(me.pipe.GetFullOutput())

	me.picture = <-me.pipe.GetInput()
	defer leakybuffer.DumpFrame(me.picture)
	me.picture.Fill(color.Black)

	for frame := range me.pipe.GetInput() {
		me.picture.SetFillTypeDoNothing()
		for n := 0; n < me.newPixel; n++ {
			x := rand.Intn(me.picture.GetWidth())
			y := rand.Intn(me.picture.GetHeight())
			p := rand.Float64()
			c := me.palette.Blend(p)
			me.picture.Set(x, y, c)
		}
		frame.CopyImageFromOther(me.picture)
		// TODO: frame counter logic
		me.pipe.GetFullOutput() <- frame
	}
}

func (me *whiteNoiseGenerator) GetID() pipepart.ID {
	return me.pipe.GetID()
}

func (me *whiteNoiseGenerator) GetPrevID() pipepart.ID {
	return me.pipe.GetPrevID()
}

func (me *whiteNoiseGenerator) Marshal() pipepart.Marshal {
	return pipepart.Marshal{
		ID:     me.GetID(),
		PrevID: me.GetPrevID(),
		Params: me.GetParams(),
	}
}

// GetParams implements PixelPiper interface
func (me *whiteNoiseGenerator) GetParams() []pipepart.PipeParam {
	pp := make([]pipepart.PipeParam, 2)
	pp[0] = pipepart.PipeParam{
		Name:  "palette",
		Type:  pipepart.NameID,
		Value: me.palette.GetName(),
	}
	pp[1] = pipepart.PipeParam{
		Name:  "newPixel",
		Type:  pipepart.UInteger,
		Value: fmt.Sprintf("%d", me.newPixel),
	}
	return pp
}

func (me *whiteNoiseGenerator) GetOutput(id pipepart.ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.pipe.GetOutput(id)
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(pipepart.OutputIDMismatchError(me.GetID(), id)))
	return nil
}

func (me *whiteNoiseGenerator) SetInput(prevID pipepart.ID, inputChan hardware.FrameSource) {
	if pipepart.IsEmptyID(prevID) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(pipepart.PipeIDEmptyError()))
	}
	me.pipe.SetInput(prevID, inputChan)
}
