package palette

import "image/color"

// Color is the enhanced color functionality
type Color interface {
	color.Color
	Equals(c color.Color) bool
}
