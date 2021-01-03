package panel

import (
	"testing"

	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

func TestPixelPiperBasicInteface(t *testing.T) {
	var _ pipepart.PixelPiperBasic = (*Panel)(nil)
}
