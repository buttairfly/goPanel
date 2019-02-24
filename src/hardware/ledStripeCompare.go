package hardware

type LedStripeCompare struct {
	change bool
	//pixelCommonnessMap map[Pixel]int
	//mainColorDiffPixel int
	otherDiffPixels []int
}

func (l *LedStripeCompare) HasChanged() bool {
	return l.change
}
func (l *LedStripeCompare) GetOtherDiffPixels() []int {
	return l.otherDiffPixels
}
