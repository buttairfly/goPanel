package generator

import (
	"context"
	"image/color"
	"sync"

	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/internal/hardware"
)

// LastBlackFrameFrameGenerator generates a black frame at the end of the program
// and closes the frame input chan on program exit
func LastBlackFrameFrameGenerator(
	cancelCtx context.Context,
	frame hardware.Frame,
	inputChan chan<- hardware.Frame,
	wg *sync.WaitGroup,
	logger *zap.Logger,
) {
	defer wg.Done()
	defer close(inputChan)

	// TODO: add leaky buffer recycling https://golang.org/doc/effective_go.html#leaky_buffer
	colorFrame := hardware.NewCopyFrameWithEmptyImage(frame)
	colorFrame.Fill(color.Black)
	select {
	case <-cancelCtx.Done():
		inputChan <- colorFrame
	}
}
