package leakybuffer

import (
	"testing"

	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

func TestSourceIntefaces(t *testing.T) {
	var _ pipepart.PixelPiper = (*Source)(nil)
}
