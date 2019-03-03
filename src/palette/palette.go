package palette

import (
	"math"
	"sort"

	"github.com/lucasb-eyer/go-colorful"
)

// Palette is the palette interface
type Palette interface {
	sort.Interface
	Add(c colorful.Color, pos float64)
	Blend(pos float64) colorful.Color
}

type palette []paletteColor

type paletteColor struct {
	color colorful.Color
	pos   float64
}

// NewPalette generates a new Palette
func NewPalette() Palette {
	p := make(palette, 0, 0)
	return &p
}

func (p *palette) Add(c colorful.Color, pos float64) {
	pos = guaranteeBetween0And1(pos)
	*p = append(*p, paletteColor{color: c, pos: pos})
	sort.Sort(p)
}

func (p *palette) Blend(pos float64) colorful.Color {
	if p.Len() == 0 {
		return colorful.Color{R: 0, G: 0, B: 0}
	}
	pos = guaranteeBetween0And1(pos)
	return p.getInterpolatedColorFor(pos)
}

// This is the meat of the gradient computation. It returns a HCL-blend between
// the two colors around `t`.
// Note: It relies heavily on the fact that the gradient keypoints are sorted.
func (p *palette) getInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < len(*p)-1; i++ {
		c1 := (*p)[i]
		c2 := (*p)[i+1]
		if c1.pos <= t && t <= c2.pos {
			// We are in between c1 and c2. Go blend them!
			t = (t - c1.pos) / (c2.pos - c1.pos)
			return c1.color.BlendHcl(c2.color, t).Clamped()
		}
	}

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return (*p)[len(*p)-1].color
}

func guaranteeBetween0And1(pos float64) float64 {
	if math.IsNaN(pos) {
		return 0.0
	}
	if pos < 0.0 {
		pos += math.Trunc(-pos) + 1.0
	}
	if pos > 1.0 {
		pos -= math.Trunc(pos)
	}
	return pos
}

func (p *palette) Len() int {
	return len(*p)
}
func (p *palette) Less(i, j int) bool {
	return (*p)[i].pos < (*p)[j].pos
}

func (p *palette) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}
