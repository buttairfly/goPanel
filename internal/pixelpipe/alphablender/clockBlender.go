package alphablender

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"strconv"
	"sync"
	"time"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
	"github.com/buttairfly/goPanel/pkg/imagefilter"
	"github.com/buttairfly/goPanel/pkg/intmath"

	"go.uber.org/zap"
)

type clockBlender struct {
	pipe        *pipepart.Pipe
	params      []pipepart.PipeParam
	minDimmer   float64
	maxDimmer   float64
	currentTime string
	blendFrame  *image.Alpha
	logger      *zap.Logger
}

const numClockDigits = 4
const clockDigitTimeDevider = 2

// NewClockBlender adds a clock over the current frame
func NewClockBlender(id pipepart.ID, minDimmer float64, maxDimmer float64, logger *zap.Logger) pipepart.PixelPiper {
	pipepart.CheckNoPlaceholderID(id, logger)
	outputChan := make(chan hardware.Frame)

	return &clockBlender{
		pipe:      pipepart.NewPipe(id, outputChan),
		minDimmer: minDimmer,
		maxDimmer: maxDimmer,
		params:    getParams(minDimmer, maxDimmer, logger),
		logger:    logger,
	}
}

func (me *clockBlender) RunPipe(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(me.pipe.GetFullOutput())

	for frame := range me.pipe.GetInput() {
		now := time.Now()
		actualTime := fmt.Sprintf("%02d%02d", now.Hour(), now.Minute())
		if actualTime != me.currentTime {
			me.logger.Info("actualtime", zap.String("time", actualTime))
			me.currentTime = actualTime
			me.changeBlendFrame(frame.Bounds())
		}

		frame.AlphaBlend(me.blendFrame)

		me.pipe.GetFullOutput() <- frame
	}
}

func (me *clockBlender) changeBlendFrame(bounds image.Rectangle) {
	me.blendFrame = image.NewAlpha(bounds)
	for numDigit, digit := range me.currentTime {
		digitNum, err := strconv.ParseInt(string(digit), 10, 64)
		if err != nil {
			me.logger.Fatal("failed to parse digit", zap.String("digit", string(digit)), zap.Error(err))
		}
		if digitNum > 9 || digitNum < 0 {
			me.logger.Fatal("digitNum out of range", zap.String("digit", string(digit)), zap.Int64("digitNum", digitNum))
		}
		digitFrame := imagefilter.Digits3x8[digitNum]

		r := me.mapDigitPosition(bounds, digitFrame.Rect, numDigit)
		for y := 0; y < r.Dy(); y++ {
			for x := 0; x < r.Dx(); x++ {
				alpha := digitFrame.AlphaAt(x, y)
				me.blendFrame.SetAlpha(r.Min.X+x, r.Min.Y+y, alpha)
			}
		}
	}
	me.logger.Debug("clock", zap.String("frame", fmt.Sprintf("%x", me.blendFrame.Pix)))
	me.applyDimmerAndInvert()
	me.logger.Debug("dimmer", zap.String("frame", fmt.Sprintf("%x", me.blendFrame.Pix)))
}

func (me *clockBlender) applyDimmerAndInvert() {
	if me.maxDimmer != 1.0 && me.minDimmer != 0.0 {
		for y := 0; y < me.blendFrame.Rect.Dy(); y++ {
			for x := 0; x < me.blendFrame.Rect.Dx(); x++ {
				currentValue := me.blendFrame.AlphaAt(x, y).A
				changeValue := uint8(intmath.Rescale(int(currentValue), 0x00, 0xff, int(me.minDimmer*255.0), int(me.maxDimmer*255.0)))

				if currentValue != changeValue {
					me.blendFrame.SetAlpha(x, y, color.Alpha{A: changeValue})
				}
			}
		}
	}
}

func (me *clockBlender) mapDigitPosition(bounds image.Rectangle, digitBounds image.Rectangle, numDigit int) image.Rectangle {
	if numDigit >= numClockDigits {
		me.logger.Fatal("numDigit exeeds numClockDigits", zap.Int("numDigit", numDigit), zap.Int("numClockDigits", numClockDigits))
	}
	digitWidth := bounds.Dx() / numClockDigits
	digitGap := bounds.Dx() % numClockDigits

	digitXOffset := numDigit * digitWidth
	if numDigit >= clockDigitTimeDevider {
		digitXOffset += digitGap
	}

	var xOffset int
	switch digitWidth {
	case 5:
		xOffset = 1
	default:
		me.logger.Fatal("digitWidth not implemented", zap.Int("digitWidth", digitWidth))
	}

	digitHeight := bounds.Dy()

	var yOffset int
	switch digitHeight {
	case 10:
		yOffset = 1
	default:
		me.logger.Fatal("digitHeight not implemented", zap.Int("digitHeight", digitHeight))
	}

	x := xOffset + digitXOffset
	y := yOffset
	b := image.Rect(x, y, x+digitBounds.Dx(), y+digitBounds.Dy())
	return b
}

func (me *clockBlender) GetID() pipepart.ID {
	return me.pipe.GetID()
}

func (me *clockBlender) GetPrevID() pipepart.ID {
	return me.pipe.GetPrevID()
}

func (me *clockBlender) Marshal() *pipepart.Marshal {
	return pipepart.MarshalFromPixelPiperInterface(me)
}

func (me *clockBlender) GetType() pipepart.PipeType {
	return pipepart.ClockBlender
}

// GetParams implements PixelPiper interface
func (me *clockBlender) GetParams() []pipepart.PipeParam {
	return me.params
}

func (me *clockBlender) GetOutput(id pipepart.ID) hardware.FrameSource {
	if id == me.GetID() {
		return me.pipe.GetOutput(id)
	}
	me.logger.Fatal("OutputIDMismatchError", zap.Error(pipepart.OutputIDMismatchError(me.GetID(), id)))
	return nil
}

func (me *clockBlender) SetInput(prevID pipepart.ID, inputChan hardware.FrameSource) {
	if pipepart.IsEmptyID(prevID) {
		me.logger.Fatal("PipeIDEmptyError", zap.Error(pipepart.PipeIDEmptyError()))
	}
	me.pipe.SetInput(prevID, inputChan)
}

func getParams(minDimmer, maxDimmer float64, logger *zap.Logger) []pipepart.PipeParam {
	minDimmer = checkDimmer("minDimmer", minDimmer, 0.0, logger)
	maxDimmer = checkDimmer("maxDimmer", maxDimmer, 1.0, logger)

	pipeParams := make([]pipepart.PipeParam, 2)

	pipeParams[0] = pipepart.PipeParam{
		Name:  "minDimmer",
		Type:  pipepart.Gauge0to1,
		Value: fmt.Sprintf("%f", minDimmer),
	}
	pipeParams[1] = pipepart.PipeParam{
		Name:  "maxDimmer",
		Type:  pipepart.Gauge0to1,
		Value: fmt.Sprintf("%f", maxDimmer),
	}
	return pipeParams
}

func checkDimmer(name string, dimmer float64, defaultValue float64, logger *zap.Logger) float64 {
	if dimmer > 1.0 || dimmer < 0.0 {
		logger.Warn("dimmer out of bounds, set to default", zap.String("name", name), zap.Float64("dimmer", dimmer), zap.Float64("default", defaultValue))
		dimmer = 0.0
	}
	return dimmer
}
