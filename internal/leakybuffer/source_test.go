package leakybuffer

import (
	"testing"

	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

func TestPixelPiperSourceInteface(t *testing.T) {
	var _ pipepart.PixelPiperSource = (*Source)(nil)
}
