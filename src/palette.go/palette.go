package palette

// Palette is a palette of colors.
type Palette []Color

// Convert returns the palette color closest to c in Euclidean R,G,B space.
func (p Palette) Convert(c Color) Color {
	if len(p) == 0 {
		return nil
	}
	return p[p.Index(c)]
}

// Index returns the index of the palette color closest to c in Euclidean
// R,G,B,A space.
func (p Palette) Index(c Color) int {
	// A batch version of this computation is in image/draw/draw.go.

	cr, cg, cb, ca := c.RGBA()
	ret, bestSum := 0, uint32(1<<32-1)
	for i, v := range p {
		vr, vg, vb, va := v.RGBA()
		sum := sqDiff(cr, vr) + sqDiff(cg, vg) + sqDiff(cb, vb) + sqDiff(ca, va)
		if sum < bestSum {
			if sum == 0 {
				return i
			}
			ret, bestSum = i, sum
		}
	}
	return ret
}

// sqDiff returns the squared-difference of x and y, shifted by 2 so that
// adding four of those won't overflow a uint32.
//
// x and y are both assumed to be in the range [0, 0xffff].
func sqDiff(x, y uint32) uint32 {
	// The canonical code of this function looks as follows:
	//
	//	var d uint32
	//	if x > y {
	//		d = x - y
	//	} else {
	//		d = y - x
	//	}
	//	return (d * d) >> 2
	//
	// Language spec guarantees the following properties of unsigned integer
	// values operations with respect to overflow/wrap around:
	//
	// > For unsigned integer values, the operations +, -, *, and << are
	// > computed modulo 2n, where n is the bit width of the unsigned
	// > integer's type. Loosely speaking, these unsigned integer operations
	// > discard high bits upon overflow, and programs may rely on ``wrap
	// > around''.
	//
	// Considering these properties and the fact that this function is
	// called in the hot paths (x,y loops), it is reduced to the below code
	// which is slightly faster. See TestSqDiff for correctness check.
	d := x - y
	return (d * d) >> 2
}
