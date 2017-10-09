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
	DeviceType device.Type                `json:"deviceType"`
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
	Rotate0   rotation = "0"
	Rotate90  rotation = "90"
	Rotate180 rotation = "180"
	Rotate270 rotation = "270"
)

type mirror string

const (
	MirrorNo mirror = "No"
	MirrorV  mirror = "Vertical"
	MirrorH  mirror = "Horizontal"
)

type lineOrder string

const (
	LineOrderXY     lineOrder = "XY"
	LineOrderSnake  lineOrder = "Snake"
	LineOrderManual lineOrder = "Manual"
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
		modules[i].deviceName = module.DeviceType
		modules[i].width = module.Width
		modules[i].height = module.Height
		modules[i].origin = module.Origin
		if modules[i].pixLUT, err = module.generatePixLUT(); err != nil {
			return nil, fmt.Errorf("module %v: %v", i, err)
		}
		modules[i].numPix = len(modules[i].pixLUT)
		if modules[i].pixCor, err = module.generatePixCor(); err != nil {
			return nil, fmt.Errorf("module %v: %v", i, err)
		}
		if modules[i].colLUT, err = module.generateColLUT(); err != nil {
			return nil, fmt.Errorf("module %v: %v", i, err)
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

func (mc *moduleConfig) generatePixCor() (pixCor map[raw.ColorPoint]float64, err error) {
	pixCor = make(map[raw.ColorPoint]float64, mc.Width*mc.Height)
	for y := 0; y < mc.Height; y++ {
		for x := 0; x < mc.Width; x++ {
			for c := range raw.RGB8Space {
				colorPoint := raw.ColorPoint{Point: image.Point{X: x, Y: y}, C: c}
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
			p := image.Point{X: x, Y: y}
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
			"no correct Mirror=(%v) set", mc.Mirror)
	}
	return
}

func (mc *moduleConfig) translateRotation(p image.Point) (newP image.Point, err error) {
	switch mc.Rotation {
	case Rotate0:
		newP = p
	case Rotate90:
		newP.Y = p.X
		newP.X = invert(p.Y, mc.Height)
	case Rotate180:
		newP.Y = invert(p.Y, mc.Height)
		newP.X = invert(p.X, mc.Width)
	case Rotate270:
		newP.Y = invert(p.X, mc.Width)
		newP.X = p.Y
	default:
		err = fmt.Errorf(
			"no correct Rotation=(%v) set", mc.Rotation)
	}
	return
}

func invert(val, maxVal int) int {
	return maxVal - val - 1
}

func isEven(val int) bool {
	return val%2 != 0
}
