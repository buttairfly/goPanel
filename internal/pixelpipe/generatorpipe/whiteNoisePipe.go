package generatorpipe

import (
	"math/rand"
	"sync"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
	"github.com/buttairfly/goPanel/pkg/palette"
)

type whiteNoisePipe struct {
	pipe    *pipepart.Pipe
	ticker  *time.Ticker
	palette palette.Palette
	logger  *zap.Logger
}

// WhiteNoisePipe generates for each tick interval a random pixel is drawn with a random color of the palette
func WhiteNoisePipe(
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
	p.AddAt(colorful.Color{R: 0xff, G: 0, B: 0}, 0)
	p.AddAt(colorful.Color{R: 0, G: 0xff, B: 0}, 0.3)
	p.AddAt(colorful.Color{R: 0, G: 0, B: 0xff}, 0.6)

	return &whiteNoisePipe{
		pipe:    pipepart.NewPipe(id, outputChan),
		ticker:  time.NewTicker(interval),
		palette: p,
		logger:  logger,
	}
}

func (me *whiteNoisePipe) RunPipe(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(me.pipe.GetFullOutput())
	defer me.ticker.Stop()
	for frame := range me.pipe.GetInput() {
		select {
		case <-me.ticker.C:
			x := rand.Intn(frame.GetWidth())
			y := rand.Intn(frame.GetHeight())
			p := rand.Float64()
			c := me.palette.Blend(p)
			frame.Set(x, y, c)
			me.logger.Sugar().Infof("output pixel x %d y %d c %v p %1.3f", x, y, c.Hex(), p)
			// TODO: frame counter logic
			me.pipe.GetFullOutput() <- frame
		}
	}
}

func (me *whiteNoisePipe) SetTickPeriod(period time.Duration) {
	me.ticker.Reset(period)
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
