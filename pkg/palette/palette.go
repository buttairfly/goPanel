package palette

import (
	"fmt"
	"math"
	"sort"

	"github.com/lucasb-eyer/go-colorful"
)

// Palette is the palette interface
type Palette interface {
	sort.Interface
	Blend(pos float64) colorful.Color
	PutAt(c colorful.Color, pos float64)
	GetKeyColorAtPos(pos float64) (*colorful.Color, error)
	DeleteAt(pos float64) error
	MoveAt(move ColorMoveMarshal) error
	ToMarshal() Marshal
	Clear() Palette
}

type palette []paletteColor

type paletteColor struct {
	color colorful.Color
	pos   float64
}

// NewPalette generates a new Palette
func NewPalette() Palette {
	p := new(palette)
	return p.Clear()
}

func (p *palette) Clear() Palette {
	plaetteSlice := make(palette, 0, 0)
	p = &plaetteSlice
	return p
}

func (p *palette) addAt(c colorful.Color, pos float64) {
	pos = guaranteeBetween0And1(pos)
	*p = append(p.slice(), paletteColor{color: c, pos: pos})
	sort.Sort(p)
}

func (p *palette) replaceAt(c colorful.Color, pos float64) error {
	index, err := p.getIndexFromPos(pos)
	if err != nil {
		return err
	}
	p.slice()[index].color = c
	return nil
}

func (p *palette) PutAt(c colorful.Color, pos float64) {
	if err := p.replaceAt(c, pos); err != nil {
		p.addAt(c, pos)
	}
}

func (p *palette) GetKeyColorAtPos(pos float64) (*colorful.Color, error) {
	index, err := p.getIndexFromPos(pos)
	if err != nil {
		return nil, err
	}
	return &(p.slice()[index].color), nil
}

func (p *palette) MoveAt(move ColorMoveMarshal) error {
	index, errFrom := p.getIndexFromPos(move.From)
	if errFrom != nil {
		return fmt.Errorf("Move from %v", errFrom)
	}

	toPos := guaranteeBetween0And1(move.To)
	_, errTo := p.getIndexFromPos(toPos)
	if errTo == nil {
		return fmt.Errorf("Move to already used")
	}
	p.slice()[index].pos = toPos
	sort.Sort(p)
	return nil
}

func (p *palette) DeleteAt(pos float64) error {
	index, err := p.getIndexFromPos(pos)
	if err != nil {
		return err
	}
	*p = append(p.slice()[:index], p.slice()[index+1:]...)
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
		c1 := p.slice()[i]
		if c1.pos > t {
			// palette does not start at 0.0 and t < c0.pos
			return p.slice()[i].color
		}
		c2 := p.slice()[i+1]
		if c1.pos <= t && t <= c2.pos {
			// We are in between c1 and c2. Go blend them!
			t12 := (t - c1.pos) / (c2.pos - c1.pos)
			return c1.color.BlendHcl(c2.color, t12).Clamped()
		}
	}
	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return p.slice()[p.Len()-1].color
}

func (p *palette) getIndexFromPos(pos float64) (int, error) {
	pos = guaranteeBetween0And1(pos)
	for i, pc := range p.slice() {
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

func (p *palette) Marshal() Marshal {
	return p.ToMarshal()
}
