package generator

import (
	"image"
	"image/color"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/palette"
)

// FrameGenerator generates a new input frame stream
func FrameGenerator(frame hardware.Frame, inputChan chan<- hardware.Frame, wg *sync.WaitGroup, exitChan <-chan bool, logger *zap.Logger) {
	wg.Add(1)
	defer wg.Done()

	mainPicture := image.NewRGBA(frame.Bounds())

	colors := make([]color.Color, 0, 10)
	colors = append(colors, color.RGBA{0xff, 0, 0, 0xff})
	colors = append(colors, color.RGBA{0xff, 0xa5, 0, 0xff})
	const granularity int = 100
	const wrapping bool = true
	fader := palette.NewFader(colors, granularity, wrapping)
	increments := fader.GetIncrements()
	for {
		select {
		case _, ok := <-exitChan:
			if !ok {
				return
			}
		default:
			for _, increment := range increments {
				color := fader.Fade(increment)
				for y := 0; y < frame.GetHeight(); y++ {
					for x := 0; x < frame.GetWidth(); x++ {
						mainPicture.Set(x, y, color)
					}
				}
				// TODO: add leaky buffer recycling https://golang.org/doc/effective_go.html#leaky_buffer
				colorFrame := hardware.NewCopyFrameFromImage(frame, mainPicture)
				time.Sleep(5 * 1000 * time.Millisecond)
				inputChan <- colorFrame
				// TODO: frame counter logic
				// logger.Sugar().Infof("send frame %d", colorFrame.GetTime())
			}
		}
	}
}
