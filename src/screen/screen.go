package screen

/*
 * screen has a byte buffer containing the image
 * the image needs to be set via the Set(.) method
 */

import (
	"fmt"
	"image"
	"image/color"

	"github.com/buttairfly/goPanel/src/device"
	"github.com/buttairfly/goPanel/src/screen/module"
	"github.com/buttairfly/goPanel/src/screen/raw"
)

type screen struct {
	device        device.SpiDevice
	width, height int
	numPix        int //actual number of pixel w/o dead space
	modules       []Module
	stride        int
	pix           []byte
}

type Module interface {
	Serialize(image image.Image) []byte
	Bounds() image.Rectangle
	GetDevice() device.Type
	GetNumPix() int
}

func New(configFile string, name device.Type) (*screen, error) {
	configuredModules, err := module.NewModulesFromConfig(configFile)
	if err != nil {
		return nil, err
	}
	s := new(screen)
	s.modules = make([]Module, 0)
	r := image.Rectangle{}
	numPix := 0
	for _, m := range configuredModules {
		if m.GetDevice() == name {
			s.modules = append(s.modules, &m)
			numPix += m.GetNumPix()
			r = r.Union(m.Bounds()).Canon()
			if r.Empty() || r.Min.X < 0 || r.Min.Y < 0 || r.Max.X < 0 || r.Max.Y < 0 {
				return nil, fmt.Errorf("module %v out of bounds %v", m, r)
			}
		}
	}
	if r.Empty() || r.Min.X != 0 || r.Min.Y != 0 {
		return nil, fmt.Errorf("screen out of bounds %v", r)
	}
	s.width = r.Dx()
	s.height = r.Dy()
	s.stride = s.width * raw.RGB8NumBytes
	s.numPix = numPix
	s.device, err = device.NewSpiDevice(name, numPix)
	s.pix = make([]byte, raw.RGB8NumBytes*s.width*s.height)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *screen) Bounds() image.Rectangle {
	return image.Rect(0, 0, s.width, s.height)
}

func (s *screen) ColorModel() color.Model {
	return raw.RGB8Model
}

func (s *screen) At(x, y int) color.Color {
	return s.RGB8At(x, y)
}

func (s *screen) RGB8At(x, y int) raw.RGB8 {
	if !(image.Point{x, y}.In(s.Bounds())) {
		return raw.RGB8{}
	}
	i := s.PixOffset(x, y)
	return raw.RGB8{R: s.pix[i+0], G: s.pix[i+1], B: s.pix[i+2]}
}

func (s *screen) PixOffset(x, y int) int {
	return y*s.stride + x*raw.RGB8NumBytes
}

func (s *screen) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(s.Bounds())) {
		return
	}
	i := s.PixOffset(x, y)
	c1 := color.RGBAModel.Convert(c).(color.RGBA)
	s.pix[i+0] = c1.R
	s.pix[i+1] = c1.G
	s.pix[i+2] = c1.B
}

func (s *screen) SetRGB8(x, y int, c raw.RGB8) {
	if !(image.Point{x, y}.In(s.Bounds())) {
		return
	}
	i := s.PixOffset(x, y)
	s.pix[i+0] = c.R
	s.pix[i+1] = c.G
	s.pix[i+2] = c.B

}

func (s *screen) Write() error {
	buffer := make([]byte, 0)
	for _, m := range s.modules {
		buffer = append(buffer, m.Serialize(s)...)
	}
	if n, err := s.device.Write(buffer); err != nil {
		return err
	} else if n != s.numPix*raw.RGB8NumBytes {
		return fmt.Errorf(
			"written bytes do not match %v != %v", n, s.numPix*raw.RGB8NumBytes)
	}
	return nil
}
