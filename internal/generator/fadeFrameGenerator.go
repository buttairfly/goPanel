package generator

import (
	"context"
	"image"
	"image/color"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/palette"
)

// FrameGenerator generates a new input frame stream
func FrameGenerator(
	cancelCtx context.Context,
	frame hardware.Frame,
	inputChan chan<- hardware.Frame,
	wg *sync.WaitGroup,
	logger *zap.Logger,
) {
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
		for _, increment := range increments {
			color := fader.Fade(increment)
			for y := 0; y < frame.GetHeight(); y++ {
				for x := 0; x < frame.GetWidth(); x++ {
					// TODO: add leaky buffer recycling https://golang.org/doc/effective_go.html#leaky_buffer
					colorFrame := hardware.NewCopyFrameFromImage(frame, mainPicture)
					colorFrame.Set(x, y, color)
					// TODO: change NewCopyFrameFromImage to check on a single pixel frame so mainPicture does not need to get changed
					mainPicture.Set(x, y, color)

					select {
					case <-cancelCtx.Done():
						return
					default:
						// TODO: frame counter logic
						inputChan <- colorFrame
					}
				}
			}

		}
	}
}
