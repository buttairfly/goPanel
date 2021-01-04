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
	Marshal() *Marshal
	Clear() Palette
	GetID() ID
	GetColors() []fixColor
	SetBlender(blender string) error
}

type palette struct {
	id        ID
	blender   Blender
	fixColors []fixColor
}

type fixColor struct {
	color colorful.Color
	pos   float64
}

// NewPalette generates a new Palette
func NewPalette(id ID, blender Blender) Palette {
	me := &palette{
		id:      id,
		blender: blender,
	}
	me.Clear()
	return me
}

func (me *palette) Clear() Palette {
	fixColorSlice := make([]fixColor, 0, 0)
	me.fixColors = fixColorSlice
	return me
}

func (me *palette) GetID() ID {
	return me.id
}

func (me *palette) GetColors() []fixColor {
	return me.fixColors
}

func (me *palette) SetBlender(blender string) error {
	b, err := FromString(blender)
	if err != nil {
		return err
	}
	me.blender = b
	return nil
}

func (me *palette) addAt(c colorful.Color, pos float64) {
	pos = guaranteeBetween0And1(pos)
	me.fixColors = append(me.slice(), fixColor{color: c, pos: pos})
	sort.Sort(me)
}

func (me *palette) replaceAt(c colorful.Color, pos float64) error {
	index, err := me.getIndexFromPos(pos)
	if err != nil {
		return err
	}
	me.slice()[index].color = c
	return nil
}

func (me *palette) PutAt(c colorful.Color, pos float64) {
	if err := me.replaceAt(c, pos); err != nil {
		me.addAt(c, pos)
	}
}

func (me *palette) GetKeyColorAtPos(pos float64) (*colorful.Color, error) {
	index, err := me.getIndexFromPos(pos)
	if err != nil {
		return nil, err
	}
	return &(me.slice()[index].color), nil
}

func (me *palette) MoveAt(move ColorMoveMarshal) error {
	index, errFrom := me.getIndexFromPos(move.From)
	if errFrom != nil {
		return fmt.Errorf("move from %v", errFrom)
	}

	toPos := guaranteeBetween0And1(move.To)
	toIndex, errTo := me.getIndexFromPos(toPos)
	if errTo == nil {
		return fmt.Errorf("move to index %d already used toPos: %f (%+v)", toIndex, toPos, move)
	}
	me.slice()[index].pos = toPos
	sort.Sort(me)
	return nil
}

func (me *palette) DeleteAt(pos float64) error {
	index, err := me.getIndexFromPos(pos)
	if err != nil {
		return err
	}
	me.fixColors = append(me.slice()[:index], me.slice()[index+1:]...)
	return nil
}

// Blend will blend colors within the palaette within [0.0,1.0]
// when the palette does not start with 0.0 the first palette color value is returned
// when the palette does not end with a 1.0 pos value, the last value is returned
func (me *palette) Blend(pos float64) colorful.Color {
	if me.Len() == 0 {
		return colorful.Color{R: 0, G: 0, B: 0}
	}
	pos = guaranteeBetween0And1(pos)
	return me.getInterpolatedColorFor(pos)
}

// This is the meat of the gradient computation. It returns a HCL-blend between
// the two colors around `t`.
// Note: It relies heavily on the fact that the gradient keypoints are sorted.
func (me *palette) getInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < me.Len()-1; i++ {
		c1 := me.slice()[i]
		if c1.pos > t {
			// palette does not start at 0.0 and t < c0.pos
			return me.slice()[i].color
		}
		c2 := me.slice()[i+1]
		if c1.pos <= t && t <= c2.pos {
			// We are in between c1 and c2. Go blend them!
			t12 := (t - c1.pos) / (c2.pos - c1.pos)
			return me.blender.Blend(c1.color, c2.color, t12).Clamped()
		}
	}
	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return me.slice()[me.Len()-1].color
}

func (me *palette) getIndexFromPos(pos float64) (int, error) {
	pos = guaranteeBetween0And1(pos)
	for i, pc := range me.slice() {
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

func (me *palette) slice() []fixColor {
	return (me.fixColors)
}

func (me *palette) Len() int {
	return len(me.slice())
}
func (me *palette) Less(i, j int) bool {
	return me.slice()[i].pos < me.slice()[j].pos
}

func (me *palette) Swap(i, j int) {
	me.slice()[i], me.slice()[j] = me.slice()[j], me.slice()[i]
}
