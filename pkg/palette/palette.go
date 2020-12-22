package palette

import (
	"math"
	"sort"

	"github.com/lucasb-eyer/go-colorful"
	"go.uber.org/zap"
)

// Palette is the palette interface
type Palette interface {
	sort.Interface
	Blend(pos float64) colorful.Color
	AddAt(c colorful.Color, pos float64)
	PutAt(c colorful.Color, pos float64)
	ReplaceAt(c colorful.Color, pos float64) error
	GetKeyColorAtPos(pos float64) (*colorful.Color, error)
	DeleteAt(pos float64) error
	MoveAt(pos, toPos float64) error
	Clear()
}

type palette []paletteColor

type paletteColor struct {
	color colorful.Color
	pos   float64
}

// NewPalette generates a new Palette
func NewPalette() Palette {
	p := new(palette)
	p.Clear()
	return p
}

func (p *palette) Clear() {
	plaetteSlice := make(palette, 0, 0)
	p = &plaetteSlice
}

func (p *palette) AddAt(c colorful.Color, pos float64) {
	pos = guaranteeBetween0And1(pos)
	*p = append(*p, paletteColor{color: c, pos: pos})
	sort.Sort(p)
}

func (p *palette) ReplaceAt(c colorful.Color, pos float64) error {
	index, err := p.getIndexFromPos(pos)
	if err != nil {
		return err
	}
	(*p)[index].color = c
	return nil
}

func (p *palette) PutAt(c colorful.Color, pos float64) {
	if err := p.ReplaceAt(c, pos); err != nil {
		p.AddAt(c, pos)
	}
}

func (p *palette) GetKeyColorAtPos(pos float64) (*colorful.Color, error) {
	index, err := p.getIndexFromPos(pos)
	if err != nil {
		return nil, err
	}
	return &((*p)[index].color), nil
}

func (p *palette) MoveAt(pos, toPos float64) error {
	index, err := p.getIndexFromPos(pos)
	if err != nil {
		return err
	}

	toPos = guaranteeBetween0And1(toPos)
	(*p)[index].pos = toPos
	sort.Sort(p)
	return nil
}

func (p *palette) DeleteAt(pos float64) error {
	index, err := p.getIndexFromPos(pos)
	if err != nil {
		return err
	}
	*p = append((*p)[:index], (*p)[index+1:]...)
	return nil
}

// Blend will blend colors within the palaette within [0.0,1.0]
// when the palette does not start with 0.0 the first palette color value is returned
// when the palette does not end with a 1.0 pos value, the last value is returned
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
	for i := 0; i < p.Len()-1; i++ {
		c1 := (*p)[i]
		if c1.pos > t {
			// palette does not start at 0.0 and t < c0.pos
			return p.slice()[i].color
		}
		c2 := (*p)[i+1]
		if c1.pos <= t && t <= c2.pos {
			// We are in between c1 and c2. Go blend them!
			t12 := (t - c1.pos) / (c2.pos - c1.pos)
			zap.L().Info("palette", zap.Float64("t12", t12))
			return c1.color.BlendHcl(c2.color, t12).Clamped()
		}
	}
	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return p.slice()[p.Len()-1].color
}

func (p *palette) getIndexFromPos(pos float64) (int, error) {
	pos = guaranteeBetween0And1(pos)
	for i, pc := range *p {
		if pc.pos == pos {
			return i, nil
		}
	}
	return -1, NoKeyColorFoundError(pos)
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

func (p *palette) slice() palette {
	return (*p)
}

func (p *palette) Len() int {
	return len(p.slice())
}
func (p *palette) Less(i, j int) bool {
	return p.slice()[i].pos < p.slice()[j].pos
}

func (p *palette) Swap(i, j int) {
	p.slice()[i], p.slice()[j] = p.slice()[j], p.slice()[i]
}
