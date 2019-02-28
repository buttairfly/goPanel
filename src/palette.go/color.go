package palette

import "image/color"

// Color is the enhanced color functionality
type Color interface {
	color.Color
	Equals(c color.Color) bool
}

// RGBA represents a traditional 32-bit alpha-premultiplied color, having 8
// bits for each of red, green, blue and alpha.
//
// An alpha-premultiplied color component C has been scaled by alpha (A), so
// has valid values 0 <= C <= A.
type RGBA struct {
	R, G, B, A uint8
}

// RGBA implements color.Color interface
func (c RGBA) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = uint32(c.A)
	a |= a << 8
	return
}

// Equals checks wheather the color of another color is identical
func (c RGBA) Equals(c2 color.Color) bool {
	cr, cg, cb, ca := c.RGBA()
	c2r, c2g, c2b, c2a := c2.RGBA()
	return cr == c2r && cg == c2g && cb == c2b && ca == c2a
}

// RGBA64 represents a 64-bit alpha-premultiplied color, having 16 bits for
// each of red, green, blue and alpha.
//
// An alpha-premultiplied color component C has been scaled by alpha (A), so
// has valid values 0 <= C <= A.
type RGBA64 struct {
	R, G, B, A uint16
}

// RGBA implements color.Color interface
func (c RGBA64) RGBA() (r, g, b, a uint32) {
	return uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
}

// Equals checks wheather the color of another color is identical
func (c RGBA64) Equals(c2 color.Color) bool {
	cr, cg, cb, ca := c.RGBA()
	c2r, c2g, c2b, c2a := c2.RGBA()
	return cr == c2r && cg == c2g && cb == c2b && ca == c2a
}

// NRGBA represents a non-alpha-premultiplied 32-bit color.
type NRGBA struct {
	R, G, B, A uint8
}

// RGBA implements color.Color interface
func (c NRGBA) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	r *= uint32(c.A)
	r /= 0xff
	g = uint32(c.G)
	g |= g << 8
	g *= uint32(c.A)
	g /= 0xff
	b = uint32(c.B)
	b |= b << 8
	b *= uint32(c.A)
	b /= 0xff
	a = uint32(c.A)
	a |= a << 8
	return
}

// Equals checks wheather the color of another color is identical
func (c NRGBA) Equals(c2 color.Color) bool {
	cr, cg, cb, ca := c.RGBA()
	c2r, c2g, c2b, c2a := c2.RGBA()
	return cr == c2r && cg == c2g && cb == c2b && ca == c2a
}

// NRGBA64 represents a non-alpha-premultiplied 64-bit color,
// having 16 bits for each of red, green, blue and alpha.
type NRGBA64 struct {
	R, G, B, A uint16
}

// RGBA implements color.Color interface
func (c NRGBA64) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r *= uint32(c.A)
	r /= 0xffff
	g = uint32(c.G)
	g *= uint32(c.A)
	g /= 0xffff
	b = uint32(c.B)
	b *= uint32(c.A)
	b /= 0xffff
	a = uint32(c.A)
	return
}

// Equals checks wheather the color of another color is identical
func (c NRGBA64) Equals(c2 color.Color) bool {
	cr, cg, cb, ca := c.RGBA()
	c2r, c2g, c2b, c2a := c2.RGBA()
	return cr == c2r && cg == c2g && cb == c2b && ca == c2a
}

// Alpha represents an 8-bit alpha color.
type Alpha struct {
	A uint8
}

// RGBA implements color.Color interface
func (c Alpha) RGBA() (r, g, b, a uint32) {
	a = uint32(c.A)
	a |= a << 8
	return a, a, a, a
}

// Equals checks wheather the color of another color is identical
func (c Alpha) Equals(c2 color.Color) bool {
	cr, cg, cb, ca := c.RGBA()
	c2r, c2g, c2b, c2a := c2.RGBA()
	return cr == c2r && cg == c2g && cb == c2b && ca == c2a
}

// Alpha16 represents a 16-bit alpha color.
type Alpha16 struct {
	A uint16
}

// RGBA implements color.Color interface
func (c Alpha16) RGBA() (r, g, b, a uint32) {
	a = uint32(c.A)
	return a, a, a, a
}

// Equals checks wheather the color of another color is identical
func (c Alpha16) Equals(c2 color.Color) bool {
	cr, cg, cb, ca := c.RGBA()
	c2r, c2g, c2b, c2a := c2.RGBA()
	return cr == c2r && cg == c2g && cb == c2b && ca == c2a
}

// Gray represents an 8-bit grayscale color.
type Gray struct {
	Y uint8
}

// RGBA implements color.Color interface
func (c Gray) RGBA() (r, g, b, a uint32) {
	y := uint32(c.Y)
	y |= y << 8
	return y, y, y, 0xffff
}

// Equals checks wheather the color of another color is identical
func (c Gray) Equals(c2 color.Color) bool {
	cr, cg, cb, ca := c.RGBA()
	c2r, c2g, c2b, c2a := c2.RGBA()
	return cr == c2r && cg == c2g && cb == c2b && ca == c2a
}

// Gray16 represents a 16-bit grayscale color.
type Gray16 struct {
	Y uint16
}

// RGBA implements color.Color interface
func (c Gray16) RGBA() (r, g, b, a uint32) {
	y := uint32(c.Y)
	return y, y, y, 0xffff
}

// Equals checks wheather the color of another color is identical
func (c Gray16) Equals(c2 color.Color) bool {
	cr, cg, cb, ca := c.RGBA()
	c2r, c2g, c2b, c2a := c2.RGBA()
	return cr == c2r && cg == c2g && cb == c2b && ca == c2a
}

// Model can convert any Color to one from its own color model. The conversion
// may be lossy.
type Model interface {
	Convert(c Color) Color
}

// ModelFunc returns a Model that invokes f to implement the conversion.
func ModelFunc(f func(Color) Color) Model {
	// Note: using *modelFunc as the implementation
	// means that callers can still use comparisons
	// like m == RGBAModel. This is not possible if
	// we use the func value directly, because funcs
	// are no longer comparable.
	return &modelFunc{f}
}

type modelFunc struct {
	f func(Color) Color
}

func (m *modelFunc) Convert(c Color) Color {
	return m.f(c)
}

// Models for the standard color types.
var (
	RGBAModel    = ModelFunc(rgbaModel)
	RGBA64Model  = ModelFunc(rgba64Model)
	NRGBAModel   = ModelFunc(nrgbaModel)
	NRGBA64Model = ModelFunc(nrgba64Model)
	AlphaModel   = ModelFunc(alphaModel)
	Alpha16Model = ModelFunc(alpha16Model)
	GrayModel    = ModelFunc(grayModel)
	Gray16Model  = ModelFunc(gray16Model)
)

func rgbaModel(c Color) Color {
	if _, ok := c.(RGBA); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	return RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func rgba64Model(c Color) Color {
	if _, ok := c.(RGBA64); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	return RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func nrgbaModel(c Color) Color {
	if _, ok := c.(NRGBA); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	if a == 0xffff {
		return NRGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 0xff}
	}
	if a == 0 {
		return NRGBA{0, 0, 0, 0}
	}
	// Since Color.RGBA returns an alpha-premultiplied color, we should have r <= a && g <= a && b <= a.
	r = (r * 0xffff) / a
	g = (g * 0xffff) / a
	b = (b * 0xffff) / a
	return NRGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func nrgba64Model(c Color) Color {
	if _, ok := c.(NRGBA64); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	if a == 0xffff {
		return NRGBA64{uint16(r), uint16(g), uint16(b), 0xffff}
	}
	if a == 0 {
		return NRGBA64{0, 0, 0, 0}
	}
	// Since Color.RGBA returns an alpha-premultiplied color, we should have r <= a && g <= a && b <= a.
	r = (r * 0xffff) / a
	g = (g * 0xffff) / a
	b = (b * 0xffff) / a
	return NRGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func alphaModel(c Color) Color {
	if _, ok := c.(Alpha); ok {
		return c
	}
	_, _, _, a := c.RGBA()
	return Alpha{uint8(a >> 8)}
}

func alpha16Model(c Color) Color {
	if _, ok := c.(Alpha16); ok {
		return c
	}
	_, _, _, a := c.RGBA()
	return Alpha16{uint16(a)}
}

func grayModel(c Color) Color {
	if _, ok := c.(Gray); ok {
		return c
	}
	r, g, b, _ := c.RGBA()

	// These coefficients (the fractions 0.299, 0.587 and 0.114) are the same
	// as those given by the JFIF specification and used by func RGBToYCbCr in
	// ycbcr.go.
	//
	// Note that 19595 + 38470 + 7471 equals 65536.
	//
	// The 24 is 16 + 8. The 16 is the same as used in RGBToYCbCr. The 8 is
	// because the return value is 8 bit color, not 16 bit color.
	y := (19595*r + 38470*g + 7471*b + 1<<15) >> 24

	return Gray{uint8(y)}
}

func gray16Model(c Color) Color {
	if _, ok := c.(Gray16); ok {
		return c
	}
	r, g, b, _ := c.RGBA()

	// These coefficients (the fractions 0.299, 0.587 and 0.114) are the same
	// as those given by the JFIF specification and used by func RGBToYCbCr in
	// ycbcr.go.
	//
	// Note that 19595 + 38470 + 7471 equals 65536.
	y := (19595*r + 38470*g + 7471*b + 1<<15) >> 16

	return Gray16{uint16(y)}
}
