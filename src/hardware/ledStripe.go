package hardware

// LedStripe interface
type LedStripe interface {
}

type ledStripe struct {
	buffer      []uint8
	pixelLength int
}
