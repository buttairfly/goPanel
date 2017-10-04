package module

import (
	"fmt"
	"image"
	"strconv"
	"strings"

	"github.com/buttairfly/goPanel/src/screen/raw"
)

var pre []string = []string{
	"x: ",
	"y: ",
	"c: ",
}

type jsonPoint string

func (jp jsonPoint) point() (image.Point, error) {
	const numPar = 2
	pointString := strings.Split(string(jp), ",")
	val := make([]int, numPar)
	var err error
	for i := 0; i < numPar; i++ {
		val[i], err = strconv.Atoi(strings.TrimLeft(pointString[i], pre[i]))
		if err != nil {
			return image.Point{}, fmt.Errorf("could not parse %v%v", pre[i], err)
		}
	}
	return image.Point{X: val[0], Y: val[1]}, nil
}

func marshalPoint(p image.Point) jsonPoint {
	return jsonPoint(fmt.Sprintf("%v%v, %v%v", pre[0], p.X, pre[1], p.Y))
}

type JsonColorPoint string

func (jcp JsonColorPoint) colorPoint() (ColorPoint, error) {
	const numPar = 2
	pointString := strings.Split(string(jcp), ",")
	val := make([]int, numPar)
	var err error
	for i := 0; i < numPar; i++ {
		val[i], err = strconv.Atoi(strings.TrimLeft(pointString[i], pre[i]))
		if err != nil {
			return ColorPoint{}, fmt.Errorf("could not parse %v%v", pre[i], err)
		}
	}
	return ColorPoint{
		image.Point{X: val[0], Y: val[1]},
		raw.RGB8Color(pointString[2]),
	}, nil
}

func marshalColorPoint(cp ColorPoint) JsonColorPoint {
	return JsonColorPoint(fmt.Sprintf("%v%v, %v%v, %v%v",
		pre[0], cp.X, pre[1], cp.Y, pre[2], cp.rgbType))
}
