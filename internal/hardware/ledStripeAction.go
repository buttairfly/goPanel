package hardware

// LedStripeAction holds all diffenences between old and new LedStripe buffer
type LedStripeAction struct {
	change          bool
	fullFrame       bool
	fillColor       Pixel
	otherDiffPixels []int
}

// HasChanged returns true when a change happened
func (l *LedStripeAction) HasChanged() bool {
	return l.change
}

// IsFullFrame returns true, when a full frame update should be done
func (l *LedStripeAction) IsFullFrame() bool {
	return l.fullFrame
}

// GetFillColor returns nil when no fullColor is available, otherwise the Pixel color
func (l *LedStripeAction) GetFillColor() Pixel {
	return l.fillColor
}

// GetOtherDiffPixels returns the slice of changed pixels
func (l *LedStripeAction) GetOtherDiffPixels() []int {
	return l.otherDiffPixels
}
