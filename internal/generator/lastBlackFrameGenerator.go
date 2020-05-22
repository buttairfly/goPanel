package generator

import (
	"image"
	"image/color"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
)

// LastBlackFrameFrameGenerator generates a black frame at the end of the program
func LastBlackFrameFrameGenerator(
	frame hardware.Frame,
	inputChan chan<- hardware.Frame,
	wg *sync.WaitGroup,
	exitChan <-chan bool,
	logger *zap.Logger,
) {
	wg.Add(1)
	defer wg.Done()
	defer close(inputChan)

	mainPicture := image.NewRGBA(frame.Bounds())
	color := color.Black
	for y := 0; y < frame.GetHeight(); y++ {
		for x := 0; x < frame.GetWidth(); x++ {
			mainPicture.Set(x, y, color)
		}
	}
	// TODO: add leaky buffer recycling https://golang.org/doc/effective_go.html#leaky_buffer
	colorFrame := hardware.NewCopyFrameFromImage(frame, mainPicture)
	for {
		select {
		case _, ok := <-exitChan:
			if !ok {
				time.Sleep(500 * time.Millisecond)
				inputChan <- colorFrame
				return
			}
		default:
			// TODO: frame counter logic
			// logger.Sugar().Infof("send frame %d", colorFrame.GetTime())
			time.Sleep(50 * time.Millisecond)
		}
	}

}
