package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"

	"github.com/buttairfly/goPanel/src/device"
	"github.com/buttairfly/goPanel/src/screen/raw"
)

type config struct {
	Modules []moduleConfig `json:"modules"`
}

type moduleConfig struct {
	DeviceName device.Name                `json:"deviceName"`
	Width      int                        `json:"width"`
	Height     int                        `json:"height"`
	Origin     image.Point                `json:"origin"`
	PixLUT     map[jsonPoint]int          `json:"pixLut"`
	PixCor     map[JsonColorPoint]float64 `json:"pixCorrection"`
	ColLUT     map[raw.RGB8Color]int      `json:"colLut"`
	Rotation   rotation                   `json:"rotation"`
	Mirror     mirror                     `json:"mirror"`
	LineOrder  lineOrder                  `json:"lineOrder"`
}
type rotation string

const (
	Rotate0   = rotation("0")
	Rotate90  = rotation("90")
	Rotate180 = rotation("180")
	Rotate270 = rotation("270")
)

type mirror string

const (
	MirrorNo = mirror("No")
	MirrorV  = mirror("Vertical")
	MirrorH  = mirror("Horizontal")
)

type lineOrder string

const (
	LineOrderXY     = lineOrder("XY")
	LineOrderSnake  = lineOrder("Snake")
	LineOrderManual = lineOrder("Manual")
)

const (
	pixLutUndef   = -1
	pixCorDefault = 1
)

func NewModulesFromConfig(path string) ([]module, error) {
	var c config
	err := c.FromFile(path)
	if err != nil {
		return nil, err
	}
	modules := make([]module, len(c.Modules))
	for i, module := range c.Modules {
		modules[i].deviceName = module.DeviceName
		modules[i].width = module.Width
		modules[i].height = module.Height
		modules[i].origin = module.Origin
		modules[i].numPix = len(module.PixLUT)
		if modules[i].pixLUT, err = module.generatePixLUT(); err != nil {
			return nil, fmt.Errorf("module %v: ", i, err)
		}
		if modules[i].pixCor, err = module.generatePixCor(); err != nil {
			return nil, fmt.Errorf("module %v: ", i, err)
		}
		if modules[i].colLUT, err = module.generateColLUT(); err != nil {
			return nil, fmt.Errorf("module %v: ", i, err)
		}
	}
	return modules, nil
}

func (c *config) FromReader(r io.Reader) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(&*c)
	if err != nil {
		return fmt.Errorf("can not decode json. error: %v", err)
	}
	return nil
}

func (c *config) FromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can not read config file %v. error: %v", path, err)
	}
	defer f.Close()
	return c.FromReader(f)
}

func (c *config) WriteToFile(path string) error {
	jsonConfig, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, jsonConfig, 0622)
}

func (mc *moduleConfig) generatePixCor() (pixCor map[ColorPoint]float64, err error) {
	pixCor = make(map[ColorPoint]float64, mc.Width*mc.Height)
	for y := 0; y < mc.Height; y++ {
		for x := 0; x < mc.Width; x++ {
			for c := range raw.RGB8Space {
				colorPoint := ColorPoint{Point: image.Point{x, y}, rgbType: c}
				if val, ok := mc.PixCor[marshalColorPoint(colorPoint)]; ok {
					if val > 1.0 || val < 0 {
						return nil, fmt.Errorf("value %v out of bounds at %+v", val, colorPoint)
					}
					pixCor[colorPoint] = val
				} else {
					pixCor[colorPoint] = pixCorDefault
				}
			}
		}
	}
	return
}

func (mc *moduleConfig) generateColLUT() (colLut map[raw.RGB8Color]int, err error) {
	colLut = make(map[raw.RGB8Color]int, len(raw.RGB8Space))
	for c := range raw.RGB8Space {
		if val, ok := mc.ColLUT[c]; ok {
			if val > len(raw.RGB8Space)-1 || val < 0 {
				return nil, fmt.Errorf("value %v out of bounds", val)
			}
			colLut[c] = val
		} else {
			return nil, fmt.Errorf("value %v not configured", c)
		}
	}

	return
}

func (mc *moduleConfig) generatePixLUT() (pixLUT map[image.Point]int, err error) {
	if mc.Rotation == Rotate90 || mc.Rotation == Rotate270 {
		mc.Width, mc.Height = mc.Height, mc.Width
	}

	pixLUT = make(map[image.Point]int, mc.Width*mc.Height)
	for y := 0; y < mc.Height; y++ {
		for x := 0; x < mc.Width; x++ {
			p := image.Point{x, y}
			var pos int
			if pos, err = mc.translateLineOrder(p); err != nil {
				return
			} else {
				if p, err = mc.translateRotation(p); err != nil {
					return
				}
				if p, err = mc.translateMirror(p); err != nil {
					return
				}
				if pos != pixLutUndef {
					pixLUT[p] = pos
				}
			}
		}
	}
	return
}

func (mc *moduleConfig) translateLineOrder(p image.Point) (pos int, err error) {
	switch mc.LineOrder {
	case LineOrderXY:
		pos = p.Y*mc.Width + p.X
	case LineOrderSnake:
		if isEven(p.Y) {
			p.X = invert(p.X, mc.Width)
		}
		pos = p.Y*mc.Width + p.X
	case LineOrderManual:
		if mc.PixLUT != nil {
			var ok bool
			pos, ok = mc.PixLUT[marshalPoint(p)]
			if !ok {
				pos = pixLutUndef
			}
		} else {
			err = errors.New("no PixLUT but LineOrder Manual set")
		}
	default:
		err = fmt.Errorf(
			"no correct LineOrder=(%v) set", mc.LineOrder)
	}
	return
}

func (mc *moduleConfig) translateMirror(p image.Point) (newP image.Point, err error) {
	newP = p
	switch mc.Mirror {
	case MirrorNo:
	case MirrorH:
		newP.Y = invert(p.Y, mc.Height)
	case MirrorV:
		newP.X = invert(p.X, mc.Width)
	default:
		err = fmt.Errorf(
			"no correct LineOrder=(%v) set", mc.LineOrder)
	}
	return
}

func (mc *moduleConfig) translateRotation(p image.Point) (newP image.Point, err error) {
	newP = p
	switch mc.Rotation {
	case Rotate0:
	case Rotate90:

	case Rotate180:
		newP.Y = invert(p.Y, mc.Height)
		newP.X = invert(p.X, mc.Width)
	case Rotate270:

	default:
		err = fmt.Errorf(
			"no correct LineOrder=(%v) set", mc.LineOrder)
	}
	return
}

func invert(val, maxVal int) int {
	return maxVal - val - 1
}

func isEven(val int) bool {
	return val%2 != 0
}

func (mc *moduleConfig) rotate(p image.Point, angle rotation) (newP image.Point) {
	midP := image.Point{
		X: mc.Width / 2,
		Y: mc.Height / 2,
	}
	newP.X = midP.X + ((p.X-midP.X)*cos(angle) - (p.Y-midP.Y)*sin(angle))
	newP.Y = midP.Y + ((p.Y-midP.Y)*cos(angle) + (p.X-midP.X)*sin(angle))

	return
}

func cos(angle rotation) int {
	switch angle {
	case Rotate0:
		return 1
	case Rotate180:
		return -1
	default:
		return 0
	}
}

func sin(angle rotation) int {
	switch angle {
	case Rotate90:
		return 1
	case Rotate270:
		return -1
	default:
		return 0
	}
}
