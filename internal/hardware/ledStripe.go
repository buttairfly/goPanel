package hardware

import (
	"go.uber.org/zap"
)

// LedStripe interface
type LedStripe interface {
	GetBuffer() []uint8
	GetPixelLength() int
	GetColorMap() map[string][]int
	GetAction() *LedStripeAction
}

type ledStripe struct {
	buffer          []uint8
	pixelLength     int
	numPixelChanges int
	pixelChangePos  []int
	fillType        FrameFillType
	logger          *zap.Logger
}

func (l *ledStripe) GetBuffer() []uint8 {
	return l.buffer
}

func (l *ledStripe) GetPixelLength() int {
	return l.pixelLength
}

// NewLedStripe creates a new led stripe buffer
func NewLedStripe(
	numLed, numPixelChanges int,
	pixelChangePos []int,
	fillType FrameFillType,
	logger *zap.Logger,
) LedStripe {
	bufferCap := numLed * NumBytePixel
	buffer := make([]uint8, bufferCap, bufferCap)
	return &ledStripe{
		buffer:          buffer,
		pixelLength:     numLed,
		numPixelChanges: numPixelChanges,
		pixelChangePos:  pixelChangePos,
		fillType:        fillType,
		logger:          logger,
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

func (l *ledStripe) GetAction() *LedStripeAction {
	switch l.fillType {
	case FillTypeFullFrame:
		return &LedStripeAction{
			change:    true,
			fullFrame: true,
		}
	case FillTypeSinglePixelChange:
		return &LedStripeAction{
			change:          true,
			otherDiffPixels: l.pixelChangePos,
		}
	case FillTypeSingleFillColor:
		return &LedStripeAction{
			change:    true,
			fillColor: NewPixelFromSlice(l.buffer, 0, l.logger),
		}
	default:
		return &LedStripeAction{}
	}
}
