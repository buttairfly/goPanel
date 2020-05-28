package marshal

import "image"

type Rectangle struct {
	Min Point `json:"min" yaml:"min"`
	Max Point `json:"max" yaml:"max"`
}

func (b *Rectangle) ToImageRectangle() image.Rectangle {
	return image.Rectangle{b.Min.ToImagePoint(), b.Max.ToImagePoint()}
}

func FromImageRectangle(ib image.Rectangle) Rectangle {
	return Rectangle{FromImagePoint(ib.Min), FromImagePoint(ib.Max)}
}
