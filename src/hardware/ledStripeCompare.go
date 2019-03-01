package hardware

// LedStripeCompare holds all diffenences between old and new LedStripe buffer
type LedStripeCompare struct {
	change          bool
	fullColor       Pixel
	otherDiffPixels []int
}

// HasChanged returns true when a change happened
func (l *LedStripeCompare) HasChanged() bool {
	return l.change
}

// GetFullColor returns nil when no fullColor is available, otherwise the Pixel
func (l *LedStripeCompare) GetFullColor() Pixel {
	return l.fullColor
}

// GetOtherDiffPixels returns the slice of changed pixels
func (l *LedStripeCompare) GetOtherDiffPixels() []int {
	return l.otherDiffPixels
}
