package hardware

import (
	"go.uber.org/zap"
)

// LedStripe interface
type LedStripe interface {
	GetBuffer() []uint8
	GetPixelLength() int
	GetColorMap() map[string][]int
	Compare(other LedStripe) *LedStripeCompare
}

type ledStripe struct {
	buffer      []uint8
	pixelLength int
	logger      *zap.Logger
}

func (l *ledStripe) GetBuffer() []uint8 {
	return l.buffer
}

func (l *ledStripe) GetPixelLength() int {
	return l.pixelLength
}

// NewLedStripe creates a new led stripe buffer
func NewLedStripe(numLed int, logger *zap.Logger) LedStripe {
	bufferCap := numLed * NumBytePixel
	buffer := make([]uint8, bufferCap, bufferCap)
	return &ledStripe{
		buffer:      buffer,
		pixelLength: numLed,
		logger:      logger,
	}
}

func (l *ledStripe) GetColorMap() map[string][]int {
	colorMap := make(map[string][]int)
	for i := 1; i < l.pixelLength; i++ {
		pix := NewPixelFromSlice(l.buffer, i, l.logger)
		if posSlice, ok := colorMap[pix.Hex()]; ok {
			colorMap[pix.Hex()] = append(posSlice, i)
		} else {
			posSlice := make([]int, 1, l.pixelLength) // make posSlice big enough for one whole screen
			posSlice[0] = i
			colorMap[pix.Hex()] = posSlice
		}
	}
	return colorMap
}

func (l *ledStripe) Compare(
	other LedStripe,
) *LedStripeCompare {
	if l.pixelLength != other.GetPixelLength() || len(l.buffer) != len(l.GetBuffer()) {
		l.logger.Fatal("pixel length is not equal",
			zap.Int("thisStripe", l.pixelLength),
			zap.Int("otherStripe", other.GetPixelLength()),
		)
		return nil
	}
	/*
		change := false
		fullColor := Pixel(nil)
		colorMap := l.GetColorMap()
		maxColor := 0
		maxColorHex := ""
		for hexColor, posSlice := range colorMap {
			maxColor = intmath.Max(maxColor, len(posSlice))
			maxColorHex = hexColor
		}
		if maxColor > l.pixelLength/2 {
			_, ok := colorMap[maxColorHex]
			if ok {
				var err error
				fullColor, err = NewPixelFromHex(maxColorHex)
				if err != nil {
					l.logger.Fatal("pixel from hex", zap.Error(err))
				}
			}
			change = true
		}
		otherDiffPixels := make([]int, 0, l.pixelLength-maxColor+1)
		oPix := fullColor
		for i := 0; i < l.pixelLength; i++ {
			lPix := NewPixelFromSlice(l.buffer, i, l.logger)
			if fullColor == nil {
				oPix = NewPixelFromSlice(other.GetBuffer(), i, l.logger)
			}
			if !lPix.Equals(oPix) {
				change = true
				otherDiffPixels = append(otherDiffPixels, i)
			}
		}
		return &LedStripeCompare{
			change:          change,
			fullColor:       fullColor,
			otherDiffPixels: otherDiffPixels,
		}
	*/
	// very dumb approach
	return &LedStripeCompare{
		change:    true,
		fullFrame: true,
	}
}
