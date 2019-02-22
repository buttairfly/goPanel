package hardware

// LedStripe interface
type LedStripe interface {
	GetBuffer() []uint8
}

type ledStripe struct {
	buffer      []uint8
	pixelLength int
}

func (l *ledStripe) GetBuffer() []uint8 {
	return l.buffer
}
