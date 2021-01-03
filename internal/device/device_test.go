package device

import (
	"testing"

	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

func TestPixelPiperSinkInteface(t *testing.T) {
	var _ pipepart.PixelPiperSink = (*printDevice)(nil)
	var _ pipepart.PixelPiperSink = (*serialDevice)(nil)
}
