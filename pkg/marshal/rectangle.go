package marshal

import "image"

// Rectangle is the marshalable version of image.Rectangle
type Rectangle struct {
	Min Point `json:"min" yaml:"min"`
	Max Point `json:"max" yaml:"max"`
}

// ToImageRectangle converts a marshalable Rectangle to image.Rectangle
func (b *Rectangle) ToImageRectangle() image.Rectangle {
	return image.Rectangle{b.Min.ToImagePoint(), b.Max.ToImagePoint()}
}

// FromImageRectangle converts a image.Rectangle to the marshal version
func FromImageRectangle(ib image.Rectangle) Rectangle {
	return Rectangle{FromImagePoint(ib.Min), FromImagePoint(ib.Max)}
}
