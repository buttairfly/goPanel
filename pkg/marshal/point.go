package marshal

import "image"

type Point struct {
	X int `json:"px" yaml:"px"`
	Y int `json:"py" yaml:"py"`
}

func (p *Point) ToImagePoint() image.Point {
	return image.Point{p.X, p.Y}
}

func FromImagePoint(ip image.Point) Point {
	return Point{ip.X, ip.Y}
}
