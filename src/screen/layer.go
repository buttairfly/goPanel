package screen

import (
	"fmt"
	"image"
	"image/draw"
)

type Layer interface {
	Blend(higher Layer) (Layer, error)
	GetZ() int
	GetImage() draw.Image
	GetBlendType() BlendType
	SetMask(mask draw.Image, maskOrigin image.Point)
}

type BlendType string

const (
	Overwrite BlendType = "overwrite"
)

type layer struct {
	zPos       int
	img        draw.Image
	blendType  BlendType
	mask       draw.Image
	maskOrigin image.Point
	//TODO: antiAliasing
}

func (l *layer) GetZ() int {
	return l.zPos
}

func (l *layer) GetImage() draw.Image {
	return l.img
}

func (l *layer) GetBlendType() BlendType {
	return l.blendType
}

func (l *layer) SetMask(mask draw.Image, maskOrigin image.Point) {
	l.mask = mask
	l.maskOrigin = maskOrigin
}

func (l *layer) Blend(higher Layer) (Layer, error) {
	switch l.blendType {
	case Overwrite:
		return l.overWrite(higher)
	default:
		return nil, fmt.Errorf("unknown blend type %v of layer %v", l.blendType, l.zPos)
	}
}

func (l *layer) overWrite(higher Layer) (Layer, error) {
	lowerRect := l.img.Bounds()
	higherRect := higher.GetImage().Bounds()
	if higherRect.Overlaps(lowerRect) {
		return nil, fmt.Errorf("higher layer %v is overlapping lower layer %v",
			higher.GetZ(), l.zPos)
	}

	if l.zPos >= higher.GetZ() {
		return nil, fmt.Errorf("higher layer %v must have a bigger z value than %v",
			higher.GetZ(), l.zPos)
	}
	draw.DrawMask(l.img, lowerRect, higher.GetImage(), l.img.Bounds().Min, l.mask, l.maskOrigin, draw.Over)

	//TODO: do return copy?!
	return l, nil
}
