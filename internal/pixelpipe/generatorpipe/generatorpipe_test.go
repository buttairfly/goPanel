package generatorpipe

import (
	"testing"

	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

func TestPixelPiperInteface(t *testing.T) {
	var _ pipepart.PixelPiper = (*drawGenerator)(nil)
	var _ pipepart.PixelPiper = (*fullFrameFadeGenerator)(nil)
	var _ pipepart.PixelPiper = (*lastBlackFrameGenerator)(nil)
	var _ pipepart.PixelPiper = (*rainbowGenerator)(nil)
	var _ pipepart.PixelPiper = (*snakeGenerator)(nil)
	var _ pipepart.PixelPiper = (*whiteNoiseGenerator)(nil)
}
