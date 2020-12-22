package generatorpipe

import (
	"image/color"
	"math/rand"
	"sync"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/leakybuffer"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
	"github.com/buttairfly/goPanel/pkg/palette"
)

type waveGenerator struct {
	pipe    *pipepart.Pipe
	ticker  *time.Ticker
	palette palette.Palette
	picture hardware.Frame
	logger  *zap.Logger
}

// WaveGenerator generates for each tick interval a random pixel is drawn with a random color of the palette
func WaveGenerator(
	id pipepart.ID,
	interval time.Duration,
	logger *zap.Logger,
) pipepart.PixelPiper {
	if pipepart.IsPlaceholderID(id) {
		logger.Fatal("PipeIDPlaceholderError", zap.Error(pipepart.PipeIDPlaceholderError(id)))
	}
	outputChan := make(chan hardware.Frame)

	//todo set palette via function
	p := palette.NewPalette()
	p.AddAt(colorful.Color{R: 0.1, G: 0, B: 0}, 0)
	p.AddAt(colorful.Color{R: 0.5, G: 0.1, B: 0}, 1.0/3)
	p.AddAt(colorful.Color{R: 0.3, G: 0, B: 0}, 2.0/3)
	p.AddAt(colorful.Color{R: 0.4, G: 0.1, B: 0}, 1.0)

	return &waveGenerator{
		pipe:    pipepart.NewPipe(id, outputChan),
		ticker:  time.NewTicker(interval),
		palette: p,
		picture: nil,
		logger:  logger,
	}
}

func (me *waveGenerator) RunPipe(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(me.pipe.GetFullOutput())
	defer me.ticker.Stop()

	me.picture = <-me.pipe.GetInput()
	defer leakybuffer.DumpFrame(me.picture)
	me.picture.Fill(color.Black)

	for frame := range me.pipe.GetInput() {
		select {
		case <-me.ticker.C:
			x := rand.Intn(me.picture.GetWidth())
			y := rand.Intn(me.picture.GetHeight())
			p := rand.Float64()
			c := me.palette.Blend(p)
			me.picture.Set(x, y, c)
			frame.CopyImageFromOther(me.picture)
			// TODO: frame counter logic
			me.pipe.GetFullOutput() <- frame
		}
	}
}

func (me *waveGenerator) GetID() pipepart.ID {
	return me.pipe.GetID()
}

func (me *waveGenerator) GetPrevID() pipepart.ID {
	return me.pipe.GetPrevID()
}

func (me *waveGenerator) GetOutput(id pipepart.ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.pipe.GetOutput(id)
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(pipepart.OutputIDMismatchError(me.GetID(), id)))
	return nil
}

func (me *waveGenerator) SetInput(prevID pipepart.ID, inputChan hardware.FrameSource) {
	if pipepart.IsEmptyID(prevID) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(pipepart.PipeIDEmptyError()))
	}
	me.pipe.SetInput(prevID, inputChan)
}
