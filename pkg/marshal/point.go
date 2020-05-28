package marshal

import "image"

// Point is the marshalable version of image.Point
type Point struct {
	X int `json:"px" yaml:"px"`
	Y int `json:"py" yaml:"py"`
}

// ToImagePoint converts a marshalable Point to image.Point
func (p *Point) ToImagePoint() image.Point {
	return image.Point{p.X, p.Y}
}

// FromImagePoint converts a image.Point to the marshal version
func FromImagePoint(ip image.Point) Point {
	return Point{ip.X, ip.Y}
}
