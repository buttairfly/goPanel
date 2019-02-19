package hardware

import (
	"image"
	"image/color"
)

// Frame is a hardware frame
type Frame interface {
	image.Image
	toLedStripe() ledStipe
}

type frame struct {
	image   *image.RGBA
	modules []Module
}

// NewFrame return new Frame
func NewFrame(modules []Module) *frame {
	r := image.ZR
	for _, m := range modules {
		r = r.Union(m.Bounds())
	}
	return &frame{
		image:   image.NewRGBA(r),
		modules: modules,
	}
}

func (f *frame) toLedStripe() *ledStipe {
	return &ledStipe{}
}

// implements image interface
func (f *frame) ColorModel() color.Model {
	return f.image.ColorModel()
}

// implements image interface
func (f *frame) Bounds() image.Rectangle {
	return f.image.Bounds()
}

// implements image interface
func (f *frame) At(x, y int) color.Color {
	return f.image.At(x, y)
}
