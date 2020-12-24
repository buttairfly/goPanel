package generatorpipe

import (
	"image/color"
	"math/rand"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/leakybuffer"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
	"github.com/buttairfly/goPanel/pkg/palette"
)

type whiteNoisePipe struct {
	pipe     *pipepart.Pipe
	newPixel int
	palette  palette.Palette
	picture  hardware.Frame
	logger   *zap.Logger
}

// WhiteNoisePipe generates for each tick interval a random pixel is drawn with a random color of the palette
func WhiteNoisePipe(
	id pipepart.ID,
	palette palette.Palette,
	newPixel int,
	logger *zap.Logger,
) pipepart.PixelPiper {
	if pipepart.IsPlaceholderID(id) {
		logger.Fatal("PipeIDPlaceholderError", zap.Error(pipepart.PipeIDPlaceholderError(id)))
	}
	outputChan := make(chan hardware.Frame)

	return &whiteNoisePipe{
		pipe:     pipepart.NewPipe(id, outputChan),
		newPixel: newPixel,
		palette:  palette,
		picture:  nil,
		logger:   logger,
	}
}

func (me *whiteNoisePipe) RunPipe(wg *sync.WaitGroup) {
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

func (me *whiteNoisePipe) GetID() pipepart.ID {
	return me.pipe.GetID()
}

func (me *whiteNoisePipe) GetPrevID() pipepart.ID {
	return me.pipe.GetPrevID()
}

func (me *whiteNoisePipe) GetOutput(id pipepart.ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.pipe.GetOutput(id)
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(pipepart.OutputIDMismatchError(me.GetID(), id)))
	return nil
}

func (me *whiteNoisePipe) SetInput(prevID pipepart.ID, inputChan hardware.FrameSource) {
	if pipepart.IsEmptyID(prevID) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(pipepart.PipeIDEmptyError()))
	}
	me.pipe.SetInput(prevID, inputChan)
}
