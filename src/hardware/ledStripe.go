package hardware

import "log"

// LedStripe interface
type LedStripe interface {
	GetBuffer() []uint8
	GetPixelLength() int
	Compare(other LedStripe) *LedStripeCompare
}

type ledStripe struct {
	buffer      []uint8
	pixelLength int
}

func (l *ledStripe) GetBuffer() []uint8 {
	return l.buffer
}

func (l *ledStripe) GetPixelLength() int {
	return l.pixelLength
}

// NewLedStripe creates a new led stripe buffer
func NewLedStripe(numLed int) LedStripe {
	bufferCap := numLed * NumBytePixel
	buffer := make([]uint8, bufferCap, bufferCap)
	return &ledStripe{
		buffer:      buffer,
		pixelLength: numLed,
	}
}

func (l *ledStripe) Compare(other LedStripe) *LedStripeCompare {
	if l.pixelLength != other.GetPixelLength() || len(l.buffer) != len(l.GetBuffer()) {
		log.Fatal("Pixel length is not equal", l, other)
		return nil
	}
	change := false
	otherDiffPixels := make([]int, 0, l.pixelLength)
	for i := 0; i < l.pixelLength; i++ {
		lPix := NewPixelFromSlice(l.buffer, i)
		oPix := NewPixelFromSlice(other.GetBuffer(), i)
		if !lPix.Equals(oPix) {
			change = true
			otherDiffPixels = append(otherDiffPixels, i)
		}
	}
	return &LedStripeCompare{
		change:          change,
		otherDiffPixels: otherDiffPixels,
	}
}
