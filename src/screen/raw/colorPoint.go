package raw

import (
	"fmt"
	"image"
)

type ColorPoint struct {
	image.Point
	C RGB8Color
}

func (cp ColorPoint) String() string {
	return fmt.Sprintf("%v %v", cp.Point, cp.C)
}
