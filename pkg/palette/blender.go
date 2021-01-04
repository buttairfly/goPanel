package palette

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
)

// Blender defines the blending color space
type Blender int

// List of all avilable blenders
const (
	RGB Blender = iota
	HCL
	HSV
	Luv
	Lab
	unknown Blender = -1
)

var s = [...]string{"RGB", "HCL", "HSV", "Luv", "Lab"}

// Blend selects the correct Blender to blend between c1, c2 with distance t
func (me Blender) Blend(c1, c2 colorful.Color, t float64) colorful.Color {
	switch me {
	case RGB:
		return c1.BlendRgb(c2, t)
	case HCL:
		return c1.BlendHcl(c2, t)
	case HSV:
		return c1.BlendHsv(c2, t)
	case Luv:
		return c1.BlendLuv(c2, t)
	case Lab:
		return c1.BlendLab(c2, t)
	default:
		return colorful.Color{}
	}
}

func (me Blender) String() string {
	return s[me]
}

// FromString transforms a blender string to Blender
func FromString(b string) (Blender, error) {
	for i, val := range s {
		if val == b {
			return Blender(i), nil
		}
	}
	return unknown, fmt.Errorf("blender unknown: %s", b)
}
