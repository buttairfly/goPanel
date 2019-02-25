package hardware

// LedStripeCompare holds all diffenences between old and new LedStripe buffer
type LedStripeCompare struct {
	change          bool
	otherDiffPixels []int
}

// HasChanged returns true when a change happened
func (l *LedStripeCompare) HasChanged() bool {
	return l.change
}

// GetOtherDiffPixels returns the slice of changed pixels
func (l *LedStripeCompare) GetOtherDiffPixels() []int {
	return l.otherDiffPixels
}
