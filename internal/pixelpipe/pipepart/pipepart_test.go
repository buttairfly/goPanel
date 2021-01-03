package pipepart

import "testing"

func TestPixelPiperInteface(t *testing.T) {
	var _ PixelPiper = (*Pipe)(nil)
}
