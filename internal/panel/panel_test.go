package panel

import (
	"testing"

	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

func TestPixelPiperBasicInteface(t *testing.T) {
	var _ pipepart.PixelPiperBase = (*Panel)(nil)
	var _ pipepart.PixelPiperWithSubPipes = (*Panel)(nil)
}
