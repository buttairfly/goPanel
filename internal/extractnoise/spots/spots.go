package spots

import (
	"image"
	"sort"
)

var _ Spots = &spots{}

// Spots are the led Panel lights
type Spots interface {
	sort.Interface
}

type spots struct {
	Points []image.Point
	Width  int
}

// NewSpots returns Spots from points
func NewSpots(width int, points []image.Point) Spots {
	return &spots{
		Points: points,
		Width:  width,
	}
}

func (s *spots) Len() int {
	return len(s.Points)
}

func (s *spots) Swap(i, j int) {
	s.Points[i], s.Points[j] = s.Points[j], s.Points[i]
}

func (s *spots) Less(i, j int) bool {
	I := s.Points[i].X + s.Width*s.Points[i].Y
	J := s.Points[j].X + s.Width*s.Points[j].Y
	return I < J
}
