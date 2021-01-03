package alphablender

import (
	"testing"

	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

func TestPixelPiperInteface(t *testing.T) {
	var _ pipepart.PixelPiper = (*clockBlender)(nil)
}
