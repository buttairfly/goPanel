package imagefilter

import "image"

// Digits3x8 is an image array for all the number digits 0-9
var Digits3x8 = [...]image.Alpha{
	zero3x8, one3x8, two3x8, tree3x8, four3x8, five3x8, six3x8, seven3x8, eight3x8, nine3x8,
}

var zeroPix3x8 = []uint8{
	0xff, 0xff, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0xff, 0xff,
}

var onePix3x8 = []uint8{
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
}

var twoPix3x8 = []uint8{
	0xff, 0xff, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0xff, 0xff, 0xff,
	0xff, 0x00, 0x00,
	0xff, 0x00, 0x00,
	0xff, 0x00, 0x00,
	0xff, 0xff, 0xff,
}

var treePix3x8 = []uint8{
	0xff, 0xff, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0xff, 0xff, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0xff, 0xff, 0xff,
}

var fourPix3x8 = []uint8{
	0xff, 0x00, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0xff, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
}

var fivePix3x8 = []uint8{
	0xff, 0xff, 0xff,
	0xff, 0x00, 0x00,
	0xff, 0x00, 0x00,
	0xff, 0xff, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0xff, 0xff, 0xff,
}

var sixPix3x8 = []uint8{
	0xff, 0xff, 0xff,
	0xff, 0x00, 0x00,
	0xff, 0x00, 0x00,
	0xff, 0xff, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0xff, 0xff,
}

var sevenPix3x8 = []uint8{
	0xff, 0xff, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
}

var eightPix3x8 = []uint8{
	0xff, 0xff, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0xff, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0xff, 0xff,
}

var ninePix3x8 = []uint8{
	0xff, 0xff, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0x00, 0xff,
	0xff, 0xff, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0x00, 0x00, 0xff,
	0xff, 0xff, 0xff,
}

// helper vars to create image.Alpha
const x3x8 = 3
const y3x8 = 8

var rect3x8 = image.Rect(0, 0, x3x8, y3x8)
var zero3x8 = image.Alpha{
	Stride: x3x8,
	Rect:   rect3x8,
	Pix:    zeroPix3x8,
}

var one3x8 = image.Alpha{
	Stride: x3x8,
	Rect:   rect3x8,
	Pix:    onePix3x8,
}

var two3x8 = image.Alpha{
	Stride: x3x8,
	Rect:   rect3x8,
	Pix:    twoPix3x8,
}

var tree3x8 = image.Alpha{
	Stride: x3x8,
	Rect:   rect3x8,
	Pix:    treePix3x8,
}

var four3x8 = image.Alpha{
	Stride: x3x8,
	Rect:   rect3x8,
	Pix:    fourPix3x8,
}

var five3x8 = image.Alpha{
	Stride: x3x8,
	Rect:   rect3x8,
	Pix:    fivePix3x8,
}

var six3x8 = image.Alpha{
	Stride: x3x8,
	Rect:   rect3x8,
	Pix:    sixPix3x8,
}

var seven3x8 = image.Alpha{
	Stride: x3x8,
	Rect:   rect3x8,
	Pix:    sevenPix3x8,
}

var eight3x8 = image.Alpha{
	Stride: x3x8,
	Rect:   rect3x8,
	Pix:    eightPix3x8,
}

var nine3x8 = image.Alpha{
	Stride: x3x8,
	Rect:   rect3x8,
	Pix:    ninePix3x8,
}
